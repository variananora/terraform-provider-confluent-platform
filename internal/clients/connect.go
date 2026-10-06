// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package clients

import (
	"context"
	"fmt"
	"net/http"
)

type ConnectConnector struct {
	Name   string            `json:"name"`
	Type   string            `json:"type"`
	Class  string            `json:"class"`
	Config map[string]string `json:"config"`
}

type ConnectConnectorCreateRequest struct {
	Name   string            `json:"name"`
	Config map[string]string `json:"config"`
}

type ConnectConnectorStatus struct {
	Name      string                `json:"name"`
	Connector ConnectConnectorState `json:"connector"`
	Tasks     []ConnectTaskStatus   `json:"tasks"`
	Type      string                `json:"type"`
}

type ConnectConnectorState struct {
	State  string `json:"state"`
	Worker string `json:"worker_id"`
}

type ConnectTaskStatus struct {
	ID     int64  `json:"id"`
	State  string `json:"state"`
	Worker string `json:"worker_id"`
	Trace  string `json:"trace,omitempty"`
}

const connectBasePath = "/connectors"

func (c *HTTPClient) CreateConnector(ctx context.Context, request ConnectConnectorCreateRequest) (ConnectConnector, error) {
	var connector ConnectConnector
	err := c.doJSON(ctx, http.MethodPost, connectBasePath, request, &connector)
	return connector, err
}

func (c *HTTPClient) GetConnector(ctx context.Context, name string) (ConnectConnector, error) {
	var connector ConnectConnector
	err := c.getJSON(ctx, fmt.Sprintf("%s/%s", connectBasePath, EscapePathSegment(name)), &connector)
	return connector, err
}

func (c *HTTPClient) UpdateConnectorConfig(ctx context.Context, name string, config map[string]string) (map[string]string, error) {
	var updated map[string]string
	err := c.doJSON(ctx, http.MethodPut, fmt.Sprintf("%s/%s/config", connectBasePath, EscapePathSegment(name)), config, &updated)
	return updated, err
}

func (c *HTTPClient) DeleteConnector(ctx context.Context, name string) error {
	response, err := c.Do(ctx, http.MethodDelete, fmt.Sprintf("%s/%s", connectBasePath, EscapePathSegment(name)), nil, http.Header{"Accept": []string{"application/json"}})
	if err != nil {
		return err
	}
	defer response.Body.Close()
	return DecodeAPIError(response)
}

func (c *HTTPClient) GetConnectorStatus(ctx context.Context, name string) (ConnectConnectorStatus, error) {
	var status ConnectConnectorStatus
	err := c.getJSON(ctx, fmt.Sprintf("%s/%s/status", connectBasePath, EscapePathSegment(name)), &status)
	return status, err
}
