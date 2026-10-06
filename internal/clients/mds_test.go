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

func TestGetMDSActiveNodes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", request.Method)
		}
		if request.URL.EscapedPath() != "/security/1.0/activenodes/https" {
			t.Errorf("unexpected path: %s", request.URL.EscapedPath())
		}
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write([]byte(`{"nodes":[{"host":"mds-1","port":8090}]}`))
	}))
	defer server.Close()

	client, err := NewHTTPClient(server.URL, time.Second)
	if err != nil {
		t.Fatalf("create client: %s", err)
	}
	payload, err := client.GetMDSActiveNodes(context.Background(), "https")
	if err != nil {
		t.Fatalf("get active nodes: %s", err)
	}
	if string(payload) != `{"nodes":[{"host":"mds-1","port":8090}]}` {
		t.Fatalf("unexpected response: %s", payload)
	}
}
