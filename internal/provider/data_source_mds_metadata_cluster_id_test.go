// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import "testing"

func TestMDSMetadataClusterIDDataSourceExists(t *testing.T) {
	if NewMDSMetadataClusterIDDataSource() == nil {
		t.Fatal("expected MDS metadata cluster ID data source")
	}
}
