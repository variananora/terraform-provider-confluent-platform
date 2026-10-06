// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/variananora/terraform-provider-confluent-platform/internal/clients"
)

type KafkaClusterDataSource struct{ client *clients.HTTPClient }
type KafkaClusterDataSourceModel struct {
	ID               types.String `tfsdk:"id"`
	ClusterID        types.String `tfsdk:"cluster_id"`
	Controller       types.Int64  `tfsdk:"controller"`
	KafkaClusterID   types.String `tfsdk:"kafka_cluster_id"`
	KafkaClusterName types.String `tfsdk:"kafka_cluster_name"`
	MetadataVersion  types.String `tfsdk:"metadata_version"`
}

func NewKafkaClusterDataSource() datasource.DataSource { return &KafkaClusterDataSource{} }
func (d *KafkaClusterDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_kafka_cluster"
}
func (d *KafkaClusterDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{MarkdownDescription: "Reads Kafka cluster metadata from REST Proxy v3.", Attributes: map[string]schema.Attribute{"id": schema.StringAttribute{Computed: true}, "cluster_id": schema.StringAttribute{Required: true}, "controller": schema.Int64Attribute{Computed: true}, "kafka_cluster_id": schema.StringAttribute{Computed: true}, "kafka_cluster_name": schema.StringAttribute{Computed: true}, "metadata_version": schema.StringAttribute{Computed: true}}}
}
func (d *KafkaClusterDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, err := restProxyClient(req.ProviderData)
	if err != nil {
		resp.Diagnostics.AddError("Unable to configure Kafka cluster data source", err.Error())
		return
	}
	d.client = client
}
func (d *KafkaClusterDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data KafkaClusterDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	cluster, err := d.client.GetCluster(ctx, data.ClusterID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to read Kafka cluster", err.Error())
		return
	}
	data.ID = types.StringValue(cluster.ClusterID)
	data.Controller = types.Int64Value(cluster.Controller)
	data.KafkaClusterID = types.StringValue(cluster.KafkaClusterID)
	data.KafkaClusterName = types.StringValue(cluster.KafkaClusterName)
	data.MetadataVersion = types.StringValue(cluster.MetadataVersion)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
