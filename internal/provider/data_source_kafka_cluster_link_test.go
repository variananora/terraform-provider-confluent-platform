// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import "testing"

func TestKafkaClusterLinkDataSourceSchema(t *testing.T) {
	if len(kafkaRawSchema("", false, false).Attributes) == 0 {
		t.Fatal("cluster link schema has no attributes")
	}
}
