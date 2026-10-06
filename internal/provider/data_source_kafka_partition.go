// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/variananora/terraform-provider-confluent-platform/internal/clients"
)

type KafkaPartitionDataSource struct{ client *clients.HTTPClient }
type KafkaPartitionDataSourceModel struct {
	ClusterID    types.String `tfsdk:"cluster_id"`
	Name         types.String `tfsdk:"name"`
	PartitionID  types.Int64  `tfsdk:"partition_id"`
	ID           types.String `tfsdk:"id"`
	ResponseJSON types.String `tfsdk:"response_json"`
}

func NewKafkaPartitionDataSource() datasource.DataSource { return &KafkaPartitionDataSource{} }
func (d *KafkaPartitionDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	rawMetadata("", "kafka_partition", req, resp)
}
func (d *KafkaPartitionDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = kafkaRawSchema("Reads Kafka partition metadata from REST Proxy v3.", true, true)
}
func (d *KafkaPartitionDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	configureKafkaRaw(&d.client, req, resp)
}
func (d *KafkaPartitionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data KafkaPartitionDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	readKafkaRaw(ctx, d.client, req, resp, fmt.Sprintf("%s/partitions/%d", rawPath(data.ClusterID.ValueString(), "topics", data.Name.ValueString()), data.PartitionID.ValueInt64()), "Kafka partition")
}
