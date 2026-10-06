// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
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

var _ resource.Resource = &MDSACLResource{}
var _ resource.ResourceWithImportState = &MDSACLResource{}

type MDSACLResource struct{ client *clients.HTTPClient }

type MDSACLResourceModel struct {
	ID             types.String `tfsdk:"id"`
	KafkaClusterID types.String `tfsdk:"kafka_cluster_id"`
	ResourceType   types.String `tfsdk:"resource_type"`
	ResourceName   types.String `tfsdk:"resource_name"`
	PatternType    types.String `tfsdk:"pattern_type"`
	Principal      types.String `tfsdk:"principal"`
	Host           types.String `tfsdk:"host"`
	Operation      types.String `tfsdk:"operation"`
	PermissionType types.String `tfsdk:"permission_type"`
}

func NewMDSACLResource() resource.Resource { return &MDSACLResource{} }

func (r *MDSACLResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mds_acl"
}

func (r *MDSACLResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{MarkdownDescription: "Manages a centralized Kafka ACL through MDS.", Attributes: map[string]schema.Attribute{
		"id":               schema.StringAttribute{Computed: true},
		"kafka_cluster_id": mdsACLReplaceString("Kafka cluster ID used in the MDS scope."),
		"resource_type":    mdsACLReplaceString("ACL resource type, such as Topic or Group."),
		"resource_name":    mdsACLReplaceString("ACL resource name."),
		"pattern_type":     mdsACLReplaceString("ACL pattern type, such as LITERAL."),
		"principal":        mdsACLReplaceString("Kafka principal, such as User:alice."),
		"host":             mdsACLReplaceString("ACL host, commonly *."),
		"operation":        mdsACLReplaceString("Kafka operation, such as Read or Write."),
		"permission_type":  mdsACLReplaceString("ACL permission, such as Allow or Deny."),
	}}
}

func mdsACLReplaceString(description string) schema.StringAttribute {
	return schema.StringAttribute{Required: true, MarkdownDescription: description, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}}
}

func (r *MDSACLResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	client, err := mdsClient(req.ProviderData)
	if err != nil {
		resp.Diagnostics.AddError("Unable to configure MDS ACL resource", err.Error())
		return
	}
	r.client = client
}

func (r *MDSACLResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data MDSACLResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.CreateMDSACL(ctx, mdsACLRequest(data)); err != nil {
		resp.Diagnostics.AddError("Unable to create MDS ACL", mdsErrorMessage(err))
		return
	}
	setMDSACLState(&data)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *MDSACLResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data MDSACLResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	bindings, err := r.client.SearchMDSACLs(ctx, mdsACLFilter(data))
	if err != nil {
		if apiErr, ok := err.(*clients.APIError); ok && apiErr.NotFound() {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Unable to read MDS ACL", mdsErrorMessage(err))
		return
	}
	if !containsMDSACL(bindings, mdsACLBinding(data)) {
		resp.State.RemoveResource(ctx)
		return
	}
	setMDSACLState(&data)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *MDSACLResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
}

func (r *MDSACLResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data MDSACLResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.DeleteMDSACL(ctx, mdsACLFilter(data)); err != nil {
		if apiErr, ok := err.(*clients.APIError); ok && apiErr.NotFound() {
			return
		}
		resp.Diagnostics.AddError("Unable to delete MDS ACL", mdsErrorMessage(err))
	}
}

func (r *MDSACLResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "|")
	if len(parts) != 8 {
		resp.Diagnostics.AddError("Invalid MDS ACL import ID", "Use kafka_cluster_id|resource_type|resource_name|pattern_type|principal|host|operation|permission_type.")
		return
	}
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	for index, attribute := range []string{"kafka_cluster_id", "resource_type", "resource_name", "pattern_type", "principal", "host", "operation", "permission_type"} {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(attribute), types.StringValue(parts[index]))...)
	}
}

func mdsACLRequest(data MDSACLResourceModel) clients.MDSACLRequest {
	return clients.MDSACLRequest{Scope: clients.MDSRoleBindingScope{Clusters: map[string]string{"kafka-cluster": data.KafkaClusterID.ValueString()}}, ACLBinding: mdsACLBinding(data)}
}

func mdsACLFilter(data MDSACLResourceModel) clients.MDSACLFilter {
	binding := mdsACLBinding(data)
	return clients.MDSACLFilter{Scope: clients.MDSRoleBindingScope{Clusters: map[string]string{"kafka-cluster": data.KafkaClusterID.ValueString()}}, ACLBindingFilter: clients.MDSACLBindingFilter{PatternFilter: binding.Pattern, EntryFilter: binding.Entry}}
}

func mdsACLBinding(data MDSACLResourceModel) clients.MDSACLBinding {
	return clients.MDSACLBinding{Pattern: clients.MDSACLPattern{ResourceType: data.ResourceType.ValueString(), Name: data.ResourceName.ValueString(), PatternType: data.PatternType.ValueString()}, Entry: clients.MDSACLEntry{Principal: data.Principal.ValueString(), Host: data.Host.ValueString(), Operation: data.Operation.ValueString(), PermissionType: data.PermissionType.ValueString()}}
}

func containsMDSACL(bindings []clients.MDSACLBinding, expected clients.MDSACLBinding) bool {
	for _, binding := range bindings {
		if binding == expected {
			return true
		}
	}
	return false
}

func setMDSACLState(data *MDSACLResourceModel) {
	data.ID = types.StringValue(fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s", data.KafkaClusterID.ValueString(), data.ResourceType.ValueString(), data.ResourceName.ValueString(), data.PatternType.ValueString(), data.Principal.ValueString(), data.Host.ValueString(), data.Operation.ValueString(), data.PermissionType.ValueString()))
}
