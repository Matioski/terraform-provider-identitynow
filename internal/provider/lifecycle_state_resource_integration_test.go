//go:build integration

package provider

import (
	"testing"

	"github.com/sailpoint-oss/golang-sdk/v2/api_v2024"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func getAllowedSource(limit int32) []api_v2024.Source {
	return getSources(limit, "(features ca \"ENABLE\" or type eq \"DelimitedFile\" or type eq \"Delimited File\")")
}

func TestIntegration_ActiveLifecycleStateResource_ResourceRecreation(t *testing.T) {
	sourceClouds := getAllowedSource(2)
	sourceCloudId := *sourceClouds[0].Id
	updatedSourceCloudId := *sourceClouds[1].Id

	identityProfileId := IDENTITY_PROFILE

	accessProfiles := getAccessProfiles(2)
	accessProfileId := *accessProfiles[0].Id
	updatedAccessProfileId := *accessProfiles[1].Id

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerIntegrationConfig + `
resource "identitynow_lifecycle_state" "test" {
  name         = "TestState"
  technical_name = "testState"
  identity_profile_id = "` + identityProfileId + `"
  description = "my test description"
  enabled = false
  identity_state = "ACTIVE"
  email_notification_option = {
    notify_managers = false
    notify_all_admins = true
    notify_specific_users = true
    email_address_list = ["` + ownerIdentityEmail + `"]
  }
  account_actions = [
    {
      action = "ENABLE"
      source_ids = ["` + updatedSourceCloudId + `", "` + sourceCloudId + `"]
    }
  ]
  access_profile_ids = ["` + accessProfileId + `", "` + updatedAccessProfileId + `"]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("identitynow_lifecycle_state.test", "id"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "name", "TestState"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "technical_name", "testState"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "identity_profile_id", identityProfileId),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "description", "my test description"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "enabled", "false"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "identity_state", "ACTIVE"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "email_notification_option.notify_managers", "false"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "email_notification_option.notify_all_admins", "true"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "email_notification_option.notify_specific_users", "true"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "email_notification_option.email_address_list.0", ownerIdentityEmail),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "account_actions.0.action", "ENABLE"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "account_actions.0.source_ids.#", "2"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "access_profile_ids.0", accessProfileId),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "access_profile_ids.1", updatedAccessProfileId),
				),
			},
			// Update and Read testing
			{
				Config: providerIntegrationConfig + `
resource "identitynow_lifecycle_state" "test" {
  name         = "TestStateUpd"
  technical_name = "testStateUpd"
  identity_profile_id = "` + identityProfileId + `"
  description = "my test description"
  enabled = false
  identity_state = "ACTIVE"
  email_notification_option = {
    notify_managers = false
    notify_all_admins = true
    notify_specific_users = true
    email_address_list = ["` + ownerIdentityEmail + `"]
  }
  account_actions = [
    {
      action = "ENABLE"
      source_ids = ["` + updatedSourceCloudId + `"]
    }
  ]
  access_profile_ids = ["` + updatedAccessProfileId + `"]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("identitynow_lifecycle_state.test", "id"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "name", "TestStateUpd"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "technical_name", "testStateUpd"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "identity_profile_id", identityProfileId),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "description", "my test description"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "enabled", "false"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "identity_state", "ACTIVE"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "email_notification_option.notify_managers", "false"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "email_notification_option.notify_all_admins", "true"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "email_notification_option.notify_specific_users", "true"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "email_notification_option.email_address_list.0", ownerIdentityEmail),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "account_actions.0.action", "ENABLE"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "account_actions.0.source_ids.#", "1"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "account_actions.0.source_ids.0", updatedSourceCloudId),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "access_profile_ids.0", updatedAccessProfileId),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestIntegration_InactiveShortTermLifecycleStateResource_ResourcePatch(t *testing.T) {
	sourceClouds := getAllowedSource(2)
	sourceCloudId := *sourceClouds[0].Id
	updatedSourceCloudId := *sourceClouds[1].Id

	identityProfileId := IDENTITY_PROFILE

	accessProfiles := getAccessProfiles(2)
	accessProfileId := *accessProfiles[0].Id
	anotherAccessProfileId := *accessProfiles[1].Id

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerIntegrationConfig + `
resource "identitynow_lifecycle_state" "test" {
  name         = "Termination State 1"
  technical_name = "terminationState1"
  identity_profile_id = "` + identityProfileId + `"
  description = "my test description"
  enabled = false
  identity_state = "INACTIVE_SHORT_TERM"
  email_notification_option = {
    notify_managers = false
    notify_all_admins = true
    notify_specific_users = true
    email_address_list = ["` + ownerIdentityEmail + `"]
  }
  account_actions = [
    {
      action = "ENABLE"
      source_ids = ["` + sourceCloudId + `"]
    }
  ]
  access_profile_ids = ["` + accessProfileId + `", "` + anotherAccessProfileId + `"]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("identitynow_lifecycle_state.test", "id"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "name", "Termination State 1"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "technical_name", "terminationState1"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "identity_profile_id", identityProfileId),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "description", "my test description"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "enabled", "false"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "identity_state", "INACTIVE_SHORT_TERM"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "email_notification_option.notify_managers", "false"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "email_notification_option.notify_all_admins", "true"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "email_notification_option.notify_specific_users", "true"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "email_notification_option.email_address_list.0", ownerIdentityEmail),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "account_actions.0.action", "ENABLE"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "account_actions.0.source_ids.0", sourceCloudId),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "access_profile_ids.0", accessProfileId),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "access_profile_ids.1", anotherAccessProfileId),
				),
			},
			// Update and Read testing
			{
				Config: providerIntegrationConfig + `
resource "identitynow_lifecycle_state" "test" {
  name                = "Termination State 1"
  technical_name      = "terminationState1"
  identity_profile_id = "` + identityProfileId + `"
  description = "my test descriptionUpd"
  enabled = true
  identity_state = "INACTIVE_SHORT_TERM"
  email_notification_option = {
    notify_managers = true
    notify_all_admins = false
    notify_specific_users = true
    email_address_list = ["` + ownerIdentityEmail + `", "` + updatedOwnerIdentityEmail + `"]
  }
  account_actions = [
    {
      action = "DISABLE"
      source_ids = ["` + updatedSourceCloudId + `"]
    }
  ]
  access_profile_ids = ["` + anotherAccessProfileId + `"]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("identitynow_lifecycle_state.test", "id"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "name", "Termination State 1"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "technical_name", "terminationState1"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "identity_profile_id", identityProfileId),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "description", "my test descriptionUpd"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "enabled", "true"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "identity_state", "INACTIVE_SHORT_TERM"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "email_notification_option.notify_managers", "true"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "email_notification_option.notify_all_admins", "false"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "email_notification_option.notify_specific_users", "true"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "email_notification_option.email_address_list.0", ownerIdentityEmail),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "email_notification_option.email_address_list.1", updatedOwnerIdentityEmail),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "account_actions.0.action", "DISABLE"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "account_actions.0.source_ids.0", updatedSourceCloudId),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "access_profile_ids.0", anotherAccessProfileId),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestIntegration_InactiveLongTermLifecycleStateResource_ResourcePatch(t *testing.T) {
	sourceClouds := getAllowedSource(2)
	sourceCloudId := *sourceClouds[0].Id
	updatedSourceCloudId := *sourceClouds[1].Id

	identityProfileId := IDENTITY_PROFILE

	accessProfiles := getAccessProfiles(2)
	accessProfileId := *accessProfiles[0].Id
	anotherAccessProfileId := *accessProfiles[1].Id

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerIntegrationConfig + `
resource "identitynow_lifecycle_state" "test" {
  name         = "Termination State X"
  technical_name = "terminationStateX"
  identity_profile_id = "` + identityProfileId + `"
  description = "my test description"
  enabled = false
  identity_state = "INACTIVE_LONG_TERM"
  email_notification_option = {
    notify_managers = false
    notify_all_admins = true
    notify_specific_users = true
    email_address_list = ["` + ownerIdentityEmail + `"]
  }
  account_actions = [
    {
      action = "ENABLE"
      source_ids = ["` + sourceCloudId + `"]
    }
  ]
  access_profile_ids = ["` + accessProfileId + `", "` + anotherAccessProfileId + `"]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("identitynow_lifecycle_state.test", "id"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "name", "Termination State X"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "technical_name", "terminationStateX"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "identity_profile_id", identityProfileId),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "description", "my test description"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "enabled", "false"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "identity_state", "INACTIVE_LONG_TERM"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "email_notification_option.notify_managers", "false"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "email_notification_option.notify_all_admins", "true"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "email_notification_option.notify_specific_users", "true"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "email_notification_option.email_address_list.0", ownerIdentityEmail),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "account_actions.0.action", "ENABLE"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "account_actions.0.source_ids.0", sourceCloudId),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "access_profile_ids.0", accessProfileId),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "access_profile_ids.1", anotherAccessProfileId),
				),
			},
			// Update and Read testing
			{
				Config: providerIntegrationConfig + `
resource "identitynow_lifecycle_state" "test" {
  name                = "Termination State X"
  technical_name      = "terminationStateX"
  identity_profile_id = "` + identityProfileId + `"
  description = "my test descriptionUpd"
  enabled = true
  identity_state = "INACTIVE_LONG_TERM"
  email_notification_option = {
    notify_managers = true
    notify_all_admins = false
    notify_specific_users = true
    email_address_list = ["` + updatedOwnerIdentityEmail + `"]
  }
  account_actions = [
    {
      action = "DISABLE"
      source_ids = ["` + updatedSourceCloudId + `"]
    }
  ]
  access_profile_ids = ["` + anotherAccessProfileId + `"]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("identitynow_lifecycle_state.test", "id"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "name", "Termination State X"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "technical_name", "terminationStateX"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "identity_profile_id", identityProfileId),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "description", "my test descriptionUpd"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "enabled", "true"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "identity_state", "INACTIVE_LONG_TERM"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "email_notification_option.notify_managers", "true"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "email_notification_option.notify_all_admins", "false"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "email_notification_option.notify_specific_users", "true"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "email_notification_option.email_address_list.0", updatedOwnerIdentityEmail),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "account_actions.0.action", "DISABLE"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "account_actions.0.source_ids.0", updatedSourceCloudId),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "access_profile_ids.0", anotherAccessProfileId),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
