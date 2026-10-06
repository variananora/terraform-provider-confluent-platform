---
page_title: "MDS Audit Routes Data Source - confluent-platform"
subcategory: "Metadata Service"
description: "Reads MDS audit routes matching a CRN query."
---

# confluent-platform_mds_audit_routes (Data Source)

```terraform
data "confluent-platform_mds_audit_routes" "topic" {
  query = "crn://mds.example.com/kafka=*/topic=orders"
}
```
