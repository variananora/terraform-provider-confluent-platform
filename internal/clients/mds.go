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

const mdsBasePath = "/security/1.0"

func (c *HTTPClient) GetMDSActiveNodes(ctx context.Context, protocol string) ([]byte, error) {
	return c.GetRawJSON(ctx, fmt.Sprintf("%s/activenodes/%s", mdsBasePath, EscapePathSegment(protocol)))
}

type MDSAuthResponse struct {
	AuthToken string `json:"auth_token"`
	TokenType string `json:"token_type"`
	ExpiresIn int64  `json:"expires_in"`
}

func (c *HTTPClient) AuthenticateMDS(ctx context.Context) (MDSAuthResponse, error) {
	var response MDSAuthResponse
	if err := c.mdsJSON(ctx, http.MethodGet, mdsBasePath+"/authenticate", nil, &response); err != nil {
		return MDSAuthResponse{}, err
	}
	return response, nil
}

func (c *HTTPClient) GetMDSMetadataClusterID(ctx context.Context) ([]byte, error) {
	return c.GetRawJSON(ctx, mdsBasePath+"/metadataClusterId")
}

func (c *HTTPClient) GetMDSJSON(ctx context.Context, path string) ([]byte, error) {
	return c.GetRawJSON(ctx, path)
}

func (c *HTTPClient) PostMDSJSON(ctx context.Context, path string, payload any) ([]byte, error) {
	var response json.RawMessage
	if err := c.mdsJSON(ctx, http.MethodPost, path, payload, &response); err != nil {
		return nil, err
	}
	return response, nil
}

func (c *HTTPClient) CreateMDSClusterRegistry(ctx context.Context, configuration any) error {
	return c.mdsJSON(ctx, http.MethodPost, mdsBasePath+"/registry/clusters", []any{configuration}, nil)
}

func (c *HTTPClient) GetMDSClusterRegistry(ctx context.Context, clusterName string) ([]byte, error) {
	path := fmt.Sprintf("%s/registry/clusters/%s", mdsBasePath, EscapePathSegment(clusterName))
	return c.GetRawJSON(ctx, path)
}

func (c *HTTPClient) DeleteMDSClusterRegistry(ctx context.Context, clusterName string) error {
	path := fmt.Sprintf("%s/registry/clusters/%s", mdsBasePath, EscapePathSegment(clusterName))
	response, err := c.Do(ctx, http.MethodDelete, path, nil, http.Header{"Accept": []string{"application/json"}})
	if err != nil {
		return err
	}
	defer response.Body.Close()
	return DecodeAPIError(response)
}

func (c *HTTPClient) GetMDSAuditConfig(ctx context.Context) ([]byte, error) {
	return c.GetRawJSON(ctx, mdsBasePath+"/audit/config")
}

func (c *HTTPClient) PutMDSAuditConfig(ctx context.Context, configuration any) ([]byte, error) {
	var response json.RawMessage
	if err := c.mdsJSON(ctx, http.MethodPut, mdsBasePath+"/audit/config", configuration, &response); err != nil {
		return nil, err
	}
	return response, nil
}

func (c *HTTPClient) GetMDSAuditRoutes(ctx context.Context, query string) ([]byte, error) {
	path := mdsBasePath + "/audit/routes?q=" + EscapeQueryValue(query)
	return c.GetRawJSON(ctx, path)
}

func (c *HTTPClient) GetMDSAuditLookup(ctx context.Context, crn string) ([]byte, error) {
	path := mdsBasePath + "/audit/lookup?crn=" + EscapeQueryValue(crn)
	return c.GetRawJSON(ctx, path)
}

type MDSRoleBindingScope struct {
	Clusters map[string]string `json:"clusters,omitempty"`
}

type MDSResourcePattern struct {
	ResourceType string `json:"resourceType"`
	Name         string `json:"name"`
	PatternType  string `json:"patternType"`
}

type MDSRoleBindingRequest struct {
	Scope            MDSRoleBindingScope  `json:"scope"`
	ResourcePatterns []MDSResourcePattern `json:"resourcePatterns,omitempty"`
}

