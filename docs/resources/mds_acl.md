---
page_title: "MDS ACL Resource - confluent-platform"
subcategory: "Metadata Service"
description: |-
  Manages a centralized Kafka ACL through MDS.
---

# confluent-platform_mds_acl (Resource)

Manages one centralized Kafka ACL binding in MDS. The ACL scope always targets the Kafka cluster; the resource identity can be a Topic, Group, or another Kafka resource supported by the deployment.

```terraform
resource "confluent-platform_mds_acl" "orders_reader" {
  kafka_cluster_id = var.kafka_cluster_id
  resource_type    = "Topic"
  resource_name    = "orders"
  pattern_type     = "LITERAL"
  principal        = "User:alice"
  host             = "*"
  operation        = "Read"
  permission_type  = "Allow"
}
```

All identity fields are immutable and changing one replaces the ACL.

## Import

```shell
terraform import confluent-platform_mds_acl.orders_reader \
  'kafka-cluster-id|Topic|orders|LITERAL|User:alice|*|Read|Allow'
```
