// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/variananora/terraform-provider-confluent-platform/internal/clients"
)

type KafkaShareGroupDataSource struct{ client *clients.HTTPClient }
type KafkaShareGroupDataSourceModel struct {
	ClusterID    types.String `tfsdk:"cluster_id"`
	Name         types.String `tfsdk:"name"`
	ID           types.String `tfsdk:"id"`
	ResponseJSON types.String `tfsdk:"response_json"`
}

func NewKafkaShareGroupDataSource() datasource.DataSource { return &KafkaShareGroupDataSource{} }
func (d *KafkaShareGroupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	rawMetadata("", "kafka_share_group", req, resp)
}
func (d *KafkaShareGroupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = kafkaRawSchema("Reads Kafka share group metadata from REST Proxy v3.", false, false)
}
func (d *KafkaShareGroupDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	configureKafkaRaw(&d.client, req, resp)
}
func (d *KafkaShareGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data KafkaShareGroupDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	readKafkaRaw(ctx, d.client, req, resp, rawPath(data.ClusterID.ValueString(), "share-groups", data.Name.ValueString()), "Kafka share group")
}
