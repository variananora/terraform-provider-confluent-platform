---
page_title: "Kafka ACL Resource - confluent-platform"
subcategory: "REST Proxy"
description: |-
  Manages a Kafka ACL binding through Confluent Platform REST Proxy v3.
---

# confluent-platform_kafka_acl (Resource)

Manages one Kafka ACL binding. All identity fields are immutable, so changing any field replaces the binding.

## Example Usage

```terraform
resource "confluent-platform_kafka_acl" "orders_reader" {
  cluster_id    = "lkc-123"
  resource_type = "TOPIC"
  resource_name = "orders"
  pattern_type  = "LITERAL"
  principal     = "User:alice"
  host          = "*"
  operation     = "READ"
  permission    = "ALLOW"
}
```

## Schema

### Required

- `cluster_id` (String) Kafka cluster ID.
- `resource_type` (String) ACL resource type, such as `TOPIC`.
- `resource_name` (String) ACL resource name.
- `pattern_type` (String) ACL pattern type, such as `LITERAL`.
- `principal` (String) Kafka principal.
- `host` (String) ACL host, commonly `*`.
- `operation` (String) Kafka operation, such as `READ`.
- `permission` (String) ACL permission, such as `ALLOW`.

### Read-Only

- `id` (String) Stable composite ACL identity.

## Import

Import uses the following format:

```shell
terraform import confluent-platform_kafka_acl.orders_reader \
  'lkc-123|TOPIC|orders|LITERAL|User:alice|*|READ|ALLOW'
```

The importing principal must have permission to read and manage ACL bindings through REST Proxy.
