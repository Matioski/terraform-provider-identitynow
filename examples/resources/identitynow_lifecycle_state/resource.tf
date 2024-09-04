resource "identitynow_lifecycle_state" "demo_active" {
  name                = "Active"
  technical_name      = "active"
  identity_profile_id = "profileId"
  description         = "my test description"
  enabled             = false
  email_notification_option = {
    notify_managers       = false
    notify_all_admins     = true
    notify_specific_users = true
    email_address_list    = ["email1", "email2"]
  }
  account_actions = [
    {
      action     = "ENABLE"
      source_ids = ["sourceId1", "sourceId2"]
    },
    {
      action     = "DISABLE"
      source_ids = ["sourceId3", "sourceId4"]
    }
  ]
  access_profile_ids = ["accessProfileId1", "accessProfileId2"]
}
