// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package clients

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestConnectorLifecycleMethods(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method == http.MethodPost && request.URL.Path == "/connectors" {
			var payload ConnectConnectorCreateRequest
			if err := json.NewDecoder(request.Body).Decode(&payload); err != nil || payload.Name != "orders" {
				t.Errorf("unexpected create payload: %#v, %v", payload, err)
			}
			_, _ = writer.Write([]byte(`{"name":"orders","type":"sink","config":{"name":"orders","password":"********"}}`))
			return
		}
		if request.Method == http.MethodPut && request.URL.Path == "/connectors/orders/config" {
			_, _ = writer.Write([]byte(`{"name":"orders","password":"********"}`))
			return
		}
		if request.Method == http.MethodGet && request.URL.Path == "/connectors/orders/status" {
			_, _ = writer.Write([]byte(`{"name":"orders","connector":{"state":"RUNNING"},"tasks":[]}`))
			return
		}
		http.NotFound(writer, request)
	}))
	defer server.Close()

	client, err := NewHTTPClient(server.URL, time.Second)
	if err != nil {
		t.Fatalf("create client: %s", err)
	}
	connector, err := client.CreateConnector(context.Background(), ConnectConnectorCreateRequest{Name: "orders", Config: map[string]string{"name": "orders"}})
	if err != nil || connector.Name != "orders" {
		t.Fatalf("unexpected connector: %#v, %v", connector, err)
	}
	if _, err := client.UpdateConnectorConfig(context.Background(), "orders", map[string]string{"name": "orders"}); err != nil {
		t.Fatalf("update connector: %s", err)
	}
	status, err := client.GetConnectorStatus(context.Background(), "orders")
	if err != nil || status.Connector.State != "RUNNING" {
		t.Fatalf("unexpected connector status: %#v, %v", status, err)
	}
}
