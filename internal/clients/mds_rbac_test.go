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

func TestMDSRoleBindingAndACLPaths(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method == http.MethodPost && request.URL.Path == "/security/1.0/principals/User:alice/roles/DeveloperRead/bindings" {
			writer.WriteHeader(http.StatusNoContent)
			return
		}
		if request.Method == http.MethodPost && request.URL.Path == "/security/1.0/acls" {
			writer.WriteHeader(http.StatusNoContent)
			return
		}
		http.NotFound(writer, request)
	}))
	defer server.Close()

	client, err := NewHTTPClient(server.URL, time.Second)
	if err != nil {
		t.Fatalf("create client: %s", err)
	}
	scope := MDSRoleBindingScope{Clusters: map[string]string{"kafka-cluster": "kafka-1"}}
	if err := client.CreateMDSRoleBinding(context.Background(), "User:alice", "DeveloperRead", MDSRoleBindingRequest{Scope: scope, ResourcePatterns: []MDSResourcePattern{{ResourceType: "Topic", Name: "orders", PatternType: "LITERAL"}}}); err != nil {
		t.Fatalf("create role binding: %s", err)
	}
	if err := client.CreateMDSACL(context.Background(), MDSACLRequest{Scope: scope}); err != nil {
		t.Fatalf("create ACL: %s", err)
	}
}
