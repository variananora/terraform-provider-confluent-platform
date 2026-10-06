// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type HTTPClient struct {
	baseURL    string
	client     *http.Client
	auth       AuthConfig
	userAgent  string
	maxRetries int64
}

type APIError struct {
	StatusCode int
	ErrorCode  any    `json:"error_code,omitempty"`
	Message    string `json:"message,omitempty"`
	Body       string
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("Confluent Platform API returned HTTP %d: %s", e.StatusCode, e.Message)
	}
	return fmt.Sprintf("Confluent Platform API returned HTTP %d", e.StatusCode)
}

func (e *APIError) NotFound() bool {
	return e.StatusCode == http.StatusNotFound
}

func (e *APIError) Unauthorized() bool {
	return e.StatusCode == http.StatusUnauthorized
}

func (e *APIError) Forbidden() bool {
	return e.StatusCode == http.StatusForbidden
}

func (c *HTTPClient) StandardClient() *http.Client {
	return c.client
}

func NewHTTPClient(baseURL string, timeout time.Duration) (*HTTPClient, error) {
	return NewHTTPClientWithConfig(HTTPClientConfig{
		BaseURL: baseURL,
		Timeout: timeout,
	})
}

func NewHTTPClientWithConfig(config HTTPClientConfig) (*HTTPClient, error) {
	parsedURL, err := url.Parse(strings.TrimSpace(config.BaseURL))
	if err != nil {
		return nil, fmt.Errorf("parse endpoint: %w", err)
	}

	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return nil, fmt.Errorf("endpoint must include a scheme and host")
	}

	if config.Timeout <= 0 {
		config.Timeout = 30 * time.Second
	}
	if config.MaxRetries < 0 {
		config.MaxRetries = 0
	}

	return &HTTPClient{
		baseURL:    strings.TrimRight(parsedURL.String(), "/"),
		client:     &http.Client{Transport: config.Transport, Timeout: config.Timeout},
		auth:       config.Auth,
		userAgent:  config.UserAgent,
		maxRetries: config.MaxRetries,
	}, nil
}

func (c *HTTPClient) Do(ctx context.Context, method string, path string, body io.Reader, headers http.Header) (*http.Response, error) {
	requestURL := c.baseURL + "/" + strings.TrimLeft(path, "/")
	request, err := http.NewRequestWithContext(ctx, method, requestURL, body)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	for key, values := range headers {
		for _, value := range values {
			request.Header.Add(key, value)
		}
	}
	if c.auth.BearerToken != "" {
		request.Header.Set("Authorization", "Bearer "+c.auth.BearerToken)
	} else if c.auth.Username != "" {
		request.SetBasicAuth(c.auth.Username, c.auth.Password)
	}
	if c.userAgent != "" {
		request.Header.Set("User-Agent", c.userAgent)
	}

	started := time.Now()
	tflog.Debug(ctx, "sending Confluent Platform API request", map[string]any{
		"method": method,
		"url":    sanitizeURL(requestURL),
	})

	response, err := c.doWithRetry(request)
	if err != nil {
		tflog.Debug(ctx, "Confluent Platform API request failed", map[string]any{
			"method":   method,
			"url":      sanitizeURL(requestURL),
			"duration": time.Since(started).String(),
			"error":    err.Error(),
		})
		return nil, err
	}

	tflog.Debug(ctx, "received Confluent Platform API response", map[string]any{
		"method":   method,
		"url":      sanitizeURL(requestURL),
		"status":   response.StatusCode,
		"duration": time.Since(started).String(),
	})

	return response, nil
}

func DecodeAPIError(response *http.Response) error {
	if response == nil || response.StatusCode < http.StatusBadRequest {
		return nil
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return &APIError{StatusCode: response.StatusCode, Body: err.Error()}
	}
	response.Body.Close()
	response.Body = io.NopCloser(strings.NewReader(string(body)))

	apiError := &APIError{StatusCode: response.StatusCode, Body: string(body)}
	if err := json.Unmarshal(body, apiError); err != nil && len(body) > 0 {
		apiError.Message = strings.TrimSpace(string(body))
	}
	return apiError
}

func EscapePathSegment(value string) string {
	return url.PathEscape(value)
}

func EscapeQueryValue(value string) string {
	return url.QueryEscape(value)
}

func sanitizeURL(rawURL string) string {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}
	query := parsedURL.Query()
	for key := range query {
		lowerKey := strings.ToLower(key)
		if strings.Contains(lowerKey, "token") || strings.Contains(lowerKey, "password") || strings.Contains(lowerKey, "secret") || strings.Contains(lowerKey, "key") {
			query.Set(key, "[REDACTED]")
		}
	}
	parsedURL.RawQuery = query.Encode()
	return parsedURL.String()
}

func (c *HTTPClient) doWithRetry(request *http.Request) (*http.Response, error) {
	var response *http.Response
	var err error

	for attempt := int64(0); attempt <= c.maxRetries; attempt++ {
		if attempt > 0 && request.GetBody != nil {
			request.Body, err = request.GetBody()
			if err != nil {
				return nil, err
			}
		}

		response, err = c.client.Do(request)
		if err != nil {
			if !isRetryableMethod(request.Method) || attempt == c.maxRetries {
				return nil, err
			}
		} else if !isRetryableStatus(response.StatusCode) || attempt == c.maxRetries {
			return response, nil
		} else {
			response.Body.Close()
		}

		select {
		case <-request.Context().Done():
			return nil, request.Context().Err()
		case <-time.After(time.Duration(attempt+1) * 100 * time.Millisecond):
		}
	}

	return response, err
}

func isRetryableMethod(method string) bool {
	return method == http.MethodGet || method == http.MethodHead || method == http.MethodOptions
}

func isRetryableStatus(statusCode int) bool {
	return statusCode == http.StatusRequestTimeout || statusCode == http.StatusTooManyRequests || statusCode >= 500
}
