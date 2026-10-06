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

func TestKafkaTopicJSONRoundTrip(t *testing.T) {
	original := KafkaTopicCreateRequest{
		TopicName:         "orders",
		PartitionsCount:   3,
		ReplicationFactor: 2,
		Configs: []KafkaTopicConfig{
			{Name: "cleanup.policy", Value: "compact"},
		},
	}

	payload, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal topic request: %s", err)
	}

	var decoded KafkaTopicCreateRequest
	if err := json.Unmarshal(payload, &decoded); err != nil {
		t.Fatalf("unmarshal topic request: %s", err)
	}
	if decoded.TopicName != original.TopicName || decoded.PartitionsCount != original.PartitionsCount || len(decoded.Configs) != 1 {
		t.Fatalf("unexpected decoded topic request: %#v", decoded)
	}
}

func TestDeleteTopicConfig(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", request.Method)
		}
		if request.URL.EscapedPath() != "/kafka/v3/clusters/cluster-1/topics/orders%2Fv1/configs/cleanup.policy" {
			t.Errorf("unexpected path: %s", request.URL.EscapedPath())
		}
		writer.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client, err := NewHTTPClient(server.URL, time.Second)
	if err != nil {
		t.Fatalf("create client: %s", err)
	}
	if err := client.DeleteTopicConfig(context.Background(), "cluster-1", "orders/v1", "cleanup.policy"); err != nil {
		t.Fatalf("delete topic config: %s", err)
	}
}
