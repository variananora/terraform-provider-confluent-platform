// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import "testing"

func TestKafkaConsumerGroupDataSourceSchema(t *testing.T) {
	if len(kafkaRawSchema("", false, false).Attributes) == 0 {
		t.Fatal("consumer group schema has no attributes")
	}
}
