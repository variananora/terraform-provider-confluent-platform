// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import "testing"

func TestKafkaClusterDataSourceSchema(t *testing.T) {
	if len(KafkaClusterDataSourceModel{}.ClusterID.ValueString()) != 0 {
		t.Fatal("unexpected cluster model value")
	}
}
