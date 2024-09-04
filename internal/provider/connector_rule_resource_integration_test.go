//go:build integration

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestIntegration_ConnectorRuleResource(t *testing.T) {
	ruleName := "integration-test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerIntegrationConfig + `
resource "identitynow_connector_rule" "test" {
  name        = "` + ruleName + `"
  description = "test rule description"
  type        = "BuildMap"
  signature = {
    input = [
      {
        name        = "firstName"
        description = "First Name of identity"
        type        = "String"
      },
      {
        name        = "lastName"
        description = "Last Name of identity"
        type        = "String"
      }
    ]
    output = {
      name        = "fullName"
      description = "Full Name of identity"
      type        = "String"
    }
  }
  source_code = {
    version = "1.0"
    script  = "return firstName + \" \" + lastName;"
  }
  attributes = jsonencode({
    attribute1 = "yes"
    attribute2 = "no"
  })
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("identitynow_connector_rule.test", "id"),
					resource.TestCheckResourceAttr("identitynow_connector_rule.test", "name", ruleName),
					resource.TestCheckResourceAttr("identitynow_connector_rule.test", "description", "test rule description"),
					resource.TestCheckResourceAttr("identitynow_connector_rule.test", "type", "BuildMap"),
					resource.TestCheckResourceAttr("identitynow_connector_rule.test", "signature.input.0.name", "firstName"),
					resource.TestCheckResourceAttr("identitynow_connector_rule.test", "signature.input.0.description", "First Name of identity"),
					resource.TestCheckResourceAttr("identitynow_connector_rule.test", "signature.input.0.type", "String"),
					resource.TestCheckResourceAttr("identitynow_connector_rule.test", "signature.input.1.name", "lastName"),
					resource.TestCheckResourceAttr("identitynow_connector_rule.test", "signature.input.1.description", "Last Name of identity"),
					resource.TestCheckResourceAttr("identitynow_connector_rule.test", "signature.input.1.type", "String"),
					resource.TestCheckResourceAttr("identitynow_connector_rule.test", "signature.output.name", "fullName"),
					resource.TestCheckResourceAttr("identitynow_connector_rule.test", "signature.output.description", "Full Name of identity"),
					resource.TestCheckResourceAttr("identitynow_connector_rule.test", "signature.output.type", "String"),
					resource.TestCheckResourceAttr("identitynow_connector_rule.test", "source_code.version", "1.0"),
					resource.TestCheckResourceAttr("identitynow_connector_rule.test", "source_code.script", "return firstName + \" \" + lastName;"),
					resource.TestCheckResourceAttr("identitynow_connector_rule.test", "attributes", "{\"attribute1\":\"yes\",\"attribute2\":\"no\"}"),
				),
			},
			// Update and Read testing
			{
				Config: providerIntegrationConfig + `
resource "identitynow_connector_rule" "test" {
  name        = "` + ruleName + `"
  description = "test rule descriptionUpd"
  type        = "BuildMap"
  signature = {
    input = [
      {
        name        = "firstNameUpd"
        description = "First Name of identityUpd"
        type        = "Boolean"
      },
      {
        name        = "lastNameUpd"
        description = "Last Name of identityUpd"
        type        = "Boolean"
      }
    ]
    output = {
      name        = "fullNameUpd"
      description = "Full Name of identityUpd"
      type        = "Boolean"
    }
  }
  source_code = {
    version = "1.1"
    script  = "return firstNameUpd + \" \" + lastNameUpd;"
  }
  attributes = jsonencode({
    attribute1 = "no"
    attribute2 = "yes"
  })
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("identitynow_connector_rule.test", "id"),
					resource.TestCheckResourceAttr("identitynow_connector_rule.test", "name", ruleName),
					resource.TestCheckResourceAttr("identitynow_connector_rule.test", "description", "test rule descriptionUpd"),
					resource.TestCheckResourceAttr("identitynow_connector_rule.test", "type", "BuildMap"),
					resource.TestCheckResourceAttr("identitynow_connector_rule.test", "signature.input.0.name", "firstNameUpd"),
					resource.TestCheckResourceAttr("identitynow_connector_rule.test", "signature.input.0.description", "First Name of identityUpd"),
					resource.TestCheckResourceAttr("identitynow_connector_rule.test", "signature.input.0.type", "Boolean"),
					resource.TestCheckResourceAttr("identitynow_connector_rule.test", "signature.input.1.name", "lastNameUpd"),
					resource.TestCheckResourceAttr("identitynow_connector_rule.test", "signature.input.1.description", "Last Name of identityUpd"),
					resource.TestCheckResourceAttr("identitynow_connector_rule.test", "signature.input.1.type", "Boolean"),
					resource.TestCheckResourceAttr("identitynow_connector_rule.test", "signature.output.name", "fullNameUpd"),
					resource.TestCheckResourceAttr("identitynow_connector_rule.test", "signature.output.description", "Full Name of identityUpd"),
					resource.TestCheckResourceAttr("identitynow_connector_rule.test", "signature.output.type", "Boolean"),
					resource.TestCheckResourceAttr("identitynow_connector_rule.test", "source_code.version", "1.1"),
					resource.TestCheckResourceAttr("identitynow_connector_rule.test", "source_code.script", "return firstNameUpd + \" \" + lastNameUpd;"),
					resource.TestCheckResourceAttr("identitynow_connector_rule.test", "attributes", "{\"attribute1\":\"no\",\"attribute2\":\"yes\"}"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
