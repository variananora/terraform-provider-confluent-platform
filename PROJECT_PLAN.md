# Confluent Platform Terraform Provider

## Goal

Build a Terraform Plugin Framework provider named `confluent-platform` for self-managed Confluent Platform. The provider will manage and read cluster-wide resources through the REST Proxy, Metadata Service, Kafka Connect, Schema Registry, and ksqlDB APIs. Flink is out of scope.

The endpoint inventory is a long-term backlog, delivered in stages. Each API surface must have a typed client model, Terraform implementation, unit tests, documentation, and examples.

## Terraform Modeling Rules

Classify endpoints by behavior, not only by HTTP method:

- **Data source:** stable reads, lists, status responses, and lookups.
- **Resource:** persistent server-side objects with stable identity and meaningful Create, Read, Update, and Delete behavior. Create/update may use POST, PUT, or PATCH.
- **Action:** imperative operations such as restart, pause, resume, reset, terminate, publish records, or token operations.
- **Function or ephemeral resource:** only when the result is naturally transient and the Terraform Framework feature is appropriate.
- **Client-only:** unsafe, internal, private UI, or otherwise unsuitable for Terraform state.

## Phases

### Phase 0: Scope and API Inventory

1. Complete the endpoint catalog for REST Proxy v3, MDS, Connect, Schema Registry, and ksqlDB.
2. Record path, method, API version, request/response model, authentication, permissions, idempotency, pagination, asynchronous behavior, delete semantics, and Terraform exposure for every endpoint.
3. Target Confluent Platform 8.0 through 8.3. Treat the REST contracts as compatible across this range unless live testing identifies a deployment-specific difference.
4. Prioritize the first release around topics, topic configs, ACLs, Connectors, Schema Registry subjects/schemas/configuration, and MDS role bindings/cluster registry.

The baseline catalog is maintained in `API_CATALOG.md`. Phase 0 planning is complete for the 8.0-8.3 target; release-specific endpoint availability, permissions, optional extensions, and response fields are verified as each client is implemented.

### Phase 1: Provider and Shared Foundation

1. Complete the provider rename to `confluent-platform` across module path, metadata, server address, examples, tests, docs generation, and README.
2. Support separate endpoints for REST Proxy, MDS, Connect, Schema Registry, and ksqlDB.
3. Add environment-variable configuration with documented precedence.
4. Add TLS options: custom CA, client certificate/key, server name override, and insecure mode.
5. Add Basic authentication, bearer tokens, OAuth/OIDC inputs, request timeout, retry limits, and user agent.
6. Replace raw `*http.Client` provider data with a shared runtime containing typed service clients, auth, transport, retries, timeouts, and diagnostics.
7. Centralize request creation, context cancellation, content negotiation, status/error decoding, safe retries, URL escaping, and 404 handling.
8. Log method, sanitized URL, status, duration, retry count, and request identifiers with `tflog`. Never log credentials, tokens, connector secrets, schema-registry credentials, raw SQL, or sensitive response data.
9. Add unit tests with `httptest` for configuration, TLS/auth, redaction, retries, errors, headers, URL handling, and diagnostics.
10. Remove stock scaffold surfaces once the real provider registration is validated.

Current foundation already includes the renamed provider identity, initial service endpoint fields, `tflog` configuration logging, and a shared `internal/clients` HTTP client with endpoint validation and timing logs. The current scaffold resources remain temporarily registered during migration.

### Phase 2: REST Proxy / Kafka v3

Implement typed clients and Terraform surfaces for:

- Clusters, brokers, replicas, replica status, and configs
- Topics, partitions, records, consumer groups, share groups, and streams groups
- ACLs and cluster linking
- Balancer status, broker tasks, broker replica exclusions, remove broker, unregister, and license operations

Prioritize resources for topics, topic configurations, ACL bindings, and cluster links. Add data sources for cluster, broker, topic, partition, config, consumer group, replica/status, and link/status endpoints. Model records and broker operations as actions or client-only operations. Normalize unordered collections and define pagination and polling behavior.

### Phase 3: Metadata Service / MDS

Implement typed MDS clients for authentication, active nodes, metadata cluster ID, features, role definitions, role-binding CRUD and summaries, centralized ACLs, cluster registry, and audit configuration.

Prioritize resources for role bindings, cluster registry entries, centralized ACLs, and audit configuration where ownership and deletion are clear. Use stable composite import IDs. Preserve audit configuration resource-version conflict behavior. Expose role definitions, role names, registered clusters, summaries, active nodes, features, and audit lookups as data sources. Treat private Control Center UI endpoints and token-management flows as actions or client helpers only after a security review.

