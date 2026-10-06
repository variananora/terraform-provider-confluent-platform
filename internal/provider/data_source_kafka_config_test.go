// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import "testing"

func TestKafkaConfigDataSourceSchema(t *testing.T) {
	if len(kafkaRawSchema("", true, false).Attributes) == 0 {
		t.Fatal("config schema has no attributes")
	}
}
