# Confluent Platform API Catalog

This is the Phase 0 baseline inventory for the provider. The initial support target is Confluent Platform 8.0 through 8.3. We assume the listed REST contracts remain compatible across that range, while still verifying endpoint availability, permissions, optional extensions, and deployment configuration during implementation. A row marked `verify` is inventoried but not yet tested against a live 8.x deployment.

## Classification

- **Resource:** persistent object with a stable identity and reconciled lifecycle.
- **Data source:** read-only lookup, list, status, or metadata result.
- **Action:** imperative operation, query, reset, restart, or asynchronous command.
- **Client-only:** internal, security-sensitive, or unsuitable for Terraform state.

## REST Proxy v3

Reference: [Kafka REST API](https://docs.confluent.io/platform/current/kafka-rest/api.html)

| API family         | Representative endpoints or operations             | Terraform exposure                                  | Contract status                          |
| ------------------ | -------------------------------------------------- | --------------------------------------------------- | ---------------------------------------- |
| Clusters           | `/kafka/v3/clusters`, cluster details              | Data source                                         | typed read client/schema                 |
| Configs            | Cluster, broker, and topic configuration endpoints | Resource for managed configs; data source for reads | JSON read client/schema; resource verify |
| Brokers            | Broker list/details and broker configuration       | Data source                                         | typed read client/schema                 |
| Replicas           | Replica list/details and replica status            | Data source                                         | JSON read client/schema; status verify   |
| ACLs               | ACL list/create/delete/filter operations           | Resource for ACL bindings; data source for searches | typed create/read/delete client/resource |
| Consumer groups    | Group list/details, members, offsets, lag          | Data source                                         | JSON read client/schema                  |
| Partitions         | Partition list/details and assignments             | Data source                                         | JSON read client/schema                  |
| Topics             | Topic list/create/details/configuration/partitions | Resource plus data source                           | typed read client/schema; CRUD verify    |
| Records            | Produce records to a topic                         | Action                                              | verify                                   |
| Cluster linking    | Link list/create/update/delete, link tasks/status  | Resource plus data source/action for operations     | JSON read client/schema; CRUD verify     |
| Share groups       | Share group list/details and members               | Data source                                         | JSON read client/schema                  |
| Balancer status    | Read balancer status                               | Data source                                         | verify                                   |
| Broker tasks       | Read and execute broker tasks                      | Data source/action                                  | verify                                   |
| Replica exclusions | Read/create/delete broker replica exclusions       | Resource or action after lifecycle review           | verify                                   |
| Remove broker      | Remove broker tasks/operations                     | Action                                              | verify                                   |
| Unregister         | Unregister broker or cluster objects               | Action                                              | verify                                   |
| License            | Read and manage license state where supported      | Data source/resource after ownership review         | verify                                   |
| Streams groups     | Streams group list/details and offsets             | Data source                                         | JSON read client/schema                  |

REST Proxy questions to resolve:

- Confirm the exact v3 base path and endpoint availability for the target Platform versions.
- Confirm which topic and cluster configuration operations are mutable and whether deletes restore defaults.
- Confirm stable identity and import formats for links and ACL bindings.
- Confirm pagination tokens, maximum page sizes, and collection ordering.
- Confirm asynchronous operations and polling endpoints for broker/link operations.

## Metadata Service / MDS

Reference: [Metadata API](https://docs.confluent.io/platform/current/security/authorization/rbac/mds-api.html)

| API family           | Representative endpoints                                                                | Terraform exposure                                   | Contract status                                            |
| -------------------- | --------------------------------------------------------------------------------------- | ---------------------------------------------------- | ---------------------------------------------------------- |
| Authentication       | `GET /security/1.0/authenticate`                                                        | Client helper/action; do not store tokens by default | verify                                                     |
| Impersonation        | `POST /security/1.0/impersonate`                                                        | Client-only/action after security review             | verify                                                     |
| Authorization        | `PUT /security/1.0/authorize`                                                           | Action/client-only                                   | verify                                                     |
| Active nodes         | `GET /security/1.0/activenodes/{protocol}`                                              | Data source                                          | compatibility-first JSON client/data source                |
| Metadata cluster ID  | `GET /security/1.0/metadataClusterId`                                                   | Data source                                          | compatibility-first JSON client/data source                |
| Features             | `GET /security/1.0/features`                                                            | Data source                                          | compatibility-first JSON data source                       |
| Role definitions     | `GET /security/1.0/roles`, `/roles/{roleName}`, `/roleNames`                            | Data sources                                         | compatibility-first JSON data sources                      |
| Role binding CRUD    | `POST/PUT/DELETE /security/1.0/principals/{principal}/roles/{roleName}` and `/bindings` | Resource with composite identity                     | typed client/resource for cluster, Topic, and Group scopes |
| Role summaries       | `/security/1.0/lookup/...` role/principal/resource operations                           | Data sources                                         | principal resource summary implemented                     |
| Centralized ACLs     | `POST/DELETE /security/1.0/acls`, `POST /security/1.0/acls:search`                      | Resource plus data source                            | typed client/resource with Kafka cluster scope             |
| Cluster registry     | `GET/POST /security/1.0/registry/clusters`, item GET/DELETE                             | Resource plus data source                            | JSON configuration resource                                |
| Audit configuration  | `GET/PUT /security/1.0/audit/config`                                                    | Resource plus data source                            | resource-version aware JSON resource                       |
| Audit lookup         | `/security/1.0/audit/routes`, `/audit/lookup`                                           | Data sources                                         | query/CRN JSON data sources                                |
| Private RBAC UI      | Visibility, managed rolebindings, creation guidelines, cached principals                | Client-only/data source only after review            | verify                                                     |
| Device authorization | `/security/1.0/oidc/device/*`                                                           | Client helper/action                                 | verify                                                     |
| Token sanitation     | `POST /security/1.0/token/sanitize`                                                     | Client helper/action                                 | verify                                                     |

MDS questions to resolve:

- Document required roles and 401/403 behavior for every exposed operation.
- Define composite IDs containing principal, role, scope, cluster, and resource patterns.
- Enforce the API rule that `clusterName` and `clusters` cannot be sent together.
- Preserve audit `resource_version` and return actionable diagnostics for 409 conflicts.
- Determine which MDS token operations may safely participate in Terraform workflows.

## Kafka Connect

Reference: [Connect REST API](https://docs.confluent.io/platform/current/connect/references/restapi.html)

| API family          | Representative endpoints                                | Terraform exposure                                | Contract status                      |
| ------------------- | ------------------------------------------------------- | ------------------------------------------------- | ------------------------------------ |
| Cluster             | `GET /`                                                 | Data source                                       | JSON data source                     |
| Connectors          | `GET/POST /connectors`, `GET/DELETE /connectors/{name}` | Connector resource plus data source               | typed connector client/resource      |
| Connector config    | `GET/PUT /connectors/{name}/config`                     | Resource update/read support                      | masked-secret-aware config lifecycle |
| Connector status    | `GET /connectors/{name}/status`                         | Data source and polling helper                    | typed status data source             |
| Connector lifecycle | `POST restart`, `PUT pause/resume/stop`                 | Actions                                           | verify                               |
| Tasks               | `GET /connectors/{name}/tasks`, task status/restart     | Data sources/actions                              | verify                               |
| Topics              | `GET /connectors/{name}/topics`, `PUT .../topics/reset` | Data source/action                                | verify                               |
| Offsets             | `GET/PATCH/DELETE /connectors/{name}/offsets`           | Data source/actions with stopped-state validation | verify                               |
| Plugins             | `GET /connector-plugins`, `PUT .../config/validate`     | Data sources                                      | verify                               |
| Log levels          | Connect log-level endpoints where supported             | Action/client-only                                | verify                               |

Connect questions to resolve:

- Identify all sensitive configuration keys and preserve masked values without state churn.
- Define behavior for 409 rebalance responses and 202 asynchronous transitions.
- Confirm connector create/update response differences and import behavior.
- Confirm offset request shapes for source and sink connectors.

## Schema Registry

Reference: [Schema Registry API](https://docs.confluent.io/platform/current/schema-registry/develop/api.html)

| API family          | Representative endpoints                                     | Terraform exposure                           | Contract status |
| ------------------- | ------------------------------------------------------------ | -------------------------------------------- | --------------- |
| Config              | `GET/PUT/DELETE /config`, subject config variants            | Resource plus data source                    | verify          |
| Contexts            | `GET /contexts`                                              | Data source                                  | verify          |
| Exporters           | `/exporters` CRUD, status/config, pause/reset/resume         | Resource, data source, actions               | verify          |
| Importers           | `/importers` CRUD, status/config, pause/reset/resume         | Resource, data source, actions               | verify          |
| Modes               | `/mode` and subject/context variants                         | Resource plus data source                    | verify          |
| Schemas             | `/schemas`, `/schemas/ids`, `/schemas/guids`, types          | Data sources                                 | verify          |
| Subjects            | `/subjects`, subject versions, metadata, references          | Data sources                                 | verify          |
| Schema registration | `POST /subjects/{subject}/versions`                          | Resource with version/ID/GUID modeling       | verify          |
| Schema lookup       | `POST /subjects/{subject}`                                   | Data source or compatibility action          | verify          |
| Schema deletion     | `DELETE /subjects/{subject}` and versions                    | Resource delete with soft/hard-delete opt-in | verify          |
| Compatibility       | `POST /compatibility/subjects/...`                           | Data source/action                           | verify          |
| KEKs                | `/dek-registry/v1/keks` CRUD, test, undelete                 | Resource, data source, action                | verify          |
| DEKs                | `/dek-registry/v1/keks/{kek}/deks` CRUD/read/delete/undelete | Resource, data source, action                | verify          |
| Policy              | `GET /dek-registry/v1/policy`                                | Data source                                  | verify          |

Schema Registry questions to resolve:

- Confirm content negotiation and required `Confluent-Accept-Unknown-Properties` behavior.
- Separate schema ID, subject version, GUID, and context in state and imports.
- Define exact soft-delete versus permanent-delete lifecycle.
- Protect masked exporter/importer credentials and encrypted key material.
- Confirm which KEK/DEK endpoints are enabled by the deployment.

## ksqlDB

Reference: [ksqlDB REST API](https://docs.confluent.io/platform/current/ksqldb/developer-guide/ksqldb-rest-api/overview.html)

| API family             | Representative endpoints  | Terraform exposure                                   | Contract status |
| ---------------------- | ------------------------- | ---------------------------------------------------- | --------------- |
| Server info            | `GET /info`               | Data source                                          | verify          |
| Cluster status         | `GET /clusterStatus`      | Data source                                          | verify          |
| Execute statement      | `POST /ksql`              | Action; resource only for reconciled durable objects | verify          |
| Statement status       | `GET /status/{commandId}` | Data source/polling helper                           | verify          |
| Query                  | `POST /query`             | Action with bounded result capture                   | verify          |
| Push/pull query stream | `POST /query-stream`      | Action with cancellation and limits                  | verify          |
| Property validity      | `GET /is_valid_property`  | Data source                                          | verify          |
| Cluster termination    | `POST /terminate`         | Explicitly guarded action                            | verify          |

ksqlDB questions to resolve:

- Confirm media types and response framing for the target Platform version.
- Define query timeouts, cancellation, maximum rows/bytes, and state behavior.
- Never log raw SQL or unbounded streaming responses.
- Determine which SQL-created objects can be reconciled safely as Terraform resources.

## Cross-Service Authentication and Permissions

| Service         | Basic auth                          | Bearer/OAuth/OIDC          | TLS/mTLS                              | Permission review             |
| --------------- | ----------------------------------- | -------------------------- | ------------------------------------- | ----------------------------- |
| REST Proxy      | verify                              | verify                     | required deployment option            | Kafka ACL/RBAC                |
| MDS             | supported in documented deployments | supported where configured | supported, including mTLS deployments | MDS roles and ACLs            |
| Connect         | deployment-dependent                | deployment-dependent       | deployment-dependent                  | Connect REST authorization    |
| Schema Registry | deployment-dependent                | deployment-dependent       | deployment-dependent                  | Schema Registry authorization |
| ksqlDB          | deployment-dependent                | deployment-dependent       | deployment-dependent                  | ksqlDB/RBAC authorization     |

The provider must document configuration precedence, avoid logging secrets, and return actionable diagnostics for authentication and authorization failures.

## First Release Priorities

1. Kafka topics and topic configuration.
2. Kafka ACL bindings.
3. Kafka cluster links where the target version supports stable lifecycle operations.
4. Kafka Connect connector lifecycle.
5. Schema Registry subjects and schema registration.
6. Schema Registry compatibility and mode configuration.
7. MDS role bindings and cluster registry.
8. Read-only cluster, status, and introspection data sources.

## Phase 0 Exit Criteria

- [x] Supported Platform version range is approved as Confluent Platform 8.0-8.3.
- [ ] Every endpoint row has a verified path and API version.
- [ ] Every endpoint row has request/response models and permission requirements.
- [ ] Every endpoint row has pagination, retry, async, and deletion semantics recorded.
- [x] Every baseline endpoint family has a Terraform exposure decision.
- [x] First-release priorities are approved.

Phase 0 is complete for planning purposes. The remaining unchecked items are release-specific contract verification tasks that will be completed alongside each typed client and live 8.x acceptance test.
