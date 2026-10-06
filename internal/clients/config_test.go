// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package clients

import (
	"context"
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRuntimeAppliesBearerAuthAndUserAgent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Header.Get("Authorization") != "Bearer token" {
			t.Errorf("expected bearer authorization, got %q", request.Header.Get("Authorization"))
		}
		if request.UserAgent() != "test-agent" {
			t.Errorf("expected test-agent, got %q", request.UserAgent())
		}
	}))
	defer server.Close()

	runtime, err := NewRuntime(RuntimeConfig{
		Endpoints: map[string]string{"rest_proxy": server.URL},
		Auth:      AuthConfig{BearerToken: "token"},
		UserAgent: "test-agent",
		Timeout:   time.Second,
	})
	if err != nil {
		t.Fatalf("unexpected runtime error: %s", err)
	}

	client, ok := runtime.Service("rest_proxy")
	if !ok {
		t.Fatal("expected REST Proxy client")
	}
	response, err := client.Do(context.Background(), http.MethodGet, "/health", nil, nil)
	if err != nil {
		t.Fatalf("unexpected request error: %s", err)
	}
	response.Body.Close()
}

func TestRuntimeRejectsInvalidClientCertificate(t *testing.T) {
	_, err := NewRuntime(RuntimeConfig{
		Endpoints: map[string]string{"rest_proxy": "https://localhost"},
		TLS: TLSConfig{
			ClientCertPEM: "invalid",
			ClientKeyPEM:  "invalid",
		},
	})
	if err == nil {
		t.Fatal("expected client certificate error")
	}
}

func TestNewTransportUsesTLS12Minimum(t *testing.T) {
	transport, err := newTransport(TLSConfig{})
	if err != nil {
		t.Fatalf("unexpected transport error: %s", err)
	}
	tlsTransport := transport.(*http.Transport)
	if tlsTransport.TLSClientConfig.MinVersion != tls.VersionTLS12 {
		t.Fatalf("expected TLS 1.2 minimum, got %d", tlsTransport.TLSClientConfig.MinVersion)
	}
}
