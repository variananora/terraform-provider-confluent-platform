// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import "testing"

func TestConnectConnectorStatusDataSourceExists(t *testing.T) {
	if NewConnectConnectorStatusDataSource() == nil {
		t.Fatal("expected connector status data source")
	}
}
