//go:build integration

package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func checkIdentityProfileIsDeleted() func(state *terraform.State) error {
	return func(state *terraform.State) error {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "identitynow_identity_profile" {
				continue
			}
			id := rs.Primary.ID
			_, response, _ := SPApiClient.V3.IdentityProfilesAPI.GetIdentityProfile(context.Background(), id).Execute()
			if response != nil && response.StatusCode == 404 {
				return nil
			}
			SPApiClient.V3.IdentityProfilesAPI.DeleteIdentityProfile(context.Background(), id).Execute()
			return fmt.Errorf("identity profile still exists: %s", id)
		}
		return nil
	}
}

func TestIntegration_IdentityProfileResource(t *testing.T) {
	//auth source has to be hardcoded since it should not be attached to an identity profile and there's a need to know which attributes are available in the source
	authSource := "760ac507e72a4fb7ade1ffa7ddeca00d"
	authSourceName := "Terraform Integration Source"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		CheckDestroy:             checkIdentityProfileIsDeleted(),
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerIntegrationConfig + `
resource "identitynow_identity_profile" "test" {
  name        = "Test Identity Profile"
  description = "Test Identity Profile Description"
  priority = 102
  owner = {
    id   = "` + ownerIdentityId + `"
    name   = "` + ownerIdentityName + `"
  }
  authoritative_source = {
    id = "` + authSource + `"
    #name = "` + authSourceName + `"
  }
  identity_attribute_config = {
    enabled = true
    attribute_transforms = [
      {
        identity_attribute_name = "name"
        transform_definition = {
          type = "accountAttribute"
          attributes = jsonencode({
            sourceName    = "Terraform Integration Source"
            attributeName = "name"
            sourceId      = "` + authSource + `"
          })
        }
      }
    ]
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("identitynow_identity_profile.test", "id"),
					resource.TestCheckResourceAttr("identitynow_identity_profile.test", "name", "Test Identity Profile"),
					resource.TestCheckResourceAttr("identitynow_identity_profile.test", "description", "Test Identity Profile Description"),
					resource.TestCheckResourceAttr("identitynow_identity_profile.test", "priority", "102"),
					resource.TestCheckResourceAttr("identitynow_identity_profile.test", "owner.id", ownerIdentityId),
					resource.TestCheckResourceAttr("identitynow_identity_profile.test", "authoritative_source.id", authSource),
					resource.TestCheckResourceAttr("identitynow_identity_profile.test", "identity_attribute_config.enabled", "true"),
					resource.TestCheckResourceAttr("identitynow_identity_profile.test", "identity_attribute_config.attribute_transforms.0.identity_attribute_name", "name"),
					resource.TestCheckResourceAttr("identitynow_identity_profile.test", "identity_attribute_config.attribute_transforms.0.transform_definition.type", "accountAttribute"),
					resource.TestCheckResourceAttr("identitynow_identity_profile.test", "identity_attribute_config.attribute_transforms.0.transform_definition.attributes", "{\"attributeName\":\"name\",\"sourceId\":\"760ac507e72a4fb7ade1ffa7ddeca00d\",\"sourceName\":\"Terraform Integration Source\"}"),
				),
			},
			// Update and Read testing
			{
				Config: providerIntegrationConfig + `
resource "identitynow_identity_profile" "test" {
  name        = "Test Identity Profile UPDATED"
  description = "Test Identity Profile Description UPDATED"
  priority = 107
  owner = {
    id   = "` + updatedOwnerIdentityId + `"
    name   = "` + updatedOwnerIdentityName + `"
  }
  authoritative_source = {
    id = "` + authSource + `"
  }
 identity_attribute_config = {
    enabled = false
    attribute_transforms = [
      {
        identity_attribute_name = "name"
        transform_definition = {
          type = "accountAttribute"
          attributes = jsonencode({
            sourceName    = "Terraform Integration Source"
            attributeName = "name"
            sourceId      = "` + authSource + `"
          })
        }
      }
    ]
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("identitynow_identity_profile.test", "id"),
					resource.TestCheckResourceAttr("identitynow_identity_profile.test", "name", "Test Identity Profile UPDATED"),
					resource.TestCheckResourceAttr("identitynow_identity_profile.test", "description", "Test Identity Profile Description UPDATED"),
					resource.TestCheckResourceAttr("identitynow_identity_profile.test", "priority", "107"),
					resource.TestCheckResourceAttr("identitynow_identity_profile.test", "owner.id", updatedOwnerIdentityId),
					resource.TestCheckResourceAttr("identitynow_identity_profile.test", "authoritative_source.id", authSource),
					resource.TestCheckResourceAttr("identitynow_identity_profile.test", "identity_attribute_config.enabled", "false"),
					resource.TestCheckResourceAttr("identitynow_identity_profile.test", "identity_attribute_config.attribute_transforms.0.identity_attribute_name", "name"),
					resource.TestCheckResourceAttr("identitynow_identity_profile.test", "identity_attribute_config.attribute_transforms.0.transform_definition.type", "accountAttribute"),
					resource.TestCheckResourceAttr("identitynow_identity_profile.test", "identity_attribute_config.attribute_transforms.0.transform_definition.attributes", "{\"attributeName\":\"name\",\"sourceId\":\"760ac507e72a4fb7ade1ffa7ddeca00d\",\"sourceName\":\"Terraform Integration Source\"}"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
