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