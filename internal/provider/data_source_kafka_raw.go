// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/variananora/terraform-provider-confluent-platform/internal/clients"
)

type kafkaRawDataSourceModel struct {
	ID           types.String `tfsdk:"id"`
	ClusterID    types.String `tfsdk:"cluster_id"`
	Name         types.String `tfsdk:"name"`
	Topic        types.String `tfsdk:"topic"`
	PartitionID  types.Int64  `tfsdk:"partition_id"`
	ResponseJSON types.String `tfsdk:"response_json"`
}

func kafkaRawSchema(description string, topic, partition bool) schema.Schema {
	attrs := map[string]schema.Attribute{"id": schema.StringAttribute{Computed: true}, "cluster_id": schema.StringAttribute{Required: true}, "name": schema.StringAttribute{Required: true}, "response_json": schema.StringAttribute{Computed: true}}
	if topic {
		attrs["topic"] = schema.StringAttribute{Required: true}
	}
	if partition {
		attrs["partition_id"] = schema.Int64Attribute{Required: true}
	}
	return schema.Schema{MarkdownDescription: description, Attributes: attrs}
}
func configureKafkaRaw(client **clients.HTTPClient, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	configured, err := restProxyClient(req.ProviderData)
	if err != nil {
		resp.Diagnostics.AddError("Unable to configure Kafka data source", err.Error())
		return
	}
	*client = configured
}
func readKafkaRaw(ctx context.Context, client *clients.HTTPClient, req datasource.ReadRequest, resp *datasource.ReadResponse, path, dataSourceName string) {
	var data kafkaRawDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	payload, err := client.GetRawJSON(ctx, path)
	if err != nil {
		resp.Diagnostics.AddError("Unable to read "+dataSourceName, err.Error())
		return
	}
	data.ID = types.StringValue(path)
	data.ResponseJSON = types.StringValue(string(payload))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
func rawMetadata(providerType, name string, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + name
}
func rawPath(cluster string, segments ...string) string {
	path := "/kafka/v3/clusters/" + clients.EscapePathSegment(cluster)
	for _, segment := range segments {
		path += "/" + clients.EscapePathSegment(segment)
	}
	return path
}
