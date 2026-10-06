---
page_title: "MDS Cluster Registry Resource - confluent-platform"
subcategory: "Metadata Service"
description: "Manages a named MDS cluster registry entry."
---

# confluent-platform_mds_cluster_registry (Resource)

```terraform
resource "confluent-platform_mds_cluster_registry" "kafka" {
  cluster_name = "primary-kafka"
  configuration_json = jsonencode({
    clusterName = "primary-kafka"
    scope = { clusters = { "kafka-cluster" = var.kafka_cluster_id } }
    hosts = [{ host = "kafka.example.com", port = 9092 }]
    protocol = "SASL_SSL"
  })
}
```

`configuration_json` is preserved as deployment-specific JSON.
