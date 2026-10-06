---
page_title: "Kafka Connect Connector Status Data Source - confluent-platform"
subcategory: "Kafka Connect"
description: "Reads Kafka Connect connector status."
---

# confluent-platform_connect_connector_status (Data Source)

```terraform
data "confluent-platform_connect_connector_status" "orders" {
  name = "orders"
}
```

The data source exposes the connector state and the response as JSON.
