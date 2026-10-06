// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const restProxyV3Path = "/kafka/v3/clusters"

type KafkaCluster struct {
	ClusterID        string `json:"cluster_id"`
	Controller       int64  `json:"controller"`
	KafkaClusterID   string `json:"kafka_cluster_id"`
	KafkaClusterName string `json:"kafka_cluster_name"`
	MetadataVersion  string `json:"metadata_version"`
}

type KafkaBroker struct {
	BrokerID  int64  `json:"broker_id"`
	ClusterID string `json:"cluster_id"`
	Host      string `json:"host"`
	Port      int64  `json:"port"`
	Rack      string `json:"rack"`
}

type KafkaTopic struct {
	ClusterID         string             `json:"cluster_id"`
	TopicName         string             `json:"topic_name"`
	IsInternal        bool               `json:"is_internal"`
	PartitionsCount   int64              `json:"partitions_count"`
	ReplicationFactor int64              `json:"replication_factor"`
	Configs           []KafkaTopicConfig `json:"configs"`
}

type KafkaTopicConfig struct {
	Name        string `json:"name"`
	Value       string `json:"value"`
	IsReadOnly  bool   `json:"is_read_only"`
	IsSensitive bool   `json:"is_sensitive"`
}

type KafkaTopicCreateRequest struct {
	TopicName         string             `json:"topic_name"`
	PartitionsCount   int64              `json:"partitions_count,omitempty"`
	ReplicationFactor int64              `json:"replication_factor,omitempty"`
	Configs           []KafkaTopicConfig `json:"configs,omitempty"`
}

type KafkaTopicUpdateRequest struct {
	PartitionsCount *int64             `json:"partitions_count,omitempty"`
	Configs         []KafkaTopicConfig `json:"configs,omitempty"`
}

func (c *HTTPClient) GetCluster(ctx context.Context, clusterID string) (KafkaCluster, error) {
	var cluster KafkaCluster
	err := c.getJSON(ctx, fmt.Sprintf("%s/%s", restProxyV3Path, EscapePathSegment(clusterID)), &cluster)
	return cluster, err
}

func (c *HTTPClient) GetBroker(ctx context.Context, clusterID string, brokerID int64) (KafkaBroker, error) {
	var broker KafkaBroker
	path := fmt.Sprintf("%s/%s/brokers/%d", restProxyV3Path, EscapePathSegment(clusterID), brokerID)
	err := c.getJSON(ctx, path, &broker)
	return broker, err
}

func (c *HTTPClient) GetTopic(ctx context.Context, clusterID string, topicName string) (KafkaTopic, error) {
	var topic KafkaTopic
	path := fmt.Sprintf("%s/%s/topics/%s", restProxyV3Path, EscapePathSegment(clusterID), EscapePathSegment(topicName))
	err := c.getJSON(ctx, path, &topic)
	return topic, err
}

func (c *HTTPClient) CreateTopic(ctx context.Context, clusterID string, request KafkaTopicCreateRequest) (KafkaTopic, error) {
	var topic KafkaTopic
	err := c.doJSON(ctx, http.MethodPost, fmt.Sprintf("%s/%s/topics", restProxyV3Path, EscapePathSegment(clusterID)), request, &topic)
	return topic, err
}

func (c *HTTPClient) UpdateTopicConfig(ctx context.Context, clusterID string, topicName string, configName string, value string) error {
	path := fmt.Sprintf("%s/%s/topics/%s/configs/%s", restProxyV3Path, EscapePathSegment(clusterID), EscapePathSegment(topicName), EscapePathSegment(configName))
	return c.doJSON(ctx, http.MethodPut, path, map[string]string{"value": value}, nil)
}

func (c *HTTPClient) DeleteTopicConfig(ctx context.Context, clusterID string, topicName string, configName string) error {
	path := fmt.Sprintf("%s/%s/topics/%s/configs/%s", restProxyV3Path, EscapePathSegment(clusterID), EscapePathSegment(topicName), EscapePathSegment(configName))
	response, err := c.Do(ctx, http.MethodDelete, path, nil, http.Header{"Accept": []string{"application/json"}})
	if err != nil {
		return err
	}
	defer response.Body.Close()
	return DecodeAPIError(response)
}

func (c *HTTPClient) DeleteTopic(ctx context.Context, clusterID string, topicName string) error {
	response, err := c.Do(ctx, http.MethodDelete, fmt.Sprintf("%s/%s/topics/%s", restProxyV3Path, EscapePathSegment(clusterID), EscapePathSegment(topicName)), nil, http.Header{"Accept": []string{"application/json"}})
	if err != nil {
		return err
	}
	defer response.Body.Close()
	return DecodeAPIError(response)
}

func (c *HTTPClient) GetRawJSON(ctx context.Context, path string) ([]byte, error) {
	response, err := c.Do(ctx, http.MethodGet, path, nil, http.Header{"Accept": []string{"application/json"}})
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if apiErr := DecodeAPIError(response); apiErr != nil {
		return nil, apiErr
	}
	var payload json.RawMessage
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("decode REST Proxy response: %w", err)
	}
	return payload, nil
}

func (c *HTTPClient) getJSON(ctx context.Context, path string, target any) error {
	return c.doJSON(ctx, http.MethodGet, path, nil, target)
}

func (c *HTTPClient) doJSON(ctx context.Context, method string, path string, payload any, target any) error {
	var body *bytes.Reader
	if payload == nil {
		body = bytes.NewReader(nil)
	} else {
		encoded, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("encode REST Proxy request: %w", err)
		}
		body = bytes.NewReader(encoded)
	}
	response, err := c.Do(ctx, method, path, body, http.Header{"Accept": []string{"application/json"}, "Content-Type": []string{"application/json"}})
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if apiErr := DecodeAPIError(response); apiErr != nil {
		return apiErr
	}
	if target == nil || response.StatusCode == http.StatusNoContent {
		return nil
	}
	if err := json.NewDecoder(response.Body).Decode(target); err != nil {
		return fmt.Errorf("decode REST Proxy response: %w", err)
	}
	return nil
}
