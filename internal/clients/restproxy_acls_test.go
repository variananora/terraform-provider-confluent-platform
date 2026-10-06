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

func TestKafkaACLClientLifecycle(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.EscapedPath() != "/kafka/v3/clusters/cluster-1/acls/TOPIC/orders/LITERAL/User:alice/%2A/READ/ALLOW" {
			t.Fatalf("unexpected ACL path: %s", request.URL.EscapedPath())
		}
		if request.Method == http.MethodGet {
			_, _ = writer.Write([]byte(`{"cluster_id":"cluster-1","resource_type":"TOPIC","resource_name":"orders","pattern_type":"LITERAL","principal":"User:alice","host":"*","operation":"READ","permission":"ALLOW"}`))
			return
		}
		writer.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client, err := NewHTTPClient(server.URL, time.Second)
	if err != nil {
		t.Fatalf("create client: %s", err)
	}
	acl := KafkaACL{ResourceType: "TOPIC", ResourceName: "orders", PatternType: "LITERAL", Principal: "User:alice", Host: "*", Operation: "READ", Permission: "ALLOW"}
	if _, err := client.GetKafkaACL(context.Background(), "cluster-1", acl); err != nil {
		t.Fatalf("get ACL: %s", err)
	}
	if err := client.DeleteKafkaACL(context.Background(), "cluster-1", acl); err != nil {
		t.Fatalf("delete ACL: %s", err)
	}
}
