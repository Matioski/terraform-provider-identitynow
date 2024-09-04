//go:build !integration

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestSourceSchemaResource_AddNew(t *testing.T) {
	t.Skipf("Unstable test with mockoon")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "identitynow_source_schema" "test" {
  name                = "account"
  source_id           = "source_id"
  native_object_type  = "User"
  identity_attribute  = "id"
  display_attribute   = "uid"
  hierarchy_attribute = "uid"
  include_permissions = false
  features = ["AUTHENTICATE"]
  configuration       = jsonencode({
    key = "value"
  })
  attributes = [
    {
      name           = "id"
      type           = "STRING"
      description    = "ID of the user"
      is_multi       = false
      is_entitlement = false
      is_group       = false
      schema         = {
        id = "yes"
        name = "name"
      }
    },
    {
      name           = "uid"
      type           = "STRING"
      description    = "UID of the user"
      is_multi       = true
      is_entitlement = true
      is_group       = true
      schema         = {
        id = "no"
        name = "name"
      }
    }
  ]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("identitynow_source_schema.test", "id"),
					resource.TestCheckResourceAttr("identitynow_source_schema.test", "name", "account"),
					resource.TestCheckResourceAttr("identitynow_source_schema.test", "source_id", "source_id"),
					resource.TestCheckResourceAttr("identitynow_source_schema.test", "native_object_type", "User"),
					resource.TestCheckResourceAttr("identitynow_source_schema.test", "identity_attribute", "id"),
					resource.TestCheckResourceAttr("identitynow_source_schema.test", "display_attribute", "uid"),
					resource.TestCheckResourceAttr("identitynow_source_schema.test", "include_permissions", "false"),
					resource.TestCheckResourceAttr("identitynow_source_schema.test", "features.0", "AUTHENTICATE"),
					resource.TestCheckResourceAttr("identitynow_source_schema.test", "configuration", "{\"key\":\"value\"}"),

					resource.TestCheckResourceAttr("identitynow_source_schema.test", "attributes.0.name", "id"),
					resource.TestCheckResourceAttr("identitynow_source_schema.test", "attributes.0.type", "STRING"),
					resource.TestCheckResourceAttr("identitynow_source_schema.test", "attributes.0.description", "ID of the user"),
					resource.TestCheckResourceAttr("identitynow_source_schema.test", "attributes.0.is_multi", "false"),
					resource.TestCheckResourceAttr("identitynow_source_schema.test", "attributes.0.is_entitlement", "false"),
					resource.TestCheckResourceAttr("identitynow_source_schema.test", "attributes.0.is_group", "false"),
					resource.TestCheckResourceAttr("identitynow_source_schema.test", "attributes.0.schema.type", "CONNECTOR_SCHEMA"),
					resource.TestCheckResourceAttr("identitynow_source_schema.test", "attributes.0.schema.id", "yes"),

					resource.TestCheckResourceAttr("identitynow_source_schema.test", "attributes.1.name", "uid"),
					resource.TestCheckResourceAttr("identitynow_source_schema.test", "attributes.1.type", "STRING"),
					resource.TestCheckResourceAttr("identitynow_source_schema.test", "attributes.1.description", "UID of the user"),
					resource.TestCheckResourceAttr("identitynow_source_schema.test", "attributes.1.is_multi", "true"),
					resource.TestCheckResourceAttr("identitynow_source_schema.test", "attributes.1.is_entitlement", "true"),
					resource.TestCheckResourceAttr("identitynow_source_schema.test", "attributes.1.is_group", "true"),
					resource.TestCheckResourceAttr("identitynow_source_schema.test", "attributes.1.schema.type", "CONNECTOR_SCHEMA"),
					resource.TestCheckResourceAttr("identitynow_source_schema.test", "attributes.1.schema.id", "no"),
				),
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "identitynow_source_schema" "test" {
  name                = "accountUpd"
  source_id           = "source_id"
  native_object_type  = "UserUpd"
  identity_attribute  = "another_idUpd"
  display_attribute   = "nameUpd"
  hierarchy_attribute = "uid"
  include_permissions = true
  features = ["AUTHENTICATE", "ENABLE"]
  configuration       = jsonencode({
    key = "valueUpd"
  })
  attributes = [
    {
      name           = "idUpd"
      type           = "BOOLEAN"
      description    = "ID of the user Upd"
      is_multi       = false
      is_entitlement = false
      is_group       = false
      schema         = {
        id = "yes"
        name = "nameUpd"
      }
    },
    {
      name           = "uidUpd"
      type           = "STRING"
      description    = "UID of the user Upd"
      is_multi       = true
      is_entitlement = true
      is_group       = true
      schema         = {
        id = "no"
        name = "nameUpd"
      }
    }
  ]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("identitynow_source_schema.test", "id"),
					resource.TestCheckResourceAttr("identitynow_source_schema.test", "name", "accountUpd"),
					resource.TestCheckResourceAttr("identitynow_source_schema.test", "source_id", "source_id"),
					resource.TestCheckResourceAttr("identitynow_source_schema.test", "native_object_type", "UserUpd"),
					resource.TestCheckResourceAttr("identitynow_source_schema.test", "identity_attribute", "another_idUpd"),
					resource.TestCheckResourceAttr("identitynow_source_schema.test", "display_attribute", "nameUpd"),
					resource.TestCheckResourceAttr("identitynow_source_schema.test", "include_permissions", "true"),
					resource.TestCheckResourceAttr("identitynow_source_schema.test", "features.0", "AUTHENTICATE"),
					resource.TestCheckResourceAttr("identitynow_source_schema.test", "features.1", "ENABLE"),
					resource.TestCheckResourceAttr("identitynow_source_schema.test", "configuration", "{\"key\":\"valueUpd\"}"),

					resource.TestCheckResourceAttr("identitynow_source_schema.test", "attributes.0.name", "idUpd"),
					resource.TestCheckResourceAttr("identitynow_source_schema.test", "attributes.0.type", "BOOLEAN"),
					resource.TestCheckResourceAttr("identitynow_source_schema.test", "attributes.0.description", "ID of the user Upd"),
					resource.TestCheckResourceAttr("identitynow_source_schema.test", "attributes.0.is_multi", "false"),
					resource.TestCheckResourceAttr("identitynow_source_schema.test", "attributes.0.is_entitlement", "false"),
					resource.TestCheckResourceAttr("identitynow_source_schema.test", "attributes.0.is_group", "false"),
					resource.TestCheckResourceAttr("identitynow_source_schema.test", "attributes.0.schema.type", "CONNECTOR_SCHEMA"),
					resource.TestCheckResourceAttr("identitynow_source_schema.test", "attributes.0.schema.id", "yes"),
					resource.TestCheckResourceAttr("identitynow_source_schema.test", "attributes.0.schema.name", "nameUpd"),

					resource.TestCheckResourceAttr("identitynow_source_schema.test", "attributes.1.name", "uidUpd"),
					resource.TestCheckResourceAttr("identitynow_source_schema.test", "attributes.1.type", "STRING"),
					resource.TestCheckResourceAttr("identitynow_source_schema.test", "attributes.1.description", "UID of the user Upd"),
					resource.TestCheckResourceAttr("identitynow_source_schema.test", "attributes.1.is_multi", "true"),
					resource.TestCheckResourceAttr("identitynow_source_schema.test", "attributes.1.is_entitlement", "true"),
					resource.TestCheckResourceAttr("identitynow_source_schema.test", "attributes.1.is_group", "true"),
					resource.TestCheckResourceAttr("identitynow_source_schema.test", "attributes.1.schema.type", "CONNECTOR_SCHEMA"),
					resource.TestCheckResourceAttr("identitynow_source_schema.test", "attributes.1.schema.id", "no"),
					resource.TestCheckResourceAttr("identitynow_source_schema.test", "attributes.1.schema.name", "nameUpd"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
