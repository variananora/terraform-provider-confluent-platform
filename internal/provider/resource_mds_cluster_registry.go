// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/variananora/terraform-provider-confluent-platform/internal/clients"
)

var _ resource.Resource = &MDSClusterRegistryResource{}
var _ resource.ResourceWithImportState = &MDSClusterRegistryResource{}

type MDSClusterRegistryResource struct{ client *clients.HTTPClient }

type MDSClusterRegistryResourceModel struct {
	ID                types.String `tfsdk:"id"`
	ClusterName       types.String `tfsdk:"cluster_name"`
	ConfigurationJSON types.String `tfsdk:"configuration_json"`
	ResponseJSON      types.String `tfsdk:"response_json"`
}

func NewMDSClusterRegistryResource() resource.Resource { return &MDSClusterRegistryResource{} }

func (r *MDSClusterRegistryResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mds_cluster_registry"
}

func (r *MDSClusterRegistryResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{MarkdownDescription: "Manages a named MDS cluster registry entry.", Attributes: map[string]schema.Attribute{
		"id":                 schema.StringAttribute{Computed: true},
		"cluster_name":       schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
		"configuration_json": schema.StringAttribute{Required: true, MarkdownDescription: "JSON cluster registry object, including scope, hosts, and protocol."},
		"response_json":      schema.StringAttribute{Computed: true},
	}}
}

func (r *MDSClusterRegistryResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	client, err := mdsClient(req.ProviderData)
	if err != nil {
		resp.Diagnostics.AddError("Unable to configure MDS cluster registry resource", err.Error())
		return
	}
	r.client = client
}

func (r *MDSClusterRegistryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data MDSClusterRegistryResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	configuration, err := decodeJSON(data.ConfigurationJSON.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid MDS cluster registry configuration", err.Error())
		return
	}
	if err := r.client.CreateMDSClusterRegistry(ctx, configuration); err != nil {
		resp.Diagnostics.AddError("Unable to create MDS cluster registry entry", mdsErrorMessage(err))
		return
	}
	r.refresh(ctx, &data, resp)
}

func (r *MDSClusterRegistryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data MDSClusterRegistryResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	payload, err := r.client.GetMDSClusterRegistry(ctx, data.ClusterName.ValueString())
	if err != nil {
		if apiErr, ok := err.(*clients.APIError); ok && apiErr.NotFound() {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Unable to read MDS cluster registry entry", mdsErrorMessage(err))
		return
	}
	data.ID = types.StringValue(data.ClusterName.ValueString())
	data.ResponseJSON = types.StringValue(string(payload))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *MDSClusterRegistryResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
}

func (r *MDSClusterRegistryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data MDSClusterRegistryResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.DeleteMDSClusterRegistry(ctx, data.ClusterName.ValueString()); err != nil {
		if apiErr, ok := err.(*clients.APIError); ok && apiErr.NotFound() {
			return
		}
		resp.Diagnostics.AddError("Unable to delete MDS cluster registry entry", mdsErrorMessage(err))
	}
}

func (r *MDSClusterRegistryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if strings.TrimSpace(req.ID) == "" {
		resp.Diagnostics.AddError("Invalid MDS cluster registry import ID", "Use the registered cluster name.")
		return
	}
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("cluster_name"), types.StringValue(req.ID))...)
}

func (r *MDSClusterRegistryResource) refresh(ctx context.Context, data *MDSClusterRegistryResourceModel, resp *resource.CreateResponse) {
	payload, err := r.client.GetMDSClusterRegistry(ctx, data.ClusterName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to read created MDS cluster registry entry", mdsErrorMessage(err))
		return
	}
	data.ID = types.StringValue(data.ClusterName.ValueString())
	data.ResponseJSON = types.StringValue(string(payload))
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func decodeJSON(value string) (any, error) {
	var result any
	if err := json.Unmarshal([]byte(value), &result); err != nil {
		return nil, fmt.Errorf("configuration_json must contain valid JSON: %w", err)
	}
	return result, nil
}
