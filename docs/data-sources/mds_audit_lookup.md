---
page_title: "MDS Audit Lookup Data Source - confluent-platform"
subcategory: "Metadata Service"
description: "Resolves the MDS audit route for a CRN."
---

# confluent-platform_mds_audit_lookup (Data Source)

```terraform
data "confluent-platform_mds_audit_lookup" "topic" {
  crn = "crn://mds.example.com/kafka=cluster-id/topic=orders"
}
```
