---
page_title: "Kafka Topic Resource - confluent-platform"
subcategory: "REST Proxy"
description: |-
  Manages a Kafka topic through Confluent Platform REST Proxy v3.
---

# confluent-platform_kafka_topic (Resource)

Creates and manages a Kafka topic. Topic names, partition count, and replication factor are immutable. Topic configuration values are reconciled through the nested `configs` map.

## Example Usage

```terraform
resource "confluent-platform_kafka_topic" "orders" {
  cluster_id         = "lkc-123"
  name               = "orders"
  partitions_count   = 3
  replication_factor = 2

  configs = {
    "cleanup.policy" = "compact"
  }
}
```

## Schema

### Required

- `cluster_id` (String) Kafka cluster ID.
- `name` (String) Topic name.

### Optional

- `configs` (Map of String) Topic configuration values.
- `partitions_count` (Number) Number of partitions.
- `replication_factor` (Number) Replication factor.

### Read-Only

- `id` (String) Composite cluster/topic identity.
- `is_internal` (Boolean) Whether the topic is internal.

## Import

```shell
terraform import confluent-platform_kafka_topic.orders 'lkc-123/orders'
```
