---
page_title: "MDS Roles Data Source - confluent-platform"
subcategory: "Metadata Service"
description: "Reads MDS role definitions as JSON."
---

# confluent-platform_mds_roles (Data Source)

```terraform
data "confluent-platform_mds_roles" "all" {}
```

Set `role_name` to read one role. The `response_json` attribute contains the response.