type MDSACLBinding struct {
	Pattern MDSACLPattern `json:"pattern"`
	Entry   MDSACLEntry   `json:"entry"`
}

type MDSACLPattern struct {
	ResourceType string `json:"resourceType"`
	Name         string `json:"name"`
	PatternType  string `json:"patternType"`
}

type MDSACLEntry struct {
	Principal      string `json:"principal"`
	Host           string `json:"host"`
	Operation      string `json:"operation"`
	PermissionType string `json:"permissionType"`
}

type MDSACLRequest struct {
	Scope      MDSRoleBindingScope `json:"scope"`
	ACLBinding MDSACLBinding       `json:"aclBinding"`
}

type MDSACLFilter struct {
	Scope            MDSRoleBindingScope `json:"scope"`
	ACLBindingFilter MDSACLBindingFilter `json:"aclBindingFilter"`
}

type MDSACLBindingFilter struct {
	PatternFilter MDSACLPattern `json:"patternFilter"`
	EntryFilter   MDSACLEntry   `json:"entryFilter"`
}

func (c *HTTPClient) CreateMDSRoleBinding(ctx context.Context, principal, roleName string, request MDSRoleBindingRequest) error {
	if err := ValidateMDSRoleBindingScope(request.Scope); err != nil {
		return err
	}
	path := fmt.Sprintf("%s/principals/%s/roles/%s", mdsBasePath, EscapePathSegment(principal), EscapePathSegment(roleName))
	if len(request.ResourcePatterns) > 0 {
		path += "/bindings"
	}
	return c.mdsJSON(ctx, http.MethodPost, path, request, nil)
}

func (c *HTTPClient) DeleteMDSRoleBinding(ctx context.Context, principal, roleName string, request MDSRoleBindingRequest) error {
	if err := ValidateMDSRoleBindingScope(request.Scope); err != nil {
		return err
	}
	path := fmt.Sprintf("%s/principals/%s/roles/%s", mdsBasePath, EscapePathSegment(principal), EscapePathSegment(roleName))
	if len(request.ResourcePatterns) > 0 {
		path += "/bindings"
	}
	return c.mdsJSON(ctx, http.MethodDelete, path, request, nil)
}

func (c *HTTPClient) ReadMDSRoleBindingResources(ctx context.Context, principal, roleName string, scope MDSRoleBindingScope) ([]MDSResourcePattern, error) {
	if err := ValidateMDSRoleBindingScope(scope); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("%s/principals/%s/roles/%s/resources", mdsBasePath, EscapePathSegment(principal), EscapePathSegment(roleName))
	var resources []MDSResourcePattern
	err := c.mdsJSON(ctx, http.MethodPost, path, scope, &resources)
	return resources, err
}

func ValidateMDSRoleBindingScope(scope MDSRoleBindingScope) error {
	if len(scope.Clusters) == 0 {
		return fmt.Errorf("MDS role binding scope requires clusters.kafka-cluster")
	}
	if _, ok := scope.Clusters["kafka-cluster"]; !ok {
		return fmt.Errorf("MDS role binding scope requires clusters.kafka-cluster")
	}
	return nil
}

func (c *HTTPClient) CreateMDSACL(ctx context.Context, request MDSACLRequest) error {
	return c.mdsJSON(ctx, http.MethodPost, mdsBasePath+"/acls", request, nil)
}

func (c *HTTPClient) DeleteMDSACL(ctx context.Context, request MDSACLFilter) error {
	return c.mdsJSON(ctx, http.MethodDelete, mdsBasePath+"/acls", request, nil)
}

func (c *HTTPClient) SearchMDSACLs(ctx context.Context, request MDSACLFilter) ([]MDSACLBinding, error) {
	var bindings []MDSACLBinding
	err := c.mdsJSON(ctx, http.MethodPost, mdsBasePath+"/acls:search", request, &bindings)
	return bindings, err
}

func (c *HTTPClient) mdsJSON(ctx context.Context, method, path string, payload, target any) error {
	var body *bytes.Reader
	if payload == nil {
		body = bytes.NewReader(nil)
	} else {
		encoded, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("encode MDS request: %w", err)
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
		return fmt.Errorf("decode MDS response: %w", err)
	}
	return nil
}
