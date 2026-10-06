// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import "testing"

func TestKafkaBrokerDataSourceSchema(t *testing.T) {
	if (&KafkaBrokerDataSourceModel{}).BrokerID.ValueInt64() != 0 {
		t.Fatal("unexpected broker model value")
	}
}
