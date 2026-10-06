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

var _ datasource.DataSource = &MDSMetadataClusterIDDataSource{}

type MDSMetadataClusterIDDataSource struct{ client *clients.HTTPClient }

type MDSMetadataClusterIDDataSourceModel struct {
	ID           types.String `tfsdk:"id"`
	ResponseJSON types.String `tfsdk:"response_json"`
}

func NewMDSMetadataClusterIDDataSource() datasource.DataSource {
	return &MDSMetadataClusterIDDataSource{}
}

func (d *MDSMetadataClusterIDDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mds_metadata_cluster_id"
}

func (d *MDSMetadataClusterIDDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{MarkdownDescription: "Reads the Metadata Service cluster ID.", Attributes: map[string]schema.Attribute{
		"id":            schema.StringAttribute{Computed: true},
		"response_json": schema.StringAttribute{Computed: true, MarkdownDescription: "Raw metadata cluster ID response from MDS."},
	}}
}

func (d *MDSMetadataClusterIDDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *MDSMetadataClusterIDDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	payload, err := d.client.GetMDSMetadataClusterID(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Unable to read MDS metadata cluster ID", err.Error())
		return
	}
	data := MDSMetadataClusterIDDataSourceModel{ID: types.StringValue("metadata-cluster-id"), ResponseJSON: types.StringValue(string(payload))}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