### Phase 4: Kafka Connect

Implement a Connect client using JSON content negotiation, connector-specific errors, masked-secret handling, rebalance conflict handling, and asynchronous status polling.

Prioritize a connector resource with name, configuration, initial state, import, full CRUD, and sensitive configuration behavior. Use POST for creation and PUT for configuration updates. Never attempt to recover secrets masked by Connect. Add data sources for cluster info, connector info/config/status, tasks, plugins, topics, offsets, and validation. Add actions for restart, pause, resume, stop, offset alteration/reset, and topic reset.

### Phase 5: Schema Registry

Implement a Schema Registry client with v1 content negotiation, context-aware paths, JSON-encoded schema strings, standard error codes, forwarding/timeouts, and soft/hard delete semantics.

Prioritize resources for schema registration, subject/global compatibility configuration, mutable modes, exporters, importers, KEKs, and DEKs where supported by the deployment. Model schema ID, version, and GUID separately. Add data sources for subjects, versions, schemas by ID/GUID/version, schema lists, compatibility, modes, contexts, exporter/importer status/config, supported types, KEKs, DEKs, and policy. Use actions for compatibility tests, exporter/importer lifecycle operations, KEK tests/undelete, and DEK undelete. Require explicit opt-in for permanent deletion.

### Phase 6: ksqlDB

Implement a ksqlDB client with service-specific media types, statement execution, status polling, server/cluster introspection, property validation, streaming cancellation, and timeouts.

Add data sources for server status, cluster status, statement/query status, and property validity. Use resources only for durable objects with a reliable reconciled lifecycle. Model statement execution, push/pull queries, and cluster termination as actions. Do not log raw SQL or unbounded streaming results.

### Phase 7: Documentation and Release

1. Generate documentation for every registered resource, data source, and action.
2. Add service-specific examples with endpoints, authentication, imports, permissions, lifecycle caveats, and destructive-operation warnings.
3. Document supported Platform versions, endpoint prerequisites, environment variables, TLS/auth combinations, sensitive state limitations, polling, and unsupported endpoints.
4. Add CI gates for formatting, vet, lint, build, unit tests, generated docs, provider schema validation, and opt-in acceptance tests.
5. Maintain the endpoint catalog as a checked-off inventory linking each endpoint to client code, Terraform surface, tests, docs, examples, and acceptance coverage.

## File Organization

- `internal/clients/`: shared transport, authentication, TLS, errors, logging, and service-specific clients.
- `internal/provider/<service>_<surface>.go`: one Terraform resource, data source, or action per API surface.
- `internal/provider/<service>_<surface>_test.go`: focused unit and acceptance tests for that surface.
- `examples/`: provider configuration and service-specific Terraform examples.
- `docs/`: generated Terraform documentation and service guidance.
- `PLAN.md`: original API scope notes.
- `PROJECT_PLAN.md`: expanded architecture and delivery plan.
- `IMPLEMENTATION_PROMPT.md`: reusable coding-agent handoff prompt.
- `API_CATALOG.md`: Phase 0 endpoint families, exposure decisions, and contract questions.

## Verification Gates

1. Run `gofmt`, `go vet`, lint, `go build ./...`, and focused unit tests after every implementation slice.
2. Use `httptest` fixtures to verify methods, paths, headers, bodies, auth, timeouts, retries, errors, and redaction.
3. Test Terraform state after Create, Read refresh, Update, Delete/not-found, and Import.
4. Run acceptance tests only with explicit environment variables and a dedicated self-managed Confluent Platform deployment.
5. Run documentation generation and fail when generated output is stale.
6. Manually verify Terraform debug logs contain request metadata but no secrets, tokens, connector passwords, schema credentials, raw SQL, or sensitive payloads.

## Decisions

- Provider name: `confluent-platform`.
- Full supplied API list is the backlog; implementation is staged.
- Basic auth, bearer/OAuth/OIDC, and TLS are first-class requirements.
- Terraform exposure is semantic rather than a direct HTTP verb mapping.
- Flink is excluded.
- Private RBAC UI and token-management endpoints require a concrete use case and security review.
- Acceptance tests must never run implicitly against production.

## Open Decisions

1. Confirm live endpoint availability and optional extensions across Confluent Platform 8.0-8.3 deployments.
2. Decide which schemas, connector values, encrypted key material, and query text may be stored in Terraform state.
3. Decide whether one provider instance represents one Platform environment or whether aliases must support multiple environments immediately.
