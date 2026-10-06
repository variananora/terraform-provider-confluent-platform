resource "confluent-platform_kafka_acl" "orders_reader" {
  cluster_id    = var.kafka_cluster_id
  resource_type = "TOPIC"
  resource_name = "orders"
  pattern_type  = "LITERAL"
  principal     = "User:alice"
  host          = "*"
  operation     = "READ"
  permission    = "ALLOW"
}