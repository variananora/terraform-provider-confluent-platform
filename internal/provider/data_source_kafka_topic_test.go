// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import "testing"

func TestKafkaTopicDataSourceSchema(t *testing.T) {
	if (&KafkaTopicDataSourceModel{}).Name.ValueString() != "" {
		t.Fatal("unexpected topic model value")
	}
}
