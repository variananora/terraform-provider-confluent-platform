// Copyright (c) HashiCorp, Inc.
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

var _ datasource.DataSource = &ConnectConnectorStatusDataSource{}

type ConnectConnectorStatusDataSource struct{ client *clients.HTTPClient }

type ConnectConnectorStatusDataSourceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	State        types.String `tfsdk:"state"`
	ResponseJSON types.String `tfsdk:"response_json"`
}

func NewConnectConnectorStatusDataSource() datasource.DataSource {
	return &ConnectConnectorStatusDataSource{}
}

func (d *ConnectConnectorStatusDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connect_connector_status"
}

func (d *ConnectConnectorStatusDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{MarkdownDescription: "Reads Kafka Connect connector status.", Attributes: map[string]schema.Attribute{
		"id":            schema.StringAttribute{Computed: true},
		"name":          schema.StringAttribute{Required: true},
		"state":         schema.StringAttribute{Computed: true},
		"response_json": schema.StringAttribute{Computed: true},
	}}
}

func (d *ConnectConnectorStatusDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, err := connectClient(req.ProviderData)
	if err != nil {
		resp.Diagnostics.AddError("Unable to configure Kafka Connect status data source", err.Error())
		return
	}
	d.client = client
}

func (d *ConnectConnectorStatusDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ConnectConnectorStatusDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	status, err := d.client.GetConnectorStatus(ctx, data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to read Kafka Connect connector status", err.Error())
		return
	}
	data.ID = types.StringValue(data.Name.ValueString())
	data.State = types.StringValue(status.Connector.State)
	data.ResponseJSON = types.StringValue(fmt.Sprintf(`{"name":%q,"state":%q}`, status.Name, status.Connector.State))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
