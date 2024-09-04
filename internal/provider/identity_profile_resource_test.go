//go:build !integration

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestIdentityProfileResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "identitynow_identity_profile" "test" {
  name        = "Test Identity Profile"
  description = "Test Identity Profile"
  priority = 5
  owner = {
    id   = "ownerId"
  }
  authoritative_source = {
    id = "auth_source"
  }
  identity_attribute_config = {
    enabled = true
    attribute_transforms = [
      {
        identity_attribute_name = "uid"
        transform_definition = {
          type = "accountAttribute"
          attributes = jsonencode({
            sourceName    = "sourceName"
            attributeName = "uid"
            sourceId      = "sourceId"
          })
        }
      },
      {
        identity_attribute_name = "email"
        transform_definition = {
          type = "accountAttribute"
          attributes = jsonencode({
            sourceName    = "sourceName"
            attributeName = "mailaddr"
            sourceId      = "sourceId"
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
					resource.TestCheckResourceAttr("identitynow_identity_profile.test", "description", "Test Identity Profile"),
					resource.TestCheckResourceAttr("identitynow_identity_profile.test", "priority", "5"),
					resource.TestCheckResourceAttr("identitynow_identity_profile.test", "owner.id", "ownerId"),
					resource.TestCheckResourceAttr("identitynow_identity_profile.test", "authoritative_source.id", "auth_source"),
					resource.TestCheckResourceAttr("identitynow_identity_profile.test", "identity_attribute_config.enabled", "true"),
					resource.TestCheckResourceAttr("identitynow_identity_profile.test", "identity_attribute_config.attribute_transforms.0.identity_attribute_name", "uid"),
					resource.TestCheckResourceAttr("identitynow_identity_profile.test", "identity_attribute_config.attribute_transforms.0.transform_definition.type", "accountAttribute"),
					resource.TestCheckResourceAttr("identitynow_identity_profile.test", "identity_attribute_config.attribute_transforms.0.transform_definition.attributes", "{\"attributeName\":\"uid\",\"sourceId\":\"sourceId\",\"sourceName\":\"sourceName\"}"),
					resource.TestCheckResourceAttr("identitynow_identity_profile.test", "identity_attribute_config.attribute_transforms.1.identity_attribute_name", "email"),
					resource.TestCheckResourceAttr("identitynow_identity_profile.test", "identity_attribute_config.attribute_transforms.1.transform_definition.type", "accountAttribute"),
					resource.TestCheckResourceAttr("identitynow_identity_profile.test", "identity_attribute_config.attribute_transforms.1.transform_definition.attributes", "{\"attributeName\":\"mailaddr\",\"sourceId\":\"sourceId\",\"sourceName\":\"sourceName\"}"),
				),
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "identitynow_identity_profile" "test" {
  name        = "Test Identity Profile UPDATED"
  description = "Test Identity Profile UPDATED"
  priority = 10
  owner = {
    id   = "ownerIdUpd"
  }
  authoritative_source = {
    id = "auth_sourceUpd"
  }
 identity_attribute_config = {
    enabled = false
    attribute_transforms = [
      {
        identity_attribute_name = "uid"
        transform_definition = {
          type = "accountAttribute"
          attributes = jsonencode({
            sourceName    = "sourceName"
            attributeName = "uid"
            sourceId      = "sourceId"
          })
        }
      },
      {
        identity_attribute_name = "email"
        transform_definition = {
          type = "accountAttribute"
          attributes = jsonencode({
            sourceName    = "sourceName"
            attributeName = "mailaddr"
            sourceId      = "sourceId"
          })
        }
      },
      {
        identity_attribute_name = "displayName"
        transform_definition = {
          type = "accountAttribute"
          attributes = jsonencode({
            sourceName    = "sourceName"
            attributeName = "displayName"
            sourceId      = "sourceId"
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
					resource.TestCheckResourceAttr("identitynow_identity_profile.test", "description", "Test Identity Profile UPDATED"),
					resource.TestCheckResourceAttr("identitynow_identity_profile.test", "priority", "10"),
					resource.TestCheckResourceAttr("identitynow_identity_profile.test", "owner.id", "ownerIdUpd"),
					resource.TestCheckResourceAttr("identitynow_identity_profile.test", "authoritative_source.id", "auth_sourceUpd"),
					resource.TestCheckResourceAttr("identitynow_identity_profile.test", "identity_attribute_config.enabled", "false"),
					resource.TestCheckResourceAttr("identitynow_identity_profile.test", "identity_attribute_config.attribute_transforms.0.identity_attribute_name", "uid"),
					resource.TestCheckResourceAttr("identitynow_identity_profile.test", "identity_attribute_config.attribute_transforms.0.transform_definition.type", "accountAttribute"),
					resource.TestCheckResourceAttr("identitynow_identity_profile.test", "identity_attribute_config.attribute_transforms.0.transform_definition.attributes", "{\"attributeName\":\"uid\",\"sourceId\":\"sourceId\",\"sourceName\":\"sourceName\"}"),
					resource.TestCheckResourceAttr("identitynow_identity_profile.test", "identity_attribute_config.attribute_transforms.1.identity_attribute_name", "email"),
					resource.TestCheckResourceAttr("identitynow_identity_profile.test", "identity_attribute_config.attribute_transforms.1.transform_definition.type", "accountAttribute"),
					resource.TestCheckResourceAttr("identitynow_identity_profile.test", "identity_attribute_config.attribute_transforms.1.transform_definition.attributes", "{\"attributeName\":\"mailaddr\",\"sourceId\":\"sourceId\",\"sourceName\":\"sourceName\"}"),
					resource.TestCheckResourceAttr("identitynow_identity_profile.test", "identity_attribute_config.attribute_transforms.2.identity_attribute_name", "displayName"),
					resource.TestCheckResourceAttr("identitynow_identity_profile.test", "identity_attribute_config.attribute_transforms.2.transform_definition.type", "accountAttribute"),
					resource.TestCheckResourceAttr("identitynow_identity_profile.test", "identity_attribute_config.attribute_transforms.2.transform_definition.attributes", "{\"attributeName\":\"displayName\",\"sourceId\":\"sourceId\",\"sourceName\":\"sourceName\"}"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
