---
page_title: "MDS Metadata Cluster ID Data Source - confluent-platform"
subcategory: "Metadata Service"
description: |-
  Reads the Metadata Service cluster ID.
---

# confluent-platform_mds_metadata_cluster_id (Data Source)

Reads the MDS metadata cluster ID and exposes the response as JSON for compatibility across Platform deployments.

## Example Usage

```terraform
data "confluent-platform_mds_metadata_cluster_id" "current" {}
```

## Schema

### Read-Only

- `id` (String) Stable lookup identity.
- `response_json` (String) Raw metadata cluster ID response from MDS.
