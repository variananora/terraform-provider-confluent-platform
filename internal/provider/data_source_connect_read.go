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

type ConnectReadDataSource struct {
	kind   string
	client *clients.HTTPClient
}

type ConnectReadDataSourceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	ResponseJSON types.String `tfsdk:"response_json"`
}

func NewConnectClusterDataSource() datasource.DataSource {
	return &ConnectReadDataSource{kind: "cluster"}
}
func NewConnectConnectorInfoDataSource() datasource.DataSource {
	return &ConnectReadDataSource{kind: "connector_info"}
}
func NewConnectConnectorConfigDataSource() datasource.DataSource {
	return &ConnectReadDataSource{kind: "connector_config"}
}

func (d *ConnectReadDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connect_" + d.kind
}

func (d *ConnectReadDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	attrs := map[string]schema.Attribute{
		"id":            schema.StringAttribute{Computed: true},
		"response_json": schema.StringAttribute{Computed: true},
	}
	if d.kind != "cluster" {
		attrs["name"] = schema.StringAttribute{Required: true}
	}
	resp.Schema = schema.Schema{MarkdownDescription: "Reads Kafka Connect " + d.kind + " information as JSON.", Attributes: attrs}
}

func (d *ConnectReadDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, err := connectClient(req.ProviderData)
	if err != nil {
		resp.Diagnostics.AddError("Unable to configure Kafka Connect data source", err.Error())
		return
	}
	d.client = client
}

func (d *ConnectReadDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ConnectReadDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	path := "/"
	if d.kind == "connector_info" {
		path = "/connectors/" + clients.EscapePathSegment(data.Name.ValueString())
	} else if d.kind == "connector_config" {
		path = "/connectors/" + clients.EscapePathSegment(data.Name.ValueString()) + "/config"
	}
	payload, err := d.client.GetRawJSON(ctx, path)
	if err != nil {
		resp.Diagnostics.AddError("Unable to read Kafka Connect "+d.kind, err.Error())
		return
	}
	data.ID = types.StringValue(fmt.Sprintf("%s%s", d.kind, path))
	data.ResponseJSON = types.StringValue(string(payload))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
