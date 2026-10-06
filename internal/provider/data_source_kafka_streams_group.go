// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/variananora/terraform-provider-confluent-platform/internal/clients"
)

type KafkaStreamsGroupDataSource struct{ client *clients.HTTPClient }
type KafkaStreamsGroupDataSourceModel struct {
	ClusterID    types.String `tfsdk:"cluster_id"`
	Name         types.String `tfsdk:"name"`
	ID           types.String `tfsdk:"id"`
	ResponseJSON types.String `tfsdk:"response_json"`
}

func NewKafkaStreamsGroupDataSource() datasource.DataSource { return &KafkaStreamsGroupDataSource{} }
func (d *KafkaStreamsGroupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	rawMetadata("", "kafka_streams_group", req, resp)
}
func (d *KafkaStreamsGroupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = kafkaRawSchema("Reads Kafka Streams group metadata from REST Proxy v3.", false, false)
}
func (d *KafkaStreamsGroupDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	configureKafkaRaw(&d.client, req, resp)
}
func (d *KafkaStreamsGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data KafkaStreamsGroupDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	readKafkaRaw(ctx, d.client, req, resp, rawPath(data.ClusterID.ValueString(), "streams-groups", data.Name.ValueString()), "Kafka Streams group")
}
