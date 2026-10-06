provider "confluent-platform" {
  rest_proxy_endpoint = var.rest_proxy_endpoint
  username            = var.confluent_username
  password            = var.confluent_password
}
