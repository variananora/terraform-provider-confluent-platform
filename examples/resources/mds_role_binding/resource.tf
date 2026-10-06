resource "confluent-platform_mds_role_binding" "orders_reader" {
  principal        = "User:alice"
  role_name        = "DeveloperRead"
  kafka_cluster_id = var.kafka_cluster_id
  resource_type    = "Topic"
  resource_name    = "orders"
  pattern_type     = "LITERAL"
}

resource "confluent-platform_mds_role_binding" "connect_admin" {
  principal        = "User:alice"
  role_name        = "CloudClusterAdmin"
  kafka_cluster_id = var.kafka_cluster_id
}