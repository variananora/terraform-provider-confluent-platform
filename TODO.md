# Confluent Platform Provider TODO

Use this checklist to track implementation progress. Check an item only after the code, focused tests, and documentation for that item are complete.

Legend: `[x]` complete, `[ ]` not started or incomplete, `[~]` partially complete.

## Current Foundation

- [x] Rename provider metadata to `confluent-platform`.
- [x] Update the Go module path to `github.com/variananora/terraform-provider-confluent-platform`.
- [x] Update the provider server address in `main.go`.
- [x] Update provider acceptance factory names.
- [x] Update the provider example configuration.
- [x] Add initial endpoint attributes for REST Proxy, MDS, Connect, Schema Registry, and ksqlDB.
- [~] Complete provider rename across all README, generated documentation, examples, and remaining scaffold identifiers.
- [x] Add initial `tflog` provider configuration logging.
- [x] Add shared HTTP client under `internal/clients`.
- [x] Validate endpoint scheme and host.
- [x] Apply HTTP client timeout defaults.
- [x] Add request context propagation.
- [x] Add request and response timing logs.
- [x] Add `httptest` coverage for shared HTTP client behavior.
- [x] Confirm `go build ./...` succeeds.
- [ ] Install or configure Terraform CLI for provider integration and acceptance tests.

## Phase 0: Scope and API Inventory

- [x] Define the supported Confluent Platform version range as 8.0-8.3.
- [~] Record version-dependent endpoint and response-field differences; no breaking API differences are currently assumed, but live deployment verification remains open.
- [~] Build the baseline endpoint catalog for REST Proxy v3 in `API_CATALOG.md`.
- [~] Build the baseline endpoint catalog for MDS in `API_CATALOG.md`.
- [~] Build the baseline endpoint catalog for Kafka Connect in `API_CATALOG.md`.
- [~] Build the baseline endpoint catalog for Schema Registry in `API_CATALOG.md`.
- [~] Build the baseline endpoint catalog for ksqlDB in `API_CATALOG.md`.
- [ ] Record authentication and permission requirements for every endpoint.
- [ ] Record pagination, idempotency, retry, async, and deletion semantics.
- [~] Classify the baseline endpoint families as resource, data source, action, or client-only.
- [x] Record first-release priorities: topics, configs, ACLs, connectors, schemas, and RBAC bindings.

## Phase 1: Provider and Client Foundation

- [x] Add environment-variable configuration and documented precedence.
- [x] Add provider request timeout configuration.
- [x] Add retry limit and retry backoff configuration.
- [x] Add user-agent configuration.
- [x] Add Basic authentication configuration.
- [x] Add bearer-token configuration.
- [~] Add OAuth/OIDC configuration where supported; bearer-token input is implemented, provider-managed token exchange remains deployment-dependent.
- [x] Add custom CA certificate configuration.
- [x] Add client certificate and private key configuration.
- [x] Add TLS server-name override and insecure-mode configuration.
- [x] Replace raw `*http.Client` provider data with a shared provider runtime.
- [ ] Create separate typed clients for REST Proxy, MDS, Connect, Schema Registry, and ksqlDB.
- [~] Add service-specific base URLs and content negotiation; base URLs and shared transport are implemented, service-specific media types remain with typed clients.
- [x] Add structured API error decoding.
- [x] Add safe transient retry classification for idempotent requests.
- [x] Add 404/not-found handling helpers.
- [x] Add URL escaping helpers.
- [x] Add authorization and secret redaction tests.
- [x] Add TLS and authentication tests.
- [x] Add shared HTTP client tests for bearer auth, user agent, TLS minimum, and invalid client certificates.
- [x] Add retry and error-decoding tests.
- [x] Remove or quarantine stock scaffold resources, data sources, actions, functions, and ephemeral resources.

## Phase 2: REST Proxy / Kafka v3

### Client and Models

- [x] Add typed cluster client and models.
- [~] Add typed broker and replica client/models; broker models are complete, replica models remain.
- [~] Add config client/model; REST Proxy read is implemented with compatibility-first JSON output, richer typed fields remain.
- [~] Add typed topic and partition client/models; topic models and CRUD/read client are complete, partition read is implemented with JSON output.
- [ ] Add typed record client/models.
- [x] Add typed ACL client/models.
- [~] Add consumer-group, cluster-link, share-group, and Streams-group read clients with compatibility-first JSON output.
- [ ] Add typed balancer and broker-task client/models.
- [ ] Add typed broker-replica-exclusion client/models.
- [ ] Add typed unregister, remove-broker, license, and replica-status models.

### Terraform Surfaces

- [x] Implement topic resource schema, Create, Read, Update, Delete, Import, and provider registration.
- [x] Implement topic configuration resource or nested configuration lifecycle.
- [x] Implement ACL resource with stable identity and import format.
- [ ] Implement cluster-link resource.
- [x] Implement and register cluster data source with REST Proxy Read implementation.
- [x] Implement and register broker data source with REST Proxy Read implementation.
- [x] Implement and register topic data source with REST Proxy Read implementation.
- [x] Implement and register partition data source with JSON response output.
- [x] Implement and register config data source with JSON response output.
- [x] Implement and register consumer-group data source with JSON response output.
- [~] Implement and register replica data source with JSON response output; replica-status-specific modeling remains.
- [~] Implement and register cluster-link data source with JSON response output; status-specific modeling remains.
- [x] Implement and register share-group and Streams-group data sources with JSON response output.
- [ ] Implement records publish action.
- [ ] Implement broker-task and broker-operation actions.
- [ ] Implement balancer and replica-exclusion actions where appropriate.
- [ ] Define pagination and collection ordering behavior.
- [x] Add mocked HTTP tests for cluster, broker, and topic REST Proxy reads.
- [ ] Add REST Proxy acceptance tests for prioritized resources.
- [~] Add REST Proxy docs and examples; topic and ACL resource coverage is documented, remaining surfaces are pending.

## Phase 3: Metadata Service / MDS

- [x] Implement MDS authentication and token handling helpers.
- [x] Implement active-nodes client.
- [x] Implement metadata-cluster-ID client.
- [x] Implement feature-discovery client.
- [x] Implement role definitions and role names data sources.
- [x] Implement role-binding resource with composite identity and Import.
- [x] Implement role-binding summary data sources.
- [x] Implement centralized ACL resource/data source behavior.
- [x] Implement cluster-registry resource with stable identity and Import.
- [x] Implement audit-configuration resource with resource-version conflict handling.
- [x] Implement audit routes and lookup data sources.
- [x] Define permission diagnostics for 401 and 403 responses.
- [x] Validate cluster-name versus cluster-ID request rules.
- [~] Add deterministic ordering for role and resource collections; compatibility-first JSON surfaces preserve API ordering pending typed collection models.
- [x] Review private RBAC UI and token operations for action/client-only treatment.
- [~] Add mocked MDS tests; core clients are covered, remaining read/resource fixtures and acceptance tests are pending.
- [ ] Add MDS acceptance tests for prioritized resources.
- [x] Add MDS docs and examples for implemented surfaces.

## Phase 4: Kafka Connect

- [x] Implement Connect typed client with JSON content negotiation.
- [ ] Implement connector error decoding.
- [x] Handle masked connector secrets without attempting recovery.
- [ ] Handle rebalance conflicts and bounded retries.
- [ ] Implement asynchronous status polling.
- [x] Implement connector resource with Create, Read, Update, Delete, and Import.
- [~] Implement connector configuration and initial-state behavior; configuration lifecycle is implemented, initial-state actions remain pending.
- [x] Implement Connect cluster data source.
- [x] Implement connector info/config/status data sources.
- [ ] Implement task list and task-status data sources.
- [ ] Implement plugin list and validation data sources.
- [ ] Implement connector topics data source.
- [ ] Implement connector offsets data source.
- [ ] Implement restart, pause, resume, stop, and topic-reset actions.
- [ ] Implement offset alteration and reset actions with stopped-state validation.
- [~] Add mocked Connect tests; connector lifecycle coverage is implemented, remaining endpoints are pending.
- [ ] Add Connect acceptance tests for a distributed cluster.
- [~] Add Connect docs and examples; connector resource coverage is documented, remaining surfaces are pending.

## Phase 5: Schema Registry

- [ ] Implement Schema Registry v1 content negotiation.
- [ ] Implement context-aware URL handling.
- [ ] Implement schema error-code decoding.
- [ ] Implement schema registration resource.
- [ ] Model schema ID, subject version, and GUID separately.
- [ ] Implement subject/global compatibility resources.
- [ ] Implement mutable subject/global mode resources where enabled.
- [ ] Implement exporter resource and lifecycle handling.
- [ ] Implement importer resource and lifecycle handling.
- [ ] Implement KEK resource.
- [ ] Implement DEK resource where deployment support is available.
- [ ] Implement subjects and versions data sources.
- [ ] Implement schema-by-ID, schema-by-GUID, and schema-by-version data sources.
- [ ] Implement schema list and supported-types data sources.
- [ ] Implement compatibility-check data source or action.
- [ ] Implement contexts and mode data sources.
- [ ] Implement exporter/importer status and config data sources.
- [ ] Implement KEK and DEK data sources.
- [ ] Implement exporter/importer pause, resume, and reset actions.
- [ ] Implement KEK test and undelete actions.
- [ ] Implement DEK undelete actions.
- [ ] Require explicit opt-in for permanent deletion.
- [ ] Test Avro, JSON, Protobuf, references, normalization, and context routing.
- [ ] Test soft-delete and hard-delete behavior.
- [ ] Add Schema Registry acceptance tests.
- [ ] Add Schema Registry docs and examples.

## Phase 6: ksqlDB

- [ ] Implement ksqlDB-specific media types.
- [ ] Implement statement execution client.
- [ ] Implement statement status polling.
- [ ] Implement server and cluster status clients.
- [ ] Implement property validation client.
- [ ] Implement streaming response cancellation and limits.
- [ ] Implement server-status data source.
- [ ] Implement cluster-status data source.
- [ ] Implement statement/query-status data source.
- [ ] Implement property-validity data source.
- [ ] Define durable ksqlDB resources only where lifecycle reconciliation is reliable.
- [ ] Implement statement execution action.
- [ ] Implement push/pull query actions with bounded results.
- [ ] Implement cluster-termination action with explicit safeguards.
- [ ] Prevent raw SQL and unbounded result logging.
- [ ] Add mocked ksqlDB tests.
- [ ] Add ksqlDB acceptance tests.
- [ ] Add ksqlDB docs and examples.

## Phase 7: Documentation, CI, and Release

- [x] Update README with provider installation and configuration.
- [x] Document supported Platform versions and service prerequisites.
- [x] Document environment variables and configuration precedence.
- [x] Document TLS and authentication combinations.
- [x] Document sensitive state limitations.
- [ ] Document async polling and eventual consistency behavior.
- [ ] Document destructive operations and permanent deletion flags.
- [~] Generate docs for every registered resource, data source, and action; manual pages are current, tfplugindocs requires Terraform CLI.
- [x] Add examples for every prioritized Terraform surface.
- [ ] Add `go vet` CI check.
- [ ] Add lint CI check.
- [ ] Add build CI check.
- [ ] Add unit-test CI check.
- [ ] Add generated-doc consistency CI check.
- [ ] Add provider schema validation CI check.
- [ ] Add opt-in acceptance-test CI workflow.
- [ ] Maintain the endpoint catalog as implementation progresses.
- [ ] Update CHANGELOG for each provider release.
- [ ] Verify registry manifest, version injection, and release packaging.

## Open Decisions

- [ ] Confirm minimum and maximum supported Confluent Platform versions.
- [ ] Decide which schemas, connector values, encrypted key material, and query text may be stored in Terraform state.
- [ ] Decide whether one provider instance represents one Platform environment or aliases support multiple environments immediately.
- [ ] Confirm exact registry namespace and published provider address.
