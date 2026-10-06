We are going to build a Terraform provider for Confluent Platform.
Confluent Platorm is a version of Confluent that runs on premise rather in the cloud

We are going to be referencing some things from the cloud version here:  
https://github.com/confluentinc/terraform-provider-confluent

Confluent Platform have a couple of REST APIs that we can use, lets break this project to those APIs.
And in each of this API lets define our requirements.

I want to have datasource and resources for each APIs.

The basic point is 
Convert all the GET API as Data source
Convert all the POST API as Resource


Confluent REST Proxy
https://docs.confluent.io/platform/current/kafka-rest/api.html

API List:
- Cell (v3)
- Cluster (v3)
- Configs (v3)
- Broker (v3)
- Replica (v3)
- ACL (v3)
- Consumer Group (v3)
- Partition (v3)
- Topic (v3)
- Records (v3)
- Cluster Linking (v3)
- Share Group (v3)
- BalancerStatus (v3)
- BrokerTask (v3)
- BrokerReplicaExclusion (v3)
- RemoveBrokerTask (v3)
- Unregister (v3)
- License (v3)
- Replica Status (v3)
- Streams Group (v3)

Metadata API
https://docs.confluent.io/platform/current/security/authorization/rbac/mds-api.html#confluent-metadata-api-reference-for-cp

API List:
- Tokens and Authentication
- Authorization
- Metadata Service Operations
- RBAC - Role Definitions
- RBAC - RoleBinding CRUD
- RBAC - RoleBinding Summaries
- Kafka ACL Management
- Cluster Registry
- Audit Log Configuration
- Private RBAC UI - Cluster Visibility
- Private RBAC UI - My RoleBindings
- Private RBAC UI - Manage RoleBindings
- Private RBAC UI - Creation Guidelines
- Private RBAC UI - Cached User Store Information
- SSO - Device Authorization
- Sanitize token

Connect REST API
https://docs.confluent.io/platform/current/connect/references/restapi.html#connect-userguide-rest

- Status and Errors
- Connect Cluster
- Connectors
- Tasks
- Topics
- Offsets
  - Alter / reset offset responses

Schema Registry API
https://docs.confluent.io/platform/current/schema-registry/develop/api.html#schemaregistry-api

- Config
- Contexts
- Exporters
- Importers
- Modes
- Schemas
- Subjects
- Key Encryption Keys
- Data Encryption Keys

ksqlDB REST API
https://docs.confluent.io/platform/current/ksqldb/developer-guide/ksqldb-rest-api/overview.html#ksqldb-rest-api

- Execute a statement
- Run a query
- Run push and pull queries
- Terminate a cluster
- Introspect query status
- Introspect server status
- Introspect cluster status
- Get the validity of a property

Flink REST API
Lets skip Flink for now because we are not using Flink at the moment.
