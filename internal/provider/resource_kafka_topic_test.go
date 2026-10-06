// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import "testing"

func TestKafkaTopicResourceSchema(t *testing.T) {
	if len(KafkaTopicResourceSchema().Attributes) == 0 {
		t.Fatal("topic resource schema has no attributes")
	}
}
