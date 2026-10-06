// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package clients

import (
	"context"
	"fmt"
	"net/http"
)

type KafkaACL struct {
	ClusterID    string `json:"cluster_id"`
	ResourceType string `json:"resource_type"`
	ResourceName string `json:"resource_name"`
	PatternType  string `json:"pattern_type"`
	Principal    string `json:"principal"`
	Host         string `json:"host"`
	Operation    string `json:"operation"`
	Permission   string `json:"permission"`
}

type KafkaACLCreateRequest struct {
	ResourceType string `json:"resource_type"`
	ResourceName string `json:"resource_name"`
	PatternType  string `json:"pattern_type"`
	Principal    string `json:"principal"`
	Host         string `json:"host"`
	Operation    string `json:"operation"`
	Permission   string `json:"permission"`
}

func (c *HTTPClient) CreateKafkaACL(ctx context.Context, clusterID string, request KafkaACLCreateRequest) (KafkaACL, error) {
	var acl KafkaACL
	err := c.doJSON(ctx, http.MethodPost, fmt.Sprintf("%s/%s/acls", restProxyV3Path, EscapePathSegment(clusterID)), request, &acl)
	return acl, err
}

func (c *HTTPClient) GetKafkaACL(ctx context.Context, clusterID string, acl KafkaACL) (KafkaACL, error) {
	var result KafkaACL
	err := c.getJSON(ctx, kafkaACLPath(clusterID, acl), &result)
	return result, err
}

func (c *HTTPClient) DeleteKafkaACL(ctx context.Context, clusterID string, acl KafkaACL) error {
	response, err := c.Do(ctx, http.MethodDelete, kafkaACLPath(clusterID, acl), nil, http.Header{"Accept": []string{"application/json"}})
	if err != nil {
		return err
	}
	defer response.Body.Close()
	return DecodeAPIError(response)
}

func kafkaACLPath(clusterID string, acl KafkaACL) string {
	return fmt.Sprintf("%s/%s/acls/%s/%s/%s/%s/%s/%s/%s",
		restProxyV3Path,
		EscapePathSegment(clusterID),
		EscapePathSegment(acl.ResourceType),
		EscapePathSegment(acl.ResourceName),
		EscapePathSegment(acl.PatternType),
		EscapePathSegment(acl.Principal),
		EscapePathSegment(acl.Host),
		EscapePathSegment(acl.Operation),
		EscapePathSegment(acl.Permission),
	)
}
