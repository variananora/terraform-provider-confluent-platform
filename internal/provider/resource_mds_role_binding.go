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

var _ resource.Resource = &MDSRoleBindingResource{}
var _ resource.ResourceWithImportState = &MDSRoleBindingResource{}

type MDSRoleBindingResource struct{ client *clients.HTTPClient }

type MDSRoleBindingResourceModel struct {
	ID             types.String `tfsdk:"id"`
	Principal      types.String `tfsdk:"principal"`
	RoleName       types.String `tfsdk:"role_name"`
	KafkaClusterID types.String `tfsdk:"kafka_cluster_id"`
	ResourceType   types.String `tfsdk:"resource_type"`
	ResourceName   types.String `tfsdk:"resource_name"`
	PatternType    types.String `tfsdk:"pattern_type"`
}

func NewMDSRoleBindingResource() resource.Resource { return &MDSRoleBindingResource{} }

func (r *MDSRoleBindingResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mds_role_binding"
}

func (r *MDSRoleBindingResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{MarkdownDescription: "Manages an MDS RBAC role binding for a Kafka cluster, topic, or consumer group.", Attributes: map[string]schema.Attribute{
		"id":               schema.StringAttribute{Computed: true},
		"principal":        replaceString("Kafka principal, such as User:alice or Group:analysts."),
		"role_name":        replaceString("Confluent RBAC role name, such as DeveloperRead or DeveloperManage."),
		"kafka_cluster_id": replaceString("Kafka cluster ID used in the MDS scope."),
		"resource_type":    optionalReplaceString("Resource type: Topic or Group. Omit for a cluster-scoped role."),
		"resource_name":    optionalReplaceString("Topic or consumer group name. Omit for a cluster-scoped role."),
		"pattern_type":     optionalReplaceString("Resource pattern type, usually LITERAL or PREFIXED."),
	}}
}

func replaceString(description string) schema.StringAttribute {
	return schema.StringAttribute{Required: true, MarkdownDescription: description, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}}
}

func optionalReplaceString(description string) schema.StringAttribute {
	return schema.StringAttribute{Optional: true, MarkdownDescription: description, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}}
}

func (r *MDSRoleBindingResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	client, err := mdsClient(req.ProviderData)
	if err != nil {
		resp.Diagnostics.AddError("Unable to configure MDS role binding resource", err.Error())
		return
	}
	r.client = client
}

func (r *MDSRoleBindingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data MDSRoleBindingResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	request, valid := roleBindingRequest(data, &resp.Diagnostics)
	if !valid {
		return
	}
	if err := r.client.CreateMDSRoleBinding(ctx, data.Principal.ValueString(), data.RoleName.ValueString(), request); err != nil {
		resp.Diagnostics.AddError("Unable to create MDS role binding", mdsErrorMessage(err))
		return
	}
	setRoleBindingState(&data)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *MDSRoleBindingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data MDSRoleBindingResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	request, valid := roleBindingRequest(data, &resp.Diagnostics)
	if !valid {
		return
	}
	resources, err := r.client.ReadMDSRoleBindingResources(ctx, data.Principal.ValueString(), data.RoleName.ValueString(), request.Scope)
	if err != nil {
		if apiErr, ok := err.(*clients.APIError); ok && apiErr.NotFound() {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Unable to read MDS role binding", mdsErrorMessage(err))
		return
	}
	if len(request.ResourcePatterns) > 0 && !containsResourcePattern(resources, request.ResourcePatterns[0]) {
		resp.State.RemoveResource(ctx)
		return
	}
	setRoleBindingState(&data)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *MDSRoleBindingResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
}

func (r *MDSRoleBindingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data MDSRoleBindingResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	request, valid := roleBindingRequest(data, &resp.Diagnostics)
	if !valid {
		return
	}
	if err := r.client.DeleteMDSRoleBinding(ctx, data.Principal.ValueString(), data.RoleName.ValueString(), request); err != nil {
		if apiErr, ok := err.(*clients.APIError); ok && apiErr.NotFound() {
			return
		}
		resp.Diagnostics.AddError("Unable to delete MDS role binding", mdsErrorMessage(err))
	}
}

func (r *MDSRoleBindingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "|")
	if len(parts) != 3 && len(parts) != 6 {
		resp.Diagnostics.AddError("Invalid MDS role binding import ID", "Use principal|role_name|kafka_cluster_id or principal|role_name|kafka_cluster_id|resource_type|resource_name|pattern_type.")
		return
	}
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	attributes := []string{"principal", "role_name", "kafka_cluster_id"}
	if len(parts) == 6 {
		attributes = append(attributes, "resource_type", "resource_name", "pattern_type")
	}
	for index, attribute := range attributes {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(attribute), types.StringValue(parts[index]))...)
	}
}

func roleBindingRequest(data MDSRoleBindingResourceModel, diagnostics interface{ AddError(string, string) }) (clients.MDSRoleBindingRequest, bool) {
	request := clients.MDSRoleBindingRequest{Scope: clients.MDSRoleBindingScope{Clusters: map[string]string{"kafka-cluster": data.KafkaClusterID.ValueString()}}}
	hasResource := !data.ResourceType.IsNull() || !data.ResourceName.IsNull() || !data.PatternType.IsNull()
	if !hasResource {
		return request, true
	}
	if data.ResourceType.IsNull() || data.ResourceName.IsNull() || data.PatternType.IsNull() {
		diagnostics.AddError("Invalid MDS role binding scope", "resource_type, resource_name, and pattern_type must be provided together.")
		return request, false
	}
	resourceType := data.ResourceType.ValueString()
	if resourceType != "Topic" && resourceType != "Group" {
		diagnostics.AddError("Invalid MDS role binding resource type", "resource_type must be Topic or Group; cluster-scoped roles omit resource fields.")
		return request, false
	}
	request.ResourcePatterns = []clients.MDSResourcePattern{{ResourceType: resourceType, Name: data.ResourceName.ValueString(), PatternType: data.PatternType.ValueString()}}
	return request, true
}

func containsResourcePattern(resources []clients.MDSResourcePattern, expected clients.MDSResourcePattern) bool {
	for _, resource := range resources {
		if resource == expected {
			return true
		}
	}
	return false
}

func setRoleBindingState(data *MDSRoleBindingResourceModel) {
	data.ID = types.StringValue(fmt.Sprintf("%s|%s|%s", data.Principal.ValueString(), data.RoleName.ValueString(), data.KafkaClusterID.ValueString()))
	if !data.ResourceType.IsNull() {
		data.ID = types.StringValue(fmt.Sprintf("%s|%s|%s|%s|%s|%s", data.Principal.ValueString(), data.RoleName.ValueString(), data.KafkaClusterID.ValueString(), data.ResourceType.ValueString(), data.ResourceName.ValueString(), data.PatternType.ValueString()))
	}
}

func mdsClient(providerData any) (*clients.HTTPClient, error) {
	runtime, ok := providerData.(*clients.Runtime)
	if !ok {
		return nil, fmt.Errorf("expected *clients.Runtime, got %T", providerData)
	}
	client, ok := runtime.Service("mds")
	if !ok {
		return nil, fmt.Errorf("MDS endpoint is not configured")
	}
	return client, nil
}
