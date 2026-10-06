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

var _ datasource.DataSource = &MDSActiveNodesDataSource{}

type MDSActiveNodesDataSource struct{ client *clients.HTTPClient }

type MDSActiveNodesDataSourceModel struct {
	ID           types.String `tfsdk:"id"`
	Protocol     types.String `tfsdk:"protocol"`
	ResponseJSON types.String `tfsdk:"response_json"`
}

func NewMDSActiveNodesDataSource() datasource.DataSource { return &MDSActiveNodesDataSource{} }

func (d *MDSActiveNodesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mds_active_nodes"
}

func (d *MDSActiveNodesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{MarkdownDescription: "Reads active Metadata Service nodes.", Attributes: map[string]schema.Attribute{
		"id":            schema.StringAttribute{Computed: true},
		"protocol":      schema.StringAttribute{Required: true, MarkdownDescription: "MDS transport protocol, such as https."},
		"response_json": schema.StringAttribute{Computed: true, MarkdownDescription: "Raw active-node response from MDS."},
	}}
}

func (d *MDSActiveNodesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	runtime, ok := req.ProviderData.(*clients.Runtime)
	if !ok {
		resp.Diagnostics.AddError("Unable to configure MDS data source", fmt.Sprintf("expected *clients.Runtime, got %T", req.ProviderData))
		return
	}
	client, ok := runtime.Service("mds")
	if !ok {
		resp.Diagnostics.AddError("Unable to configure MDS data source", "MDS endpoint is not configured")
		return
	}
	d.client = client
}

func (d *MDSActiveNodesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data MDSActiveNodesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	payload, err := d.client.GetMDSActiveNodes(ctx, data.Protocol.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to read MDS active nodes", err.Error())
		return
	}
	data.ID = types.StringValue(fmt.Sprintf("active-nodes/%s", data.Protocol.ValueString()))
	data.ResponseJSON = types.StringValue(string(payload))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
