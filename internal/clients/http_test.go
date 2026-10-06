// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package clients

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHTTPClientDo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", request.Method)
		}

		if request.URL.Path != "/v3/clusters" {
			t.Fatalf("expected /v3/clusters, got %s", request.URL.Path)
		}

		writer.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, err := NewHTTPClient(server.URL+"/", time.Second)
	if err != nil {
		t.Fatalf("unexpected client error: %s", err)
	}

	response, err := client.Do(context.Background(), http.MethodGet, "/v3/clusters", nil, http.Header{
		"Accept": []string{"application/json"},
	})
	if err != nil {
		t.Fatalf("unexpected request error: %s", err)
	}
	defer response.Body.Close()
}

func TestNewHTTPClientRejectsInvalidEndpoint(t *testing.T) {
	if _, err := NewHTTPClient("localhost:8082", time.Second); err == nil {
		t.Fatal("expected endpoint validation error")
	}
}
