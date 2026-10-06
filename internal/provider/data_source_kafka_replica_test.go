// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import "testing"

func TestKafkaReplicaDataSourceSchema(t *testing.T) {
	if len(kafkaRawSchema("", true, true).Attributes) == 0 {
		t.Fatal("replica schema has no attributes")
	}
}
