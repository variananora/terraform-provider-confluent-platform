// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import "testing"

func TestKafkaRawPath(t *testing.T) {
	if got := rawPath("cluster", "topics", "orders"); got != "/kafka/v3/clusters/cluster/topics/orders" {
		t.Fatalf("unexpected raw path: %s", got)
	}
}
