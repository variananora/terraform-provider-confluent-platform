// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/variananora/terraform-provider-confluent-platform/internal/clients"
)

type KafkaClusterLinkDataSource struct{ client *clients.HTTPClient }
type KafkaClusterLinkDataSourceModel struct {
	ClusterID    types.String `tfsdk:"cluster_id"`
	Name         types.String `tfsdk:"name"`
	ID           types.String `tfsdk:"id"`
	ResponseJSON types.String `tfsdk:"response_json"`
}

func NewKafkaClusterLinkDataSource() datasource.DataSource { return &KafkaClusterLinkDataSource{} }
func (d *KafkaClusterLinkDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	rawMetadata("", "kafka_cluster_link", req, resp)
}
func (d *KafkaClusterLinkDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = kafkaRawSchema("Reads Kafka cluster link metadata from REST Proxy v3.", false, false)
}
func (d *KafkaClusterLinkDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	configureKafkaRaw(&d.client, req, resp)
}
func (d *KafkaClusterLinkDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data KafkaClusterLinkDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	readKafkaRaw(ctx, d.client, req, resp, rawPath(data.ClusterID.ValueString(), "links", data.Name.ValueString()), "Kafka cluster link")
}
