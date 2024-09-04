resource "identitynow_source_aggregation_schedule" "test_account_schedule" {
  source_cloud_id  = identitynow_source.demo_source.cloud_external_id
  cron_expression  = "0 5 0 * * ?"
  aggregation_type = "account"
}
