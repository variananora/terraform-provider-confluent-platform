// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package clients

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"time"
)

type AuthConfig struct {
	Username    string
	Password    string
	BearerToken string
}

type TLSConfig struct {
	CACertPEM          string
	ClientCertPEM      string
	ClientKeyPEM       string
	ServerName         string
	InsecureSkipVerify bool
}

type HTTPClientConfig struct {
	BaseURL    string
	Timeout    time.Duration
	MaxRetries int64
	UserAgent  string
	Auth       AuthConfig
	Transport  http.RoundTripper
}

type RuntimeConfig struct {
	Endpoints  map[string]string
	Timeout    time.Duration
	MaxRetries int64
	UserAgent  string
	Auth       AuthConfig
	TLS        TLSConfig
}

type Runtime struct {
	services map[string]*HTTPClient
}

func NewRuntime(config RuntimeConfig) (*Runtime, error) {
	transport, err := newTransport(config.TLS)
	if err != nil {
		return nil, err
	}

	runtime := &Runtime{services: make(map[string]*HTTPClient)}
	for service, endpoint := range config.Endpoints {
		if endpoint == "" {
			continue
		}

		client, err := NewHTTPClientWithConfig(HTTPClientConfig{
			BaseURL:    endpoint,
			Timeout:    config.Timeout,
			MaxRetries: config.MaxRetries,
			UserAgent:  config.UserAgent,
			Auth:       config.Auth,
			Transport:  transport,
		})
		if err != nil {
			return nil, fmt.Errorf("configure %s client: %w", service, err)
		}
		runtime.services[service] = client
	}

	return runtime, nil
}

func (r *Runtime) Service(name string) (*HTTPClient, bool) {
	client, ok := r.services[name]
	return client, ok
}

func (r *Runtime) StandardClient() *http.Client {
	for _, client := range r.services {
		return client.StandardClient()
	}
	return http.DefaultClient
}

func newTransport(config TLSConfig) (http.RoundTripper, error) {
	rootCAs, err := x509.SystemCertPool()
	if err != nil || rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}
	if config.CACertPEM != "" && !rootCAs.AppendCertsFromPEM([]byte(config.CACertPEM)) {
		return nil, fmt.Errorf("append custom CA certificate")
	}

	tlsConfig := &tls.Config{
		RootCAs:            rootCAs,
		ServerName:         config.ServerName,
		InsecureSkipVerify: config.InsecureSkipVerify,
		MinVersion:         tls.VersionTLS12,
	}
	if config.ClientCertPEM != "" || config.ClientKeyPEM != "" {
		certificate, err := tls.X509KeyPair([]byte(config.ClientCertPEM), []byte(config.ClientKeyPEM))
		if err != nil {
			return nil, fmt.Errorf("configure client certificate: %w", err)
		}
		tlsConfig.Certificates = []tls.Certificate{certificate}
	}

	return &http.Transport{TLSClientConfig: tlsConfig}, nil
}
