// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"

	"github.com/variananora/terraform-provider-confluent-platform/internal/clients"
)

func mdsErrorMessage(err error) string {
	if apiErr, ok := err.(*clients.APIError); ok {
		if apiErr.Unauthorized() {
			return "MDS authentication failed (401). Check Basic credentials, bearer token, or MDS authentication configuration."
		}
		if apiErr.Forbidden() {
			return "MDS authorization failed (403). The calling principal lacks the required MDS role or ACL permission."
		}
	}
	return err.Error()
}

func restProxyClient(providerData any) (*clients.HTTPClient, error) {
	runtime, ok := providerData.(*clients.Runtime)
	if !ok {
		return nil, fmt.Errorf("expected *clients.Runtime, got %T", providerData)
	}
	client, ok := runtime.Service("rest_proxy")
	if !ok {
		return nil, fmt.Errorf("REST Proxy endpoint is not configured")
	}
	return client, nil
}
