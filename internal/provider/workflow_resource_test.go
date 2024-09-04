//go:build !integration

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestWorfklowResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "identitynow_workflow" "test" {
  name        = "test workflow"
  description = "test workflow from terraform"
  owner = {
    id = "ownerId"
  }
  enabled = false
  definition = {
    start = "Get List of Identities"
    steps = jsonencode({
      "Get List of Identities" = {
        displayName   = "Get List of Identities"
        type          = "action"
        actionId      = "sp:get-identities"
        versionNumber = "2"
        attributes = {
          inputQuery = "*"
          searchBy   = "searchQuery"
        }
        nextStep = "End Step"
      }
      "End Step" = {
        displayName = ""
        type        = "success"
      }
    })
  }
  trigger = {
    type = "EVENT"
    attributes = {
      id = "idn:identity-created"
      filter = "filter"

      name = "name"
      description = "description"

      cron_string = "cron string"
    }
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("identitynow_workflow.test", "id"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "name", "test workflow"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "description", "test workflow from terraform"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "owner.type", "IDENTITY"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "owner.id", "ownerId"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "enabled", "false"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "definition.start", "Get List of Identities"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "definition.steps", "{\"End Step\":{\"displayName\":\"\",\"type\":\"success\"},\"Get List of Identities\":{\"actionId\":\"sp:get-identities\",\"attributes\":{\"inputQuery\":\"*\",\"searchBy\":\"searchQuery\"},\"displayName\":\"Get List of Identities\",\"nextStep\":\"End Step\",\"type\":\"action\",\"versionNumber\":\"2\"}}"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "trigger.type", "EVENT"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "trigger.attributes.id", "idn:identity-created"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "trigger.attributes.filter", "filter"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "trigger.attributes.name", "name"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "trigger.attributes.description", "description"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "trigger.attributes.cron_string", "cron string"),
				),
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "identitynow_workflow" "test" {
  name        = "test workflowUpd"
  description = "test workflow from terraformUpd"
  owner = {
    id = "ownerId"
  }
  enabled = true
  definition = {
    start = "Get List of IdentitiesUpd"
    steps = jsonencode({
      "Get List of Identities" = {
        displayName   = "Get List of IdentitiesUpd"
        type          = "action"
        actionId      = "sp:get-identities"
        versionNumber = "2"
        attributes = {
          inputQuery = "*"
          searchBy   = "searchQuery"
        }
        nextStep = "End Step"
      }
      "End Step" = {
        displayName = ""
        type        = "success"
      }
    })
  }
  trigger = {
    type = "EXTERNAL"
    attributes = {
      id = "idn:identity-createdUpd"
      filter = "filterUpd"

      name = "nameUpd"
      description = "descriptionUpd"

      cron_string = "cron stringUpd"
    }
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("identitynow_workflow.test", "id"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "name", "test workflowUpd"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "description", "test workflow from terraformUpd"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "owner.type", "IDENTITY"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "owner.id", "ownerId"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "enabled", "true"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "definition.start", "Get List of IdentitiesUpd"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "definition.steps", "{\"End Step\":{\"displayName\":\"\",\"type\":\"success\"},\"Get List of Identities\":{\"actionId\":\"sp:get-identities\",\"attributes\":{\"inputQuery\":\"*\",\"searchBy\":\"searchQuery\"},\"displayName\":\"Get List of IdentitiesUpd\",\"nextStep\":\"End Step\",\"type\":\"action\",\"versionNumber\":\"2\"}}"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "trigger.type", "EXTERNAL"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "trigger.attributes.id", "idn:identity-createdUpd"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "trigger.attributes.filter", "filterUpd"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "trigger.attributes.name", "nameUpd"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "trigger.attributes.description", "descriptionUpd"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "trigger.attributes.cron_string", "cron stringUpd"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
