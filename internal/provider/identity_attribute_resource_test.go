//go:build !integration

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestIdentityAttributeResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "identitynow_identity_attribute" "test" {
  name         = "test"
  display_name = "test"
  standard     = false
  type         = "string"
  multi = false
  searchable = false
  system = false
  sources = [
    {
      type = "rule"
      properties = jsonencode({
        ruleType = "IdentityAttribute"
        ruleName = "Cloud Promote Identity Attribute"
      })
    }
  ]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("identitynow_identity_attribute.test", "name", "test"),
					resource.TestCheckResourceAttr("identitynow_identity_attribute.test", "display_name", "test"),
					resource.TestCheckResourceAttr("identitynow_identity_attribute.test", "standard", "false"),
					resource.TestCheckResourceAttr("identitynow_identity_attribute.test", "type", "string"),
					resource.TestCheckResourceAttr("identitynow_identity_attribute.test", "multi", "false"),
					resource.TestCheckResourceAttr("identitynow_identity_attribute.test", "searchable", "false"),
					resource.TestCheckResourceAttr("identitynow_identity_attribute.test", "system", "false"),
					resource.TestCheckResourceAttr("identitynow_identity_attribute.test", "sources.0.type", "rule"),
					resource.TestCheckResourceAttr("identitynow_identity_attribute.test", "sources.0.properties", "{\"ruleName\":\"Cloud Promote Identity Attribute\",\"ruleType\":\"IdentityAttribute\"}"),
				),
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "identitynow_identity_attribute" "test" {
  name         = "test"
  display_name = "test update"
  standard     = true
  type         = "boolean"
  multi = true
  searchable = true
  system = true
  sources = [
    {
      type = "rule"
      properties = jsonencode({
        ruleType = "IdentityAttribute"
        ruleName = "Cloud Promote Identity Attribute Update"
        anotherAttribute = "anotherAttribute"
      })
    }
  ]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("identitynow_identity_attribute.test", "name", "test"),
					resource.TestCheckResourceAttr("identitynow_identity_attribute.test", "display_name", "test update"),
					resource.TestCheckResourceAttr("identitynow_identity_attribute.test", "standard", "true"),
					resource.TestCheckResourceAttr("identitynow_identity_attribute.test", "type", "boolean"),
					resource.TestCheckResourceAttr("identitynow_identity_attribute.test", "multi", "true"),
					resource.TestCheckResourceAttr("identitynow_identity_attribute.test", "searchable", "true"),
					resource.TestCheckResourceAttr("identitynow_identity_attribute.test", "system", "true"),
					resource.TestCheckResourceAttr("identitynow_identity_attribute.test", "sources.0.type", "rule"),
					resource.TestCheckResourceAttr("identitynow_identity_attribute.test", "sources.0.properties", "{\"anotherAttribute\":\"anotherAttribute\",\"ruleName\":\"Cloud Promote Identity Attribute Update\",\"ruleType\":\"IdentityAttribute\"}"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
