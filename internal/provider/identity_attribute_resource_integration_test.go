//go:build integration

package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func checkIdentityAttributeIsDeleted() func(state *terraform.State) error {
	return func(state *terraform.State) error {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "identitynow_identity_attribute" {
				continue
			}
			name := rs.Primary.ID
			_, response, _ := SPApiClient.Beta.IdentityAttributesAPI.GetIdentityAttribute(context.Background(), name).Execute()
			if response != nil && response.StatusCode == 404 {
				return nil
			}
			SPApiClient.Beta.IdentityAttributesAPI.DeleteIdentityAttribute(context.Background(), name).Execute()
			return fmt.Errorf("identity attribute still exists: %s", name)
		}
		return nil
	}
}

func TestIntegration_IdentityAttributeResource_Rule(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		CheckDestroy:             checkIdentityAttributeIsDeleted(),
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerIntegrationConfig + `
resource "identitynow_identity_attribute" "test" {
  name         = "TestIdentityAttributeResourceIntegration_Rule"
  display_name = "att_IdentityAttribute"
  standard     = false
  multi        = false
  type         = "string"
  searchable   = false
  system       = false
  sources      = [
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
					resource.TestCheckResourceAttr("identitynow_identity_attribute.test", "name", "TestIdentityAttributeResourceIntegration_Rule"),
					resource.TestCheckResourceAttr("identitynow_identity_attribute.test", "display_name", "att_IdentityAttribute"),
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
				Config: providerIntegrationConfig + `
resource "identitynow_identity_attribute" "test" {
  name         = "TestIdentityAttributeResourceIntegration_Rule"
  display_name = "att_IdentityAttribute_RuleUpd"
  standard     = true
  type         = "string"
  searchable   = true
  sources      = [
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
					resource.TestCheckResourceAttr("identitynow_identity_attribute.test", "name", "TestIdentityAttributeResourceIntegration_Rule"),
					resource.TestCheckResourceAttr("identitynow_identity_attribute.test", "display_name", "att_IdentityAttribute_RuleUpd"),
					resource.TestCheckResourceAttr("identitynow_identity_attribute.test", "standard", "true"),
					resource.TestCheckResourceAttr("identitynow_identity_attribute.test", "type", "string"),
					resource.TestCheckResourceAttr("identitynow_identity_attribute.test", "multi", "false"),
					resource.TestCheckResourceAttr("identitynow_identity_attribute.test", "searchable", "true"),
					resource.TestCheckResourceAttr("identitynow_identity_attribute.test", "system", "false"),
					resource.TestCheckResourceAttr("identitynow_identity_attribute.test", "sources.0.type", "rule"),
					resource.TestCheckResourceAttr("identitynow_identity_attribute.test", "sources.0.properties", "{\"ruleName\":\"Cloud Promote Identity Attribute\",\"ruleType\":\"IdentityAttribute\"}"),
				),
			},
			{
				Config: providerIntegrationConfig + `
resource "identitynow_identity_attribute" "test" {
  name         = "TestIdentityAttributeResourceIntegration_Rule"
  display_name = "att_IdentityAttribute_RuleUpd"
  standard     = false
  type         = "string"
  multi        = false
  searchable   = false
  system       = false
  sources      = [
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
					resource.TestCheckResourceAttr("identitynow_identity_attribute.test", "name", "TestIdentityAttributeResourceIntegration_Rule"),
					resource.TestCheckResourceAttr("identitynow_identity_attribute.test", "display_name", "att_IdentityAttribute_RuleUpd"),
					//standard has to be false for delete operation to work
					resource.TestCheckResourceAttr("identitynow_identity_attribute.test", "standard", "false"),
					resource.TestCheckResourceAttr("identitynow_identity_attribute.test", "type", "string"),
					resource.TestCheckResourceAttr("identitynow_identity_attribute.test", "multi", "false"),
					resource.TestCheckResourceAttr("identitynow_identity_attribute.test", "searchable", "false"),
					resource.TestCheckResourceAttr("identitynow_identity_attribute.test", "system", "false"),
					resource.TestCheckResourceAttr("identitynow_identity_attribute.test", "sources.0.type", "rule"),
					resource.TestCheckResourceAttr("identitynow_identity_attribute.test", "sources.0.properties", "{\"ruleName\":\"Cloud Promote Identity Attribute\",\"ruleType\":\"IdentityAttribute\"}"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
func TestIntegration_IdentityAttributeResource_ApplicationMapping(t *testing.T) {
	sourceCloud := getSources(2, context.Background())
	source1 := sourceCloud[0].ConnectorAttributes["cloudDisplayName"].(string)
	source2 := sourceCloud[1].ConnectorAttributes["cloudDisplayName"].(string)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		CheckDestroy:             checkIdentityAttributeIsDeleted(),
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerIntegrationConfig + `
resource "identitynow_identity_attribute" "test" {
  name         = "TestIdentityAttributeResourceIntegration_appMapping"
  display_name = "att_IdentityAttribute_appMapping"
  standard     = false
  multi        = false
  type         = "string"
  searchable   = false
  system       = false
  sources      = [
    {
      type = "applicationMapping"
      properties = jsonencode({
                attribute = "null"
                sourceName = "` + source1 + `"
      })
    },
    {
      type = "applicationMapping"
      properties = jsonencode({
                attribute = "null"
                sourceName = "` + source2 + `"
      })
    }
  ]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("identitynow_identity_attribute.test", "name", "TestIdentityAttributeResourceIntegration_appMapping"),
					resource.TestCheckResourceAttr("identitynow_identity_attribute.test", "display_name", "att_IdentityAttribute_appMapping"),
					resource.TestCheckResourceAttr("identitynow_identity_attribute.test", "standard", "false"),
					resource.TestCheckResourceAttr("identitynow_identity_attribute.test", "type", "string"),
					resource.TestCheckResourceAttr("identitynow_identity_attribute.test", "multi", "false"),
					resource.TestCheckResourceAttr("identitynow_identity_attribute.test", "searchable", "false"),
					resource.TestCheckResourceAttr("identitynow_identity_attribute.test", "system", "false"),
					resource.TestCheckResourceAttr("identitynow_identity_attribute.test", "sources.0.type", "applicationMapping"),
					resource.TestCheckResourceAttr("identitynow_identity_attribute.test", "sources.0.properties", "{\"attribute\":\"null\",\"sourceName\":\""+source1+"\"}"),
					resource.TestCheckResourceAttr("identitynow_identity_attribute.test", "sources.1.type", "applicationMapping"),
					resource.TestCheckResourceAttr("identitynow_identity_attribute.test", "sources.1.properties", "{\"attribute\":\"null\",\"sourceName\":\""+source2+"\"}"),
				),
			},
			// Update and Read testing
			{
				Config: providerIntegrationConfig + `
resource "identitynow_identity_attribute" "test" {
  name         = "TestIdentityAttributeResourceIntegration_appMapping"
  display_name = "att_IdentityAttributeUpd"
  standard     = true
  type         = "string"
  searchable   = true
  sources      = [
    {
      type = "applicationMapping"
      properties = jsonencode({
                attribute = "null"
                sourceName = "` + source2 + `"
      })
    }
  ]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("identitynow_identity_attribute.test", "name", "TestIdentityAttributeResourceIntegration_appMapping"),
					resource.TestCheckResourceAttr("identitynow_identity_attribute.test", "display_name", "att_IdentityAttributeUpd"),
					resource.TestCheckResourceAttr("identitynow_identity_attribute.test", "standard", "true"),
					resource.TestCheckResourceAttr("identitynow_identity_attribute.test", "type", "string"),
					resource.TestCheckResourceAttr("identitynow_identity_attribute.test", "multi", "false"),
					resource.TestCheckResourceAttr("identitynow_identity_attribute.test", "searchable", "true"),
					resource.TestCheckResourceAttr("identitynow_identity_attribute.test", "system", "false"),
					resource.TestCheckResourceAttr("identitynow_identity_attribute.test", "sources.0.type", "applicationMapping"),
					resource.TestCheckResourceAttr("identitynow_identity_attribute.test", "sources.0.properties", "{\"attribute\":\"null\",\"sourceName\":\""+source2+"\"}"),
				),
			},
			// Delete testing automatically occurs in TestCase
			{
				Config: providerIntegrationConfig + `
resource "identitynow_identity_attribute" "test" {
  name         = "TestIdentityAttributeResourceIntegration_appMapping"
  display_name = "att_IdentityAttribute_appMappingUpd"
  standard     = false
  type         = "string"
  multi        = false
  searchable   = false
  system       = false
  sources      = [
    {
      type = "applicationMapping"
      properties = jsonencode({
                attribute = "null"
                sourceName = "` + source1 + `"
      })
    }
  ]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("identitynow_identity_attribute.test", "name", "TestIdentityAttributeResourceIntegration_appMapping"),
					resource.TestCheckResourceAttr("identitynow_identity_attribute.test", "display_name", "att_IdentityAttribute_appMappingUpd"),
					//standard has to be false for delete operation to work
					resource.TestCheckResourceAttr("identitynow_identity_attribute.test", "standard", "false"),
					resource.TestCheckResourceAttr("identitynow_identity_attribute.test", "type", "string"),
					resource.TestCheckResourceAttr("identitynow_identity_attribute.test", "multi", "false"),
					resource.TestCheckResourceAttr("identitynow_identity_attribute.test", "searchable", "false"),
					resource.TestCheckResourceAttr("identitynow_identity_attribute.test", "system", "false"),
					resource.TestCheckResourceAttr("identitynow_identity_attribute.test", "sources.0.type", "applicationMapping"),
					resource.TestCheckResourceAttr("identitynow_identity_attribute.test", "sources.0.properties", "{\"attribute\":\"null\",\"sourceName\":\""+source1+"\"}"),
				),
			},
		},
	})
}
