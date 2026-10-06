---
page_title: "MDS Role Binding Resource - confluent-platform"
subcategory: "Metadata Service"
description: |-
  Manages an MDS RBAC role binding for a Kafka cluster, topic, or consumer group.
---

# confluent-platform_mds_role_binding (Resource)

Use `resource_type = "Topic"` or `resource_type = "Group"` for resource-scoped RBAC. Omit the resource fields for cluster-scoped roles, including roles for Connect, Schema Registry, ksqlDB, and other cluster-level services.

## Topic or Group Example

```terraform
resource "confluent-platform_mds_role_binding" "orders_reader" {
  principal        = "User:alice"
  role_name        = "DeveloperRead"
  kafka_cluster_id = var.kafka_cluster_id
  resource_type    = "Topic"
  resource_name    = "orders"
  pattern_type     = "LITERAL"
}
```

## Cluster-Scoped Example

```terraform
resource "confluent-platform_mds_role_binding" "connect_admin" {
  principal        = "User:alice"
  role_name        = "CloudClusterAdmin"
  kafka_cluster_id = var.kafka_cluster_id
}
```

## Schema

### Required

- `principal` (String) Kafka principal.
- `role_name` (String) MDS RBAC role name.
- `kafka_cluster_id` (String) Kafka cluster ID in the MDS scope.

### Optional

- `resource_type` (String) `Topic` or `Group`; omit for cluster-scoped roles.
- `resource_name` (String) Topic or consumer group name.
- `pattern_type` (String) `LITERAL` or `PREFIXED`.

### Read-Only

- `id` (String) Stable composite binding identity.

## Import

Cluster scope:

```shell
terraform import confluent-platform_mds_role_binding.connect_admin \
  'User:alice|CloudClusterAdmin|kafka-cluster-id'
```

Topic or group scope:

```shell
terraform import confluent-platform_mds_role_binding.orders_reader \
  'User:alice|DeveloperRead|kafka-cluster-id|Topic|orders|LITERAL'
```
