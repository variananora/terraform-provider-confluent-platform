// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import "testing"

func TestKafkaStreamsGroupDataSourceSchema(t *testing.T) {
	if len(kafkaRawSchema("", false, false).Attributes) == 0 {
		t.Fatal("streams group schema has no attributes")
	}
}
