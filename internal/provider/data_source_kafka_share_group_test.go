// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import "testing"

func TestKafkaShareGroupDataSourceSchema(t *testing.T) {
	if len(kafkaRawSchema("", false, false).Attributes) == 0 {
		t.Fatal("share group schema has no attributes")
	}
}
