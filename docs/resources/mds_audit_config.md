---
page_title: "MDS Audit Configuration Resource - confluent-platform"
subcategory: "Metadata Service"
description: "Manages the MDS audit configuration."
---

# confluent-platform_mds_audit_config (Resource)

Supply the full current audit configuration, including `metadata.resource_version`:

```terraform
resource "confluent-platform_mds_audit_config" "current" {
  configuration_json = file("audit-config.json")
}
```

Concurrent updates return a diagnostic requiring the latest resource version to be merged before retrying.
