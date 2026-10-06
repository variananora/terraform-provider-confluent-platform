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

func TestAuthenticateMDS(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/security/1.0/authenticate" {
			t.Fatalf("unexpected path: %s", request.URL.Path)
		}
		_, _ = writer.Write([]byte(`{"auth_token":"token","token_type":"Bearer","expires_in":900}`))
	}))
	defer server.Close()

	client, err := NewHTTPClient(server.URL, time.Second)
	if err != nil {
		t.Fatalf("create client: %s", err)
	}
	response, err := client.AuthenticateMDS(context.Background())
	if err != nil || response.AuthToken != "token" || response.TokenType != "Bearer" {
		t.Fatalf("unexpected authentication response: %#v, %v", response, err)
	}
}

func TestValidateMDSRoleBindingScope(t *testing.T) {
	if err := ValidateMDSRoleBindingScope(MDSRoleBindingScope{Clusters: map[string]string{"connect-cluster": "id"}}); err == nil {
		t.Fatal("expected missing kafka cluster scope error")
	}
}
