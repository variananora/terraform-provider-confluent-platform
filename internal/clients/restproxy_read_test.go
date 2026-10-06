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

func TestRESTProxyReadMethods(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", request.Method)
		}
		writer.Header().Set("Content-Type", "application/json")
		switch request.URL.EscapedPath() {
		case "/kafka/v3/clusters/cluster-1":
			_, _ = writer.Write([]byte(`{"cluster_id":"cluster-1","controller":1}`))
		case "/kafka/v3/clusters/cluster-1/brokers/1":
			_, _ = writer.Write([]byte(`{"broker_id":1,"cluster_id":"cluster-1","host":"broker-1","port":9092}`))
		case "/kafka/v3/clusters/cluster-1/topics/orders%2Fv1":
			_, _ = writer.Write([]byte(`{"cluster_id":"cluster-1","topic_name":"orders/v1","partitions_count":3}`))
		default:
			http.NotFound(writer, request)
		}
	}))
	defer server.Close()

	client, err := NewHTTPClient(server.URL, time.Second)
	if err != nil {
		t.Fatalf("create client: %s", err)
	}

	cluster, err := client.GetCluster(context.Background(), "cluster-1")
	if err != nil || cluster.ClusterID != "cluster-1" {
		t.Fatalf("unexpected cluster: %#v, %v", cluster, err)
	}
	broker, err := client.GetBroker(context.Background(), "cluster-1", 1)
	if err != nil || broker.Host != "broker-1" {
		t.Fatalf("unexpected broker: %#v, %v", broker, err)
	}
	topic, err := client.GetTopic(context.Background(), "cluster-1", "orders/v1")
	if err != nil || topic.TopicName != "orders/v1" {
		t.Fatalf("unexpected topic: %#v, %v", topic, err)
	}
}
