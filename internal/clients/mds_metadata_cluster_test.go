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

func TestGetMDSMetadataClusterID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/security/1.0/metadataClusterId" {
			t.Errorf("unexpected path: %s", request.URL.Path)
		}
		_, _ = writer.Write([]byte(`"cluster-id-1"`))
	}))
	defer server.Close()

	client, err := NewHTTPClient(server.URL, time.Second)
	if err != nil {
		t.Fatalf("create client: %s", err)
	}
	payload, err := client.GetMDSMetadataClusterID(context.Background())
	if err != nil || string(payload) != `"cluster-id-1"` {
		t.Fatalf("unexpected metadata cluster ID: %s, %v", payload, err)
	}
}
