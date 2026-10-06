// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/variananora/terraform-provider-confluent-platform/internal/clients"
)

type KafkaReplicaDataSource struct{ client *clients.HTTPClient }
type KafkaReplicaDataSourceModel struct {
	ClusterID    types.String `tfsdk:"cluster_id"`
	Name         types.String `tfsdk:"name"`
	PartitionID  types.Int64  `tfsdk:"partition_id"`
	ID           types.String `tfsdk:"id"`
	ResponseJSON types.String `tfsdk:"response_json"`
}

func NewKafkaReplicaDataSource() datasource.DataSource { return &KafkaReplicaDataSource{} }
func (d *KafkaReplicaDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	rawMetadata("", "kafka_replica", req, resp)
}
func (d *KafkaReplicaDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = kafkaRawSchema("Reads Kafka replica metadata from REST Proxy v3.", true, true)
}
func (d *KafkaReplicaDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	configureKafkaRaw(&d.client, req, resp)
}
func (d *KafkaReplicaDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data KafkaReplicaDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	readKafkaRaw(ctx, d.client, req, resp, fmt.Sprintf("%s/replicas", rawPath(data.ClusterID.ValueString(), "topics", data.Name.ValueString(), "partitions", fmt.Sprint(data.PartitionID.ValueInt64()))), "Kafka replica")
}
