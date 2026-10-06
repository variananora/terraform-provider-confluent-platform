// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type KafkaTopicResourceModel struct {
	ID                types.String `tfsdk:"id"`
	ClusterID         types.String `tfsdk:"cluster_id"`
	Name              types.String `tfsdk:"name"`
	PartitionsCount   types.Int64  `tfsdk:"partitions_count"`
	ReplicationFactor types.Int64  `tfsdk:"replication_factor"`
	Configs           types.Map    `tfsdk:"configs"`
	IsInternal        types.Bool   `tfsdk:"is_internal"`
}

func KafkaTopicResourceSchema() schema.Schema {
	return schema.Schema{MarkdownDescription: "Manages a Kafka topic through the Confluent Platform REST Proxy v3 API.", Attributes: map[string]schema.Attribute{"id": schema.StringAttribute{Computed: true}, "cluster_id": schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}}, "name": schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}}, "partitions_count": schema.Int64Attribute{Optional: true, Computed: true, PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()}}, "replication_factor": schema.Int64Attribute{Optional: true, Computed: true, PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()}}, "configs": schema.MapAttribute{Optional: true, Computed: true, ElementType: types.StringType}, "is_internal": schema.BoolAttribute{Computed: true}}}
}
