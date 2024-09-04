//go:build integration

package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	sailpointBeta "github.com/sailpoint-oss/golang-sdk/v2/api_beta"
	"terraform-provider-identitynow/internal/util"
	"testing"
)

func TestIntegration_LifecycleStateResource_ResourceRecreation(t *testing.T) {
	t.Skip("1949455")

	sourceClouds := getSources(2, context.Background())
	sourceCloudId := *sourceClouds[0].Id
	updatedSourceCloudId := *sourceClouds[1].Id

	identityProfiles := getIdentityProfile(2, context.Background())
	identityProfileId := *identityProfiles[0].Id
	updatedIdentityProfileId := *identityProfiles[1].Id

	accessProfiles := getAccessProfiles(2, context.Background())
	accessProfileId := *accessProfiles[0].Id
	updatedAccessProfileId := *accessProfiles[1].Id

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerIntegrationConfig + `
resource "identitynow_lifecycle_state" "test" {
  name         = "Active"
  technical_name = "active"
  identity_profile_id = "` + identityProfileId + `"
  description = "my test description"
  enabled = false
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
  access_profile_ids = ["` + accessProfileId + `", "` + updatedAccessProfileId + `"]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("identitynow_lifecycle_state.test", "id"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "name", "Active"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "technical_name", "active"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "identity_profile_id", identityProfileId),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "description", "my test description"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "enabled", "false"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "email_notification_option.notify_managers", "false"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "email_notification_option.notify_all_admins", "true"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "email_notification_option.notify_specific_users", "true"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "email_notification_option.email_address_list.0", ownerIdentityEmail),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "account_actions.0.action", "ENABLE"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "account_actions.0.source_ids.0", sourceCloudId),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "access_profile_ids.0", accessProfileId),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "access_profile_ids.1", updatedAccessProfileId),
				),
			},
			// Update and Read testing
			{
				Config: providerIntegrationConfig + `
resource "identitynow_lifecycle_state" "test" {
  name         = "ActiveUpd"
  technical_name = "activeUpd"
  identity_profile_id = "` + updatedIdentityProfileId + `"
  description = "my test description"
  enabled = false
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
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "name", "ActiveUpd"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "technical_name", "activeUpd"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "identity_profile_id", updatedIdentityProfileId),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "description", "my test description"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "enabled", "false"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "email_notification_option.notify_managers", "false"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "email_notification_option.notify_all_admins", "true"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "email_notification_option.notify_specific_users", "true"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "email_notification_option.email_address_list.0", ownerIdentityEmail),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "account_actions.0.action", "ENABLE"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "account_actions.0.source_ids.0", updatedSourceCloudId),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "access_profile_ids.0", updatedAccessProfileId),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestIntegration_LifecycleStateResource_ResourcePatch(t *testing.T) {
	t.Skip("1949455")

	sourceClouds := getSources(2, context.Background())
	sourceCloudId := *sourceClouds[0].Id
	updatedSourceCloudId := *sourceClouds[1].Id

	identityProfileId := *getIdentityProfile(1, context.Background())[0].Id

	accessProfiles := getAccessProfiles(2, context.Background())
	accessProfileId := *accessProfiles[0].Id
	anotherAccessProfileId := *accessProfiles[1].Id

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerIntegrationConfig + `
resource "identitynow_lifecycle_state" "test" {
  name         = "ActiveX"
  technical_name = "activeX"
  identity_profile_id = "` + identityProfileId + `"
  description = "my test description"
  enabled = false
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
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "name", "ActiveX"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "technical_name", "activeX"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "identity_profile_id", identityProfileId),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "description", "my test description"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "enabled", "false"),
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
  name                = "ActiveX"
  technical_name      = "activeX"
  identity_profile_id = "` + identityProfileId + `"
  description = "my test descriptionUpd"
  enabled = true
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
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "name", "ActiveX"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "technical_name", "activeX"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "identity_profile_id", identityProfileId),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "description", "my test descriptionUpd"),
					resource.TestCheckResourceAttr("identitynow_lifecycle_state.test", "enabled", "true"),
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

func getIdentityProfile(limit int32, ctx context.Context) []sailpointBeta.IdentityProfile {
	identityProfiles, spResp, err := SPApiClient.Beta.IdentityProfilesAPI.ListIdentityProfiles(context.Background()).Limit(limit).Execute()
	if err != nil {
		fmt.Errorf("Error fetching identity profiles: %s\n%s", err, util.GetBody(spResp))
	}
	if len(identityProfiles) < int(limit) {
		fmt.Printf("Unable to provide %d identity profile(s); %d found\n", limit, len(identityProfiles))
	}

	return identityProfiles
}
