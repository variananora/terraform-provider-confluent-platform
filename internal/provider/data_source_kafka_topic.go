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

type KafkaTopicDataSource struct{ client *clients.HTTPClient }
type KafkaTopicDataSourceModel struct {
	ID                types.String `tfsdk:"id"`
	ClusterID         types.String `tfsdk:"cluster_id"`
	Name              types.String `tfsdk:"name"`
	PartitionsCount   types.Int64  `tfsdk:"partitions_count"`
	ReplicationFactor types.Int64  `tfsdk:"replication_factor"`
	Configs           types.Map    `tfsdk:"configs"`
	IsInternal        types.Bool   `tfsdk:"is_internal"`
}

func NewKafkaTopicDataSource() datasource.DataSource { return &KafkaTopicDataSource{} }
func (d *KafkaTopicDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_kafka_topic"
}
func (d *KafkaTopicDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{MarkdownDescription: "Reads Kafka topic metadata from REST Proxy v3.", Attributes: map[string]schema.Attribute{"id": schema.StringAttribute{Computed: true}, "cluster_id": schema.StringAttribute{Required: true}, "name": schema.StringAttribute{Required: true}, "partitions_count": schema.Int64Attribute{Computed: true}, "replication_factor": schema.Int64Attribute{Computed: true}, "configs": schema.MapAttribute{Computed: true, ElementType: types.StringType}, "is_internal": schema.BoolAttribute{Computed: true}}}
}
func (d *KafkaTopicDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, err := restProxyClient(req.ProviderData)
	if err != nil {
		resp.Diagnostics.AddError("Unable to configure Kafka topic data source", err.Error())
		return
	}
	d.client = client
}
func (d *KafkaTopicDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data KafkaTopicDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	topic, err := d.client.GetTopic(ctx, data.ClusterID.ValueString(), data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to read Kafka topic", err.Error())
		return
	}
	data.ID = types.StringValue(fmt.Sprintf("%s/%s", topic.ClusterID, topic.TopicName))
	data.PartitionsCount = types.Int64Value(topic.PartitionsCount)
	data.ReplicationFactor = types.Int64Value(topic.ReplicationFactor)
	data.IsInternal = types.BoolValue(topic.IsInternal)
	configs := make(map[string]string, len(topic.Configs))
	for _, config := range topic.Configs {
		configs[config.Name] = config.Value
	}
	data.Configs, _ = types.MapValueFrom(ctx, types.StringType, configs)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
