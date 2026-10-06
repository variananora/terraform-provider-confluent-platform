resource "confluent-platform_kafka_topic" "orders" {
  cluster_id         = var.kafka_cluster_id
  name               = "orders"
  partitions_count   = 3
  replication_factor = 2

  configs = {
    "cleanup.policy" = "compact"
  }
}