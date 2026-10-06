// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/variananora/terraform-provider-confluent-platform/internal/clients"
)

var _ provider.Provider = &ConfluentPlatformProvider{}
var _ provider.ProviderWithFunctions = &ConfluentPlatformProvider{}
var _ provider.ProviderWithEphemeralResources = &ConfluentPlatformProvider{}
var _ provider.ProviderWithActions = &ConfluentPlatformProvider{}

// ConfluentPlatformProvider defines the provider implementation.
type ConfluentPlatformProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// ConfluentPlatformProviderModel describes the provider data model.
type ConfluentPlatformProviderModel struct {
	RestProxyEndpoint      types.String `tfsdk:"rest_proxy_endpoint"`
	MdsEndpoint            types.String `tfsdk:"mds_endpoint"`
	ConnectEndpoint        types.String `tfsdk:"connect_endpoint"`
	SchemaRegistryEndpoint types.String `tfsdk:"schema_registry_endpoint"`
	KsqlDBEndpoint         types.String `tfsdk:"ksqldb_endpoint"`
	Username               types.String `tfsdk:"username"`
	Password               types.String `tfsdk:"password"`
	BearerToken            types.String `tfsdk:"bearer_token"`
	CACertPEM              types.String `tfsdk:"ca_cert_pem"`
	ClientCertPEM          types.String `tfsdk:"client_cert_pem"`
	ClientKeyPEM           types.String `tfsdk:"client_key_pem"`
	TLSServerName          types.String `tfsdk:"tls_server_name"`
	InsecureSkipVerify     types.Bool   `tfsdk:"insecure_skip_verify"`
	TimeoutSeconds         types.Int64  `tfsdk:"timeout_seconds"`
	MaxRetries             types.Int64  `tfsdk:"max_retries"`
	UserAgent              types.String `tfsdk:"user_agent"`
}

func (p *ConfluentPlatformProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "confluent-platform"
	resp.Version = p.version
}

func (p *ConfluentPlatformProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"rest_proxy_endpoint": schema.StringAttribute{
				MarkdownDescription: "Kafka REST Proxy endpoint.",
				Optional:            true,
			},
			"mds_endpoint": schema.StringAttribute{
				MarkdownDescription: "Metadata Service endpoint.",
				Optional:            true,
			},
			"connect_endpoint": schema.StringAttribute{
				MarkdownDescription: "Kafka Connect REST endpoint.",
				Optional:            true,
			},
			"schema_registry_endpoint": schema.StringAttribute{
				MarkdownDescription: "Schema Registry endpoint.",
				Optional:            true,
			},
			"ksqldb_endpoint": schema.StringAttribute{
				MarkdownDescription: "ksqlDB REST endpoint.",
				Optional:            true,
			},
			"username":             schema.StringAttribute{Optional: true, MarkdownDescription: "Username for Basic authentication."},
			"password":             schema.StringAttribute{Optional: true, Sensitive: true, MarkdownDescription: "Password for Basic authentication."},
			"bearer_token":         schema.StringAttribute{Optional: true, Sensitive: true, MarkdownDescription: "Bearer token for API authentication."},
			"ca_cert_pem":          schema.StringAttribute{Optional: true, Sensitive: true, MarkdownDescription: "Custom CA certificate in PEM format."},
			"client_cert_pem":      schema.StringAttribute{Optional: true, Sensitive: true, MarkdownDescription: "Client certificate in PEM format."},
			"client_key_pem":       schema.StringAttribute{Optional: true, Sensitive: true, MarkdownDescription: "Client private key in PEM format."},
			"tls_server_name":      schema.StringAttribute{Optional: true, MarkdownDescription: "TLS server name override."},
			"insecure_skip_verify": schema.BoolAttribute{Optional: true, MarkdownDescription: "Disable TLS certificate verification."},
			"timeout_seconds":      schema.Int64Attribute{Optional: true, MarkdownDescription: "HTTP request timeout in seconds."},
			"max_retries":          schema.Int64Attribute{Optional: true, MarkdownDescription: "Maximum retries for safe idempotent requests."},
			"user_agent":           schema.StringAttribute{Optional: true, MarkdownDescription: "User-Agent header for API requests."},
		},
	}
}

func (p *ConfluentPlatformProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data ConfluentPlatformProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	endpoints := map[string]string{
		"rest_proxy":      resolveString(data.RestProxyEndpoint, "CONFLUENT_PLATFORM_REST_PROXY_ENDPOINT"),
		"mds":             resolveString(data.MdsEndpoint, "CONFLUENT_PLATFORM_MDS_ENDPOINT"),
		"connect":         resolveString(data.ConnectEndpoint, "CONFLUENT_PLATFORM_CONNECT_ENDPOINT"),
		"schema_registry": resolveString(data.SchemaRegistryEndpoint, "CONFLUENT_PLATFORM_SCHEMA_REGISTRY_ENDPOINT"),
		"ksqldb":          resolveString(data.KsqlDBEndpoint, "CONFLUENT_PLATFORM_KSQLDB_ENDPOINT"),
	}

	tflog.Debug(ctx, "configured Confluent Platform provider", map[string]any{
		"rest_proxy_configured":      endpoints["rest_proxy"] != "",
		"mds_configured":             endpoints["mds"] != "",
		"connect_configured":         endpoints["connect"] != "",
		"schema_registry_configured": endpoints["schema_registry"] != "",
		"ksqldb_configured":          endpoints["ksqldb"] != "",
	})

	timeoutSeconds := resolveInt64(data.TimeoutSeconds, "CONFLUENT_PLATFORM_TIMEOUT_SECONDS", 30)
	maxRetries := resolveInt64(data.MaxRetries, "CONFLUENT_PLATFORM_MAX_RETRIES", 2)
	insecureSkipVerify := resolveBool(data.InsecureSkipVerify, "CONFLUENT_PLATFORM_INSECURE_SKIP_VERIFY", false)
	runtime, err := clients.NewRuntime(clients.RuntimeConfig{
		Endpoints:  endpoints,
		Timeout:    time.Duration(timeoutSeconds) * time.Second,
		MaxRetries: maxRetries,
		UserAgent:  resolveString(data.UserAgent, "CONFLUENT_PLATFORM_USER_AGENT"),
		Auth: clients.AuthConfig{
			Username:    resolveString(data.Username, "CONFLUENT_PLATFORM_USERNAME"),
			Password:    resolveString(data.Password, "CONFLUENT_PLATFORM_PASSWORD"),
			BearerToken: resolveString(data.BearerToken, "CONFLUENT_PLATFORM_BEARER_TOKEN"),
		},
		TLS: clients.TLSConfig{
			CACertPEM:          resolveString(data.CACertPEM, "CONFLUENT_PLATFORM_CA_CERT_PEM"),
			ClientCertPEM:      resolveString(data.ClientCertPEM, "CONFLUENT_PLATFORM_CLIENT_CERT_PEM"),
			ClientKeyPEM:       resolveString(data.ClientKeyPEM, "CONFLUENT_PLATFORM_CLIENT_KEY_PEM"),
			ServerName:         resolveString(data.TLSServerName, "CONFLUENT_PLATFORM_TLS_SERVER_NAME"),
			InsecureSkipVerify: insecureSkipVerify,
		},
	})
	if err != nil {
		resp.Diagnostics.AddError("Unable to configure Confluent Platform clients", err.Error())
		return
	}

	resp.DataSourceData = runtime
	resp.ResourceData = runtime
}

func resolveString(value types.String, environmentVariable string) string {
	if !value.IsNull() && !value.IsUnknown() {
		return value.ValueString()
	}
	return os.Getenv(environmentVariable)
}

func resolveInt64(value types.Int64, environmentVariable string, fallback int64) int64 {
	if !value.IsNull() && !value.IsUnknown() {
		return value.ValueInt64()
	}
	if parsed, err := strconv.ParseInt(os.Getenv(environmentVariable), 10, 64); err == nil {
		return parsed
	}
	return fallback
}

func resolveBool(value types.Bool, environmentVariable string, fallback bool) bool {
	if !value.IsNull() && !value.IsUnknown() {
		return value.ValueBool()
	}
	if parsed, err := strconv.ParseBool(os.Getenv(environmentVariable)); err == nil {
		return parsed
	}
	return fallback
}

func (p *ConfluentPlatformProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewKafkaTopicResource,
		NewKafkaACLResource,
		NewMDSRoleBindingResource,
		NewMDSACLResource,
		NewMDSClusterRegistryResource,
		NewMDSAuditConfigResource,
		NewConnectConnectorResource,
	}
}

func (p *ConfluentPlatformProvider) EphemeralResources(ctx context.Context) []func() ephemeral.EphemeralResource {
	return nil
}

func (p *ConfluentPlatformProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewConnectClusterDataSource,
		NewConnectConnectorInfoDataSource,
		NewConnectConnectorConfigDataSource,
		NewConnectConnectorStatusDataSource,
		NewMDSActiveNodesDataSource,
		NewMDSMetadataClusterIDDataSource,
		NewMDSFeaturesDataSource,
		NewMDSRolesDataSource,
		NewMDSRoleNamesDataSource,
		NewMDSRoleBindingSummaryDataSource,
		NewMDSAuditRoutesDataSource,
		NewMDSAuditLookupDataSource,
		NewKafkaClusterDataSource,
		NewKafkaBrokerDataSource,
		NewKafkaTopicDataSource,
		NewKafkaConfigDataSource,
		NewKafkaPartitionDataSource,
		NewKafkaConsumerGroupDataSource,
		NewKafkaReplicaDataSource,
		NewKafkaClusterLinkDataSource,
		NewKafkaShareGroupDataSource,
		NewKafkaStreamsGroupDataSource,
	}
}

func (p *ConfluentPlatformProvider) Functions(ctx context.Context) []func() function.Function {
	return nil
}

func (p *ConfluentPlatformProvider) Actions(ctx context.Context) []func() action.Action {
	return nil
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &ConfluentPlatformProvider{
			version: version,
		}
	}
}
