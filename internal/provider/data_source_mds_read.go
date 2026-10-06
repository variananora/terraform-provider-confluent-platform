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

type MDSReadDataSource struct {
	kind   string
	client *clients.HTTPClient
}

type MDSReadDataSourceModel struct {
	ID             types.String `tfsdk:"id"`
	Principal      types.String `tfsdk:"principal"`
	RoleName       types.String `tfsdk:"role_name"`
	ResourceType   types.String `tfsdk:"resource_type"`
	ResourceName   types.String `tfsdk:"resource_name"`
	KafkaClusterID types.String `tfsdk:"kafka_cluster_id"`
	Query          types.String `tfsdk:"query"`
	CRN            types.String `tfsdk:"crn"`
	ResponseJSON   types.String `tfsdk:"response_json"`
}

func NewMDSFeaturesDataSource() datasource.DataSource  { return &MDSReadDataSource{kind: "features"} }
func NewMDSRolesDataSource() datasource.DataSource     { return &MDSReadDataSource{kind: "roles"} }
func NewMDSRoleNamesDataSource() datasource.DataSource { return &MDSReadDataSource{kind: "role_names"} }
func NewMDSRoleBindingSummaryDataSource() datasource.DataSource {
	return &MDSReadDataSource{kind: "role_binding_summary"}
}
func NewMDSAuditRoutesDataSource() datasource.DataSource {
	return &MDSReadDataSource{kind: "audit_routes"}
}
func NewMDSAuditLookupDataSource() datasource.DataSource {
	return &MDSReadDataSource{kind: "audit_lookup"}
}

func (d *MDSReadDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mds_" + d.kind
}

func (d *MDSReadDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	attrs := map[string]schema.Attribute{
		"id":               schema.StringAttribute{Computed: true},
		"principal":        schema.StringAttribute{Optional: true},
		"role_name":        schema.StringAttribute{Optional: true},
		"resource_type":    schema.StringAttribute{Optional: true},
		"resource_name":    schema.StringAttribute{Optional: true},
		"kafka_cluster_id": schema.StringAttribute{Optional: true},
		"query":            schema.StringAttribute{Optional: true},
		"crn":              schema.StringAttribute{Optional: true},
		"response_json":    schema.StringAttribute{Computed: true},
	}
	resp.Schema = schema.Schema{MarkdownDescription: "Reads Metadata Service " + d.kind + " information as JSON.", Attributes: attrs}
}

func (d *MDSReadDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, err := mdsClient(req.ProviderData)
	if err != nil {
		resp.Diagnostics.AddError("Unable to configure MDS data source", err.Error())
		return
	}
	d.client = client
}

func (d *MDSReadDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data MDSReadDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	path, payload, err := d.read(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError("Unable to read MDS "+d.kind, err.Error())
		return
	}
	data.ID = types.StringValue(path)
	data.ResponseJSON = types.StringValue(string(payload))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (d *MDSReadDataSource) read(ctx context.Context, data MDSReadDataSourceModel) (string, []byte, error) {
	scope := clients.MDSRoleBindingScope{Clusters: map[string]string{"kafka-cluster": data.KafkaClusterID.ValueString()}}
	switch d.kind {
	case "features":
		path := "/security/1.0/features"
		payload, err := d.client.GetMDSJSON(ctx, path)
		return path, payload, err
	case "roles":
		path := "/security/1.0/roles"
		if !data.RoleName.IsNull() && data.RoleName.ValueString() != "" {
			path += "/" + clients.EscapePathSegment(data.RoleName.ValueString())
		}
		payload, err := d.client.GetMDSJSON(ctx, path)
		return path, payload, err
	case "role_names":
		path := "/security/1.0/roleNames"
		payload, err := d.client.GetMDSJSON(ctx, path)
		return path, payload, err
	case "role_binding_summary":
		if data.Principal.IsNull() || data.Principal.ValueString() == "" {
			return "", nil, fmt.Errorf("principal is required")
		}
		path := "/security/1.0/lookup/principal/" + clients.EscapePathSegment(data.Principal.ValueString()) + "/resources"
		payload, err := d.client.PostMDSJSON(ctx, path, scope)
		return path, payload, err
	case "audit_routes":
		path := "/security/1.0/audit/routes"
		payload, err := d.client.GetMDSAuditRoutes(ctx, data.Query.ValueString())
		return path, payload, err
	case "audit_lookup":
		if data.CRN.IsNull() || data.CRN.ValueString() == "" {
			return "", nil, fmt.Errorf("crn is required")
		}
		path := "/security/1.0/audit/lookup"
		payload, err := d.client.GetMDSAuditLookup(ctx, data.CRN.ValueString())
		return path, payload, err
	default:
		return "", nil, fmt.Errorf("unsupported MDS data source %q", d.kind)
	}
}
