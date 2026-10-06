# Terraform Provider for Confluent Platform

This provider manages self-managed Confluent Platform 8.0 through 8.3 using the Terraform Plugin Framework. REST Proxy, Metadata Service, Kafka Connect, Schema Registry, and ksqlDB support is being delivered incrementally. Flink is currently out of scope.

The provider manages cluster-wide Confluent Platform resources through REST Proxy, Metadata Service, and Kafka Connect. Schema Registry and ksqlDB endpoints are configured for the next implementation phases.

The repository contains:

- Shared service clients (`internal/clients/`) and Terraform surfaces (`internal/provider/`),
- Examples (`examples/`) and generated documentation (`docs/`),
- Miscellaneous planning and release metadata.

Once you've written your provider, you'll want to [publish it on the Terraform Registry](https://developer.hashicorp.com/terraform/registry/providers/publishing) so that others can use it.

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.24
- A self-managed Confluent Platform 8.0-8.3 deployment for live use.

## Building the Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Using the Provider

Provider endpoints may be configured explicitly or through environment variables. Explicit Terraform configuration takes precedence.

```terraform
terraform {
	required_providers {
		confluent-platform = {
			source = "variananora/confluent-platform"
		}
	}
}

provider "confluent-platform" {
	rest_proxy_endpoint       = "https://kafka-rest.example.com"
	mds_endpoint              = "https://mds.example.com"
	connect_endpoint          = "https://connect.example.com"
	schema_registry_endpoint = "https://schema-registry.example.com"
	ksqldb_endpoint           = "https://ksqldb.example.com"

	username        = var.confluent_username
	password        = var.confluent_password
	timeout_seconds = 30
	max_retries     = 2
}
```

Authentication and transport settings include Basic authentication, bearer tokens, custom CA certificates, client certificates, TLS server-name overrides, insecure TLS mode, request timeouts, bounded retries, and a custom user agent. Sensitive values should be supplied through variables or environment variables.

Environment variables use the `CONFLUENT_PLATFORM_` prefix, including `CONFLUENT_PLATFORM_REST_PROXY_ENDPOINT`, `CONFLUENT_PLATFORM_MDS_ENDPOINT`, `CONFLUENT_PLATFORM_CONNECT_ENDPOINT`, `CONFLUENT_PLATFORM_SCHEMA_REGISTRY_ENDPOINT`, `CONFLUENT_PLATFORM_KSQLDB_ENDPOINT`, `CONFLUENT_PLATFORM_USERNAME`, `CONFLUENT_PLATFORM_PASSWORD`, and `CONFLUENT_PLATFORM_BEARER_TOKEN`.

Explicit provider arguments take precedence over environment variables. TLS settings include `CONFLUENT_PLATFORM_CA_CERT_PEM`, `CONFLUENT_PLATFORM_CLIENT_CERT_PEM`, `CONFLUENT_PLATFORM_CLIENT_KEY_PEM`, `CONFLUENT_PLATFORM_TLS_SERVER_NAME`, and `CONFLUENT_PLATFORM_INSECURE_SKIP_VERIFY`. Runtime settings include `CONFLUENT_PLATFORM_TIMEOUT_SECONDS`, `CONFLUENT_PLATFORM_MAX_RETRIES`, and `CONFLUENT_PLATFORM_USER_AGENT`.

## Implemented Surfaces

Current Terraform resources include Kafka topics and ACLs, MDS role bindings, centralized MDS ACLs, MDS cluster registry, MDS audit configuration, and Kafka Connect connectors.

Current data sources include Kafka cluster/broker/topic/config/partition/consumer-group/replica/cluster-link/share-group/Streams-group reads; MDS active nodes, metadata cluster ID, features, roles, role names, role summaries, audit routes, and audit lookup; and Kafka Connect cluster, connector info/config/status.

Schema Registry and ksqlDB are configured endpoints but are not yet implemented. Kafka Connect task, plugin, topic, offset, and lifecycle action surfaces are also pending.

## Security and State

Credentials, bearer tokens, TLS key material, and connector configuration maps are sensitive. Kafka Connect may return masked secret placeholders; the provider retains prior state for those keys and never attempts to recover or log them. Raw JSON data sources are intended for compatibility and may contain deployment-specific fields, so review state handling before storing sensitive API responses.

MDS role bindings use stable composite import IDs documented on the resource page. Cluster registry and audit configuration JSON must be reviewed for credentials and endpoint data before committing Terraform state.

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `make generate`.

`make generate` requires both Terraform and the Terraform Plugin Docs toolchain. It formats examples with Terraform and regenerates `docs/`; the command fails when Terraform is unavailable.

In order to run the full suite of Acceptance tests, run `make testacc`.

_Note:_ Acceptance tests create real resources, and often cost money to run.

Acceptance tests require `TF_ACC=1`, explicit `CONFLUENT_PLATFORM_*` environment variables, and a dedicated non-production Confluent Platform deployment. They must not run implicitly against production.

```shell
make testacc
```
