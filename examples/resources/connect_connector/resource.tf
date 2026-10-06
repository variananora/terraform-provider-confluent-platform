resource "confluent-platform_connect_connector" "orders" {
  name = "orders"

  config = {
    "connector.class" = "io.confluent.connect.jdbc.JdbcSinkConnector"
    "topics"          = "orders"
    "connection.url"  = var.jdbc_url
  }
}