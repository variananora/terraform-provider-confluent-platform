// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/variananora/terraform-provider-confluent-platform/internal/clients"
)

type KafkaBrokerDataSource struct{ client *clients.HTTPClient }
type KafkaBrokerDataSourceModel struct {
	ID        types.String `tfsdk:"id"`
	ClusterID types.String `tfsdk:"cluster_id"`
	BrokerID  types.Int64  `tfsdk:"broker_id"`
	Host      types.String `tfsdk:"host"`
	Port      types.Int64  `tfsdk:"port"`
	Rack      types.String `tfsdk:"rack"`
}

func NewKafkaBrokerDataSource() datasource.DataSource { return &KafkaBrokerDataSource{} }
func (d *KafkaBrokerDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_kafka_broker"
}
func (d *KafkaBrokerDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{MarkdownDescription: "Reads Kafka broker metadata from REST Proxy v3.", Attributes: map[string]schema.Attribute{"id": schema.StringAttribute{Computed: true}, "cluster_id": schema.StringAttribute{Required: true}, "broker_id": schema.Int64Attribute{Required: true}, "host": schema.StringAttribute{Computed: true}, "port": schema.Int64Attribute{Computed: true}, "rack": schema.StringAttribute{Computed: true}}}
}
func (d *KafkaBrokerDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, err := restProxyClient(req.ProviderData)
	if err != nil {
		resp.Diagnostics.AddError("Unable to configure Kafka broker data source", err.Error())
		return
	}
	d.client = client
}
func (d *KafkaBrokerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data KafkaBrokerDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	broker, err := d.client.GetBroker(ctx, data.ClusterID.ValueString(), data.BrokerID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Unable to read Kafka broker", err.Error())
		return
	}
	data.ID = types.StringValue(fmt.Sprintf("%s/%d", broker.ClusterID, broker.BrokerID))
	data.Host = types.StringValue(broker.Host)
	data.Port = types.Int64Value(broker.Port)
	data.Rack = types.StringValue(broker.Rack)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
