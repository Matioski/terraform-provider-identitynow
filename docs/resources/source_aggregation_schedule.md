---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "identitynow_source_aggregation_schedule Resource - terraform-provider-identitynow"
subcategory: ""
description: |-
  
---

# identitynow_source_aggregation_schedule (Resource)



## Example Usage

```terraform
resource "identitynow_source_aggregation_schedule" "test_account_schedule" {
  source_cloud_id  = identitynow_source.demo_source.cloud_external_id
  cron_expression  = "0 5 0 * * ?"
  aggregation_type = "account"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `aggregation_type` (String) Aggregation type one of 'account' or 'entitlement'
- `cron_expression` (String) Cron Expression for the Schedule
- `source_cloud_id` (String) Legacy Source ID
