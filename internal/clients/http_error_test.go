// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package clients

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestDecodeAPIError(t *testing.T) {
	response := &http.Response{
		StatusCode: http.StatusConflict,
		Body:       io.NopCloser(strings.NewReader(`{"error_code":409,"message":"conflict"}`)),
	}

	err := DecodeAPIError(response)
	apiError, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiError.StatusCode != http.StatusConflict || apiError.Message != "conflict" {
		t.Fatalf("unexpected API error: %#v", apiError)
	}
	if apiError.NotFound() {
		t.Fatal("conflict must not be reported as not found")
	}
}

func TestDecodeAPIErrorNotFound(t *testing.T) {
	response := &http.Response{
		StatusCode: http.StatusNotFound,
		Body:       io.NopCloser(strings.NewReader("not found")),
	}

	err := DecodeAPIError(response)
	if !err.(*APIError).NotFound() {
		t.Fatal("expected not-found API error")
	}
}

func TestEscapePathSegment(t *testing.T) {
	if got := EscapePathSegment("principal/User:alice"); got != "principal%2FUser:alice" {
		t.Fatalf("unexpected escaped path segment: %s", got)
	}
}

func TestSanitizeURL(t *testing.T) {
	got := sanitizeURL("https://example.test/v1?access_token=secret&subject=orders")
	if strings.Contains(got, "secret") || !strings.Contains(got, "subject=orders") {
		t.Fatalf("URL was not safely sanitized: %s", got)
	}
}
