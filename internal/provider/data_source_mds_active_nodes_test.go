// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import "testing"

func TestMDSActiveNodesDataSourceSchema(t *testing.T) {
	dataSource := NewMDSActiveNodesDataSource()
	if dataSource == nil {
		t.Fatal("expected MDS active nodes data source")
	}
}
