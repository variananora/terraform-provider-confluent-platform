// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/variananora/terraform-provider-confluent-platform/internal/clients"
)

var _ resource.Resource = &KafkaTopicResource{}
var _ resource.ResourceWithImportState = &KafkaTopicResource{}

type KafkaTopicResource struct{ client *clients.HTTPClient }

func NewKafkaTopicResource() resource.Resource { return &KafkaTopicResource{} }

func (r *KafkaTopicResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_kafka_topic"
}

func (r *KafkaTopicResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = KafkaTopicResourceSchema()
}

func (r *KafkaTopicResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	client, err := restProxyClient(req.ProviderData)
	if err != nil {
		resp.Diagnostics.AddError("Unable to configure Kafka topic resource", err.Error())
		return
	}
	r.client = client
}

func (r *KafkaTopicResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data KafkaTopicResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	configs := topicConfigs(ctx, data.Configs, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	topic, err := r.client.CreateTopic(ctx, data.ClusterID.ValueString(), clients.KafkaTopicCreateRequest{
		TopicName: data.Name.ValueString(), PartitionsCount: data.PartitionsCount.ValueInt64(), ReplicationFactor: data.ReplicationFactor.ValueInt64(), Configs: configs,
	})
	if err != nil {
		resp.Diagnostics.AddError("Unable to create Kafka topic", err.Error())
		return
	}
	setTopicResourceModel(ctx, &data, topic)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *KafkaTopicResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data KafkaTopicResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if (data.ClusterID.IsNull() || data.Name.IsNull()) && !data.ID.IsNull() {
		parts := strings.SplitN(data.ID.ValueString(), "/", 2)
		if len(parts) == 2 {
			data.ClusterID = types.StringValue(parts[0])
			data.Name = types.StringValue(parts[1])
		}
	}
	topic, err := r.client.GetTopic(ctx, data.ClusterID.ValueString(), data.Name.ValueString())
	if err != nil {
		if apiErr, ok := err.(*clients.APIError); ok && apiErr.NotFound() {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Unable to read Kafka topic", err.Error())
		return
	}
	setTopicResourceModel(ctx, &data, topic)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *KafkaTopicResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data KafkaTopicResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	configs := topicConfigs(ctx, data.Configs, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	var prior KafkaTopicResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &prior)...)
	if resp.Diagnostics.HasError() {
		return
	}
	priorConfigs := topicConfigValues(ctx, prior.Configs, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	configured := make(map[string]struct{}, len(configs))
	for _, config := range configs {
		configured[config.Name] = struct{}{}
		if err := r.client.UpdateTopicConfig(ctx, data.ClusterID.ValueString(), data.Name.ValueString(), config.Name, config.Value); err != nil {
			resp.Diagnostics.AddError("Unable to update Kafka topic configuration", err.Error())
			return
		}
	}
	for configName := range priorConfigs {
		if _, exists := configured[configName]; exists {
			continue
		}
		if err := r.client.DeleteTopicConfig(ctx, data.ClusterID.ValueString(), data.Name.ValueString(), configName); err != nil {
			resp.Diagnostics.AddError("Unable to delete Kafka topic configuration", err.Error())
			return
		}
	}
	topic, err := r.client.GetTopic(ctx, data.ClusterID.ValueString(), data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to read updated Kafka topic", err.Error())
		return
	}
	setTopicResourceModel(ctx, &data, topic)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *KafkaTopicResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data KafkaTopicResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.DeleteTopic(ctx, data.ClusterID.ValueString(), data.Name.ValueString()); err != nil {
		if apiErr, ok := err.(*clients.APIError); ok && apiErr.NotFound() {
			return
		}
		resp.Diagnostics.AddError("Unable to delete Kafka topic", err.Error())
	}
}

func (r *KafkaTopicResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if !strings.Contains(req.ID, "/") {
		resp.Diagnostics.AddError("Invalid Kafka topic import ID", "Use the format <cluster_id>/<topic_name>.")
		return
	}
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func topicConfigs(ctx context.Context, value types.Map, diagnostics interface{ AddError(string, string) }) []clients.KafkaTopicConfig {
	raw := topicConfigValues(ctx, value, diagnostics)
	configs := make([]clients.KafkaTopicConfig, 0, len(raw))
	for name, value := range raw {
		configs = append(configs, clients.KafkaTopicConfig{Name: name, Value: value})
	}
	return configs
}

func topicConfigValues(ctx context.Context, value types.Map, diagnostics interface{ AddError(string, string) }) map[string]string {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}
	var raw map[string]string
	if diags := value.ElementsAs(ctx, &raw, false); diags.HasError() {
		diagnostics.AddError("Invalid Kafka topic configs", fmt.Sprintf("Unable to decode topic configs: %s", diags.Errors()[0].Summary()))
		return nil
	}
	return raw
}

func setTopicResourceModel(ctx context.Context, data *KafkaTopicResourceModel, topic clients.KafkaTopic) {
	data.ID = types.StringValue(fmt.Sprintf("%s/%s", topic.ClusterID, topic.TopicName))
	data.ClusterID = types.StringValue(topic.ClusterID)
	data.Name = types.StringValue(topic.TopicName)
	data.PartitionsCount = types.Int64Value(topic.PartitionsCount)
	data.ReplicationFactor = types.Int64Value(topic.ReplicationFactor)
	data.IsInternal = types.BoolValue(topic.IsInternal)
	configValues := make(map[string]string, len(topic.Configs))
	for _, config := range topic.Configs {
		if !config.IsSensitive {
			configValues[config.Name] = config.Value
		}
	}
	data.Configs, _ = types.MapValueFrom(ctx, types.StringType, configValues)
}
