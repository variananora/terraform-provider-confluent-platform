---
page_title: "Kafka Connect Connector Resource - confluent-platform"
subcategory: "Kafka Connect"
description: "Manages a Kafka Connect connector."
---

# confluent-platform_connect_connector (Resource)

Creates connectors with `POST /connectors` and updates configuration with `PUT /connectors/{name}/config`.

```terraform
resource "confluent-platform_connect_connector" "orders" {
  name = "orders"

  config = {
    "connector.class" = "io.confluent.connect.jdbc.JdbcSinkConnector"
    "topics"          = "orders"
    "connection.url"  = var.jdbc_url
  }
}
```

Connector secrets are marked sensitive. Connect values returned as masked placeholders are retained from prior Terraform state and are never recovered or logged.

## Import

```shell
terraform import confluent-platform_connect_connector.orders orders
```
