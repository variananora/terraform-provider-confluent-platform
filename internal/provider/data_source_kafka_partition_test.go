// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package provider

import "testing"

func TestKafkaPartitionDataSourceSchema(t *testing.T) {
	if len(kafkaRawSchema("", true, true).Attributes) == 0 {
		t.Fatal("partition schema has no attributes")
	}
}
