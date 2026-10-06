// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import "testing"

func TestConnectReadDataSourcesExist(t *testing.T) {
	if NewConnectClusterDataSource() == nil || NewConnectConnectorInfoDataSource() == nil || NewConnectConnectorConfigDataSource() == nil {
		t.Fatal("expected Connect read data sources")
	}
}
