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

var _ resource.Resource = &ConnectConnectorResource{}
var _ resource.ResourceWithImportState = &ConnectConnectorResource{}

type ConnectConnectorResource struct{ client *clients.HTTPClient }

type ConnectConnectorResourceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Config       types.Map    `tfsdk:"config"`
	InitialState types.String `tfsdk:"initial_state"`
	Type         types.String `tfsdk:"type"`
	State        types.String `tfsdk:"state"`
}

func NewConnectConnectorResource() resource.Resource { return &ConnectConnectorResource{} }

func (r *ConnectConnectorResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connect_connector"
}

func (r *ConnectConnectorResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{MarkdownDescription: "Manages a Kafka Connect connector.", Attributes: map[string]schema.Attribute{
		"id":            schema.StringAttribute{Computed: true},
		"name":          schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
		"config":        schema.MapAttribute{Required: true, Sensitive: true, ElementType: types.StringType},
		"initial_state": schema.StringAttribute{Optional: true},
		"type":          schema.StringAttribute{Computed: true},
		"state":         schema.StringAttribute{Computed: true},
	}}
}

func (r *ConnectConnectorResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	client, err := connectClient(req.ProviderData)
	if err != nil {
		resp.Diagnostics.AddError("Unable to configure Kafka Connect connector resource", err.Error())
		return
	}
	r.client = client
}

func (r *ConnectConnectorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ConnectConnectorResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	config := connectorConfig(ctx, data.Config, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	connector, err := r.client.CreateConnector(ctx, clients.ConnectConnectorCreateRequest{Name: data.Name.ValueString(), Config: config})
	if err != nil {
		resp.Diagnostics.AddError("Unable to create Kafka Connect connector", err.Error())
		return
	}
	setConnectorModel(ctx, &data, connector, "")
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ConnectConnectorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ConnectConnectorResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	connector, err := r.client.GetConnector(ctx, data.Name.ValueString())
	if err != nil {
		if apiErr, ok := err.(*clients.APIError); ok && apiErr.NotFound() {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Unable to read Kafka Connect connector", err.Error())
		return
	}
	status, statusErr := r.client.GetConnectorStatus(ctx, data.Name.ValueString())
	statusValue := ""
	if statusErr == nil {
		statusValue = status.Connector.State
	}
	setConnectorModel(ctx, &data, connector, statusValue)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ConnectConnectorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ConnectConnectorResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	config := connectorConfig(ctx, data.Config, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	updated, err := r.client.UpdateConnectorConfig(ctx, data.Name.ValueString(), config)
	if err != nil {
		resp.Diagnostics.AddError("Unable to update Kafka Connect connector configuration", err.Error())
		return
	}
	connector := clients.ConnectConnector{Name: data.Name.ValueString(), Config: updated}
	setConnectorModel(ctx, &data, connector, "")
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ConnectConnectorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ConnectConnectorResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.DeleteConnector(ctx, data.Name.ValueString()); err != nil {
		if apiErr, ok := err.(*clients.APIError); ok && apiErr.NotFound() {
			return
		}
		resp.Diagnostics.AddError("Unable to delete Kafka Connect connector", err.Error())
	}
}

func (r *ConnectConnectorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if strings.TrimSpace(req.ID) == "" {
		resp.Diagnostics.AddError("Invalid Kafka Connect connector import ID", "Use the connector name.")
		return
	}
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), types.StringValue(req.ID))...)
}

func connectClient(providerData any) (*clients.HTTPClient, error) {
	runtime, ok := providerData.(*clients.Runtime)
	if !ok {
		return nil, fmt.Errorf("expected *clients.Runtime, got %T", providerData)
	}
	client, ok := runtime.Service("connect")
	if !ok {
		return nil, fmt.Errorf("Connect endpoint is not configured")
	}
	return client, nil
}

func connectorConfig(ctx context.Context, value types.Map, diagnostics interface{ AddError(string, string) }) map[string]string {
	var config map[string]string
	if diags := value.ElementsAs(ctx, &config, false); diags.HasError() {
		diagnostics.AddError("Invalid Kafka Connect configuration", "config must be a map of string values.")
		return nil
	}
	return config
}

func setConnectorModel(ctx context.Context, data *ConnectConnectorResourceModel, connector clients.ConnectConnector, status string) {
	data.ID = types.StringValue(connector.Name)
	data.Type = types.StringValue(connector.Type)
	data.State = types.StringValue(status)
	config := make(map[string]string, len(connector.Config))
	if !data.Config.IsNull() && !data.Config.IsUnknown() {
		_ = data.Config.ElementsAs(ctx, &config, false)
	}
	for key, value := range connector.Config {
		if value != "********" && value != "[hidden]" {
			config[key] = value
		}
	}
	data.Config, _ = types.MapValueFrom(ctx, types.StringType, config)
}
