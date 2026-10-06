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

var _ resource.Resource = &KafkaACLResource{}
var _ resource.ResourceWithImportState = &KafkaACLResource{}

type KafkaACLResource struct{ client *clients.HTTPClient }

type KafkaACLResourceModel struct {
	ID           types.String `tfsdk:"id"`
	ClusterID    types.String `tfsdk:"cluster_id"`
	ResourceType types.String `tfsdk:"resource_type"`
	ResourceName types.String `tfsdk:"resource_name"`
	PatternType  types.String `tfsdk:"pattern_type"`
	Principal    types.String `tfsdk:"principal"`
	Host         types.String `tfsdk:"host"`
	Operation    types.String `tfsdk:"operation"`
	Permission   types.String `tfsdk:"permission"`
}

func NewKafkaACLResource() resource.Resource { return &KafkaACLResource{} }

func (r *KafkaACLResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_kafka_acl"
}

func (r *KafkaACLResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{MarkdownDescription: "Manages a Kafka ACL through REST Proxy v3.", Attributes: map[string]schema.Attribute{
		"id":            schema.StringAttribute{Computed: true},
		"cluster_id":    aclRequiredReplaceAttribute("Kafka cluster ID."),
		"resource_type": aclRequiredReplaceAttribute("ACL resource type."),
		"resource_name": aclRequiredReplaceAttribute("ACL resource name."),
		"pattern_type":  aclRequiredReplaceAttribute("ACL pattern type."),
		"principal":     aclRequiredReplaceAttribute("ACL principal."),
		"host":          aclRequiredReplaceAttribute("ACL host."),
		"operation":     aclRequiredReplaceAttribute("ACL operation."),
		"permission":    aclRequiredReplaceAttribute("ACL permission."),
	}}
}

func aclRequiredReplaceAttribute(description string) schema.StringAttribute {
	return schema.StringAttribute{Required: true, MarkdownDescription: description, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}}
}

func (r *KafkaACLResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	client, err := restProxyClient(req.ProviderData)
	if err != nil {
		resp.Diagnostics.AddError("Unable to configure Kafka ACL resource", err.Error())
		return
	}
	r.client = client
}

func (r *KafkaACLResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data KafkaACLResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	acl, err := r.client.CreateKafkaACL(ctx, data.ClusterID.ValueString(), kafkaACLRequest(data))
	if err != nil {
		resp.Diagnostics.AddError("Unable to create Kafka ACL", err.Error())
		return
	}
	setKafkaACLModel(&data, acl)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *KafkaACLResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data KafkaACLResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	acl, err := r.client.GetKafkaACL(ctx, data.ClusterID.ValueString(), kafkaACL(data))
	if err != nil {
		if apiErr, ok := err.(*clients.APIError); ok && apiErr.NotFound() {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Unable to read Kafka ACL", err.Error())
		return
	}
	setKafkaACLModel(&data, acl)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *KafkaACLResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
}

func (r *KafkaACLResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data KafkaACLResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.DeleteKafkaACL(ctx, data.ClusterID.ValueString(), kafkaACL(data)); err != nil {
		if apiErr, ok := err.(*clients.APIError); ok && apiErr.NotFound() {
			return
		}
		resp.Diagnostics.AddError("Unable to delete Kafka ACL", err.Error())
	}
}

func (r *KafkaACLResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "|")
	if len(parts) != 8 {
		resp.Diagnostics.AddError("Invalid Kafka ACL import ID", "Use cluster_id|resource_type|resource_name|pattern_type|principal|host|operation|permission.")
		return
	}
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	for index, attribute := range []string{"cluster_id", "resource_type", "resource_name", "pattern_type", "principal", "host", "operation", "permission"} {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(attribute), types.StringValue(parts[index]))...)
	}
}

func kafkaACLRequest(data KafkaACLResourceModel) clients.KafkaACLCreateRequest {
	return clients.KafkaACLCreateRequest{ResourceType: data.ResourceType.ValueString(), ResourceName: data.ResourceName.ValueString(), PatternType: data.PatternType.ValueString(), Principal: data.Principal.ValueString(), Host: data.Host.ValueString(), Operation: data.Operation.ValueString(), Permission: data.Permission.ValueString()}
}

func kafkaACL(data KafkaACLResourceModel) clients.KafkaACL {
	return clients.KafkaACL{ClusterID: data.ClusterID.ValueString(), ResourceType: data.ResourceType.ValueString(), ResourceName: data.ResourceName.ValueString(), PatternType: data.PatternType.ValueString(), Principal: data.Principal.ValueString(), Host: data.Host.ValueString(), Operation: data.Operation.ValueString(), Permission: data.Permission.ValueString()}
}

func setKafkaACLModel(data *KafkaACLResourceModel, acl clients.KafkaACL) {
	data.ID = types.StringValue(fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s", acl.ClusterID, acl.ResourceType, acl.ResourceName, acl.PatternType, acl.Principal, acl.Host, acl.Operation, acl.Permission))
	data.ClusterID = types.StringValue(acl.ClusterID)
	data.ResourceType = types.StringValue(acl.ResourceType)
	data.ResourceName = types.StringValue(acl.ResourceName)
	data.PatternType = types.StringValue(acl.PatternType)
	data.Principal = types.StringValue(acl.Principal)
	data.Host = types.StringValue(acl.Host)
	data.Operation = types.StringValue(acl.Operation)
	data.Permission = types.StringValue(acl.Permission)
}
