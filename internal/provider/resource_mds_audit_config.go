// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/variananora/terraform-provider-confluent-platform/internal/clients"
)

var _ resource.Resource = &MDSAuditConfigResource{}

type MDSAuditConfigResource struct{ client mdsAuditClient }

type MDSAuditConfigResourceModel struct {
	ID                types.String `tfsdk:"id"`
	ConfigurationJSON types.String `tfsdk:"configuration_json"`
	ResponseJSON      types.String `tfsdk:"response_json"`
	ResourceVersion   types.String `tfsdk:"resource_version"`
}

type mdsAuditClient interface {
	GetMDSAuditConfig(context.Context) ([]byte, error)
	PutMDSAuditConfig(context.Context, any) ([]byte, error)
}

func NewMDSAuditConfigResource() resource.Resource { return &MDSAuditConfigResource{} }

func (r *MDSAuditConfigResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mds_audit_config"
}

func (r *MDSAuditConfigResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{MarkdownDescription: "Manages the MDS audit configuration. configuration_json must include the current metadata.resource_version.", Attributes: map[string]schema.Attribute{
		"id":                 schema.StringAttribute{Computed: true},
		"configuration_json": schema.StringAttribute{Required: true},
		"response_json":      schema.StringAttribute{Computed: true},
		"resource_version":   schema.StringAttribute{Computed: true},
	}}
}

func (r *MDSAuditConfigResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	client, err := mdsClient(req.ProviderData)
	if err != nil {
		resp.Diagnostics.AddError("Unable to configure MDS audit resource", err.Error())
		return
	}
	r.client = client
}

func (r *MDSAuditConfigResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data MDSAuditConfigResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	configuration, err := decodeJSON(data.ConfigurationJSON.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid MDS audit configuration", err.Error())
		return
	}
	payload, err := r.client.PutMDSAuditConfig(ctx, configuration)
	if err != nil {
		resp.Diagnostics.AddError("Unable to update MDS audit configuration", auditError(err))
		return
	}
	data.ID = types.StringValue("audit-config")
	data.ResponseJSON = types.StringValue(string(payload))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *MDSAuditConfigResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data MDSAuditConfigResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	payload, err := r.client.GetMDSAuditConfig(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Unable to read MDS audit configuration", auditError(err))
		return
	}
	data.ID = types.StringValue("audit-config")
	data.ResponseJSON = types.StringValue(string(payload))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *MDSAuditConfigResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data MDSAuditConfigResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	configuration, err := decodeJSON(data.ConfigurationJSON.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid MDS audit configuration", err.Error())
		return
	}
	payload, err := r.client.PutMDSAuditConfig(ctx, configuration)
	if err != nil {
		resp.Diagnostics.AddError("Unable to update MDS audit configuration", auditError(err))
		return
	}
	data.ID = types.StringValue("audit-config")
	data.ResponseJSON = types.StringValue(string(payload))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *MDSAuditConfigResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}

func auditError(err error) string {
	if apiErr, ok := err.(*clients.APIError); ok && apiErr.StatusCode == 409 {
		return "MDS audit configuration changed concurrently; refresh the current configuration and preserve its metadata.resource_version before retrying: " + apiErr.Error()
	}
	return err.Error()
}
