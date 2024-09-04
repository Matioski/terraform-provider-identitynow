//go:build !integration

package provider

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"
)

func TestLifecycleStateResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "identitynow_lifecycle_state" "test" {
  name         = "Active"
  technical_name = "active"
  identity_profile_id = "profileId"
  description = "my test description"
  enabled = false
  email_notification_option = {
    notify_managers = false
    notify_all_admins = true
    notify_specific_users = true
    email_address_list = ["email1", "email2"]
  }
  account_actions = [
    {
      action = "ENABLE"
      source_ids = ["sourceId1", "sourceId2"]
    },
    {
      action = "DISABLE"
      source_ids = ["sourceId3", "sourceId4"]
    }
  ]
  access_profile_ids = ["accessProfileId1", "accessProfileId2"]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("identitynow_lifecycle_state.test", "id"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "name", "Active"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "technical_name", "active"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "identity_profile_id", "profileId"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "description", "my test description"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "enabled", "false"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "email_notification_option.notify_managers", "false"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "email_notification_option.notify_all_admins", "true"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "email_notification_option.notify_specific_users", "true"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "email_notification_option.email_address_list.0", "email1"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "email_notification_option.email_address_list.1", "email2"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "account_actions.0.action", "ENABLE"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "account_actions.0.source_ids.0", "sourceId1"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "account_actions.0.source_ids.1", "sourceId2"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "account_actions.1.action", "DISABLE"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "account_actions.1.source_ids.0", "sourceId3"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "account_actions.1.source_ids.1", "sourceId4"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "access_profile_ids.0", "accessProfileId1"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "access_profile_ids.1", "accessProfileId2"),
				),
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "identitynow_lifecycle_state" "test" {
  name                = "Active"
  technical_name      = "active"
  identity_profile_id = "profileId"
  description         = "my test descriptionUpd"
  enabled             = true
  email_notification_option = {
    notify_managers       = true
    notify_all_admins     = false
    notify_specific_users = false
    email_address_list    = ["email1"]
  }
  account_actions = [
    {
      action     = "ENABLE"
      source_ids = ["sourceId1"]
    }
  ]
  access_profile_ids = ["accessProfileId1", "accessProfileId2", "accessProfileId3"]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("identitynow_lifecycle_state.test", "id"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "name", "Active"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "technical_name", "active"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "identity_profile_id", "profileId"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "description", "my test descriptionUpd"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "enabled", "true"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "email_notification_option.notify_managers", "true"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "email_notification_option.notify_all_admins", "false"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "email_notification_option.notify_specific_users", "false"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "email_notification_option.email_address_list.0", "email1"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "account_actions.0.action", "ENABLE"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "account_actions.0.source_ids.0", "sourceId1"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "access_profile_ids.0", "accessProfileId1"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "access_profile_ids.1", "accessProfileId2"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "access_profile_ids.2", "accessProfileId3"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
