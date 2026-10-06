---
page_title: "MDS Active Nodes Data Source - confluent-platform"
subcategory: "Metadata Service"
description: |-
  Reads active Metadata Service nodes.
---

# confluent-platform_mds_active_nodes (Data Source)

Reads the active-node response from the Metadata Service. The response is retained as JSON because node fields can vary with deployment and Platform version.

## Example Usage

```terraform
data "confluent-platform_mds_active_nodes" "https" {
  protocol = "https"
}
```

## Schema

### Required

- `protocol` (String) MDS transport protocol, such as `https`.

### Read-Only

- `id` (String) Active-node lookup identity.
- `response_json` (String) Raw active-node response from MDS.
