---
page_title: "MDS Role Binding Summary Data Source - confluent-platform"
subcategory: "Metadata Service"
description: "Reads effective MDS role-binding resources for a principal."
---

# confluent-platform_mds_role_binding_summary (Data Source)

```terraform
data "confluent-platform_mds_role_binding_summary" "alice" {
  principal        = "User:alice"
  kafka_cluster_id = var.kafka_cluster_id
}
```
