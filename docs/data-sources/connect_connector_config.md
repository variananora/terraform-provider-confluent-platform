---
page_title: "Kafka Connect Connector Config Data Source - confluent-platform"
subcategory: "Kafka Connect"
description: "Reads Kafka Connect connector configuration without recovering masked values."
---

# confluent-platform_connect_connector_config (Data Source)

```terraform
data "confluent-platform_connect_connector_config" "orders" {
  name = "orders"
}
```

Masked secret values remain masked in the JSON response.
