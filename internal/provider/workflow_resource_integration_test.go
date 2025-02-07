//go:build integration

package provider

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func checkWorkflowIsDeleted() func(state *terraform.State) error {
	return func(state *terraform.State) error {
		// Verify the source is destroyed
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "identitynow_workflow" {
				continue
			}
			id := rs.Primary.ID
			_, response, _ := SPApiClient.Beta.WorkflowsAPI.GetWorkflow(context.Background(), id).Execute()
			if response != nil && response.StatusCode == 404 {
				return nil
			}
			return fmt.Errorf("workflow still exists: %s", id)
		}
		return nil
	}
}

func TestIntegration_WorkflowResource_Name(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		CheckDestroy:             checkWorkflowIsDeleted(),
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerIntegrationConfig + `
resource "identitynow_workflow" "test" {
  name        = "testWorkflow_changeName"
  description = "test workflow from terraform"
  owner = {
    id   = "` + ownerIdentityId + `"
  }
  enabled = true
  definition = {
    start = "Get List of Identities"
    steps = jsonencode({
      "Get List of Identities" = {
        displayName   = "Get List of Identities"
        type          = "action"
        actionId      = "sp:get-identities"
        versionNumber = 2
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
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("identitynow_workflow.test", "id"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "name", "testWorkflow_changeName"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "description", "test workflow from terraform"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "owner.type", "IDENTITY"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "owner.id", ownerIdentityId),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "enabled", "true"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "definition.start", "Get List of Identities"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "definition.steps", "{\"End Step\":{\"displayName\":\"\",\"type\":\"success\"},\"Get List of Identities\":{\"actionId\":\"sp:get-identities\",\"attributes\":{\"inputQuery\":\"*\",\"searchBy\":\"searchQuery\"},\"displayName\":\"Get List of Identities\",\"nextStep\":\"End Step\",\"type\":\"action\",\"versionNumber\":2}}"),
				),
			},
			// Update and Read testing
			{
				Config: providerIntegrationConfig + `
resource "identitynow_workflow" "test" {
  name        = "testWorkflow_changeNameUpd"
  description = "test workflow from terraformUpd"
  owner = {
    id   = "` + ownerIdentityId + `"
  }
  enabled = false
  definition = {
    start = "Get List of Identities"
    steps = jsonencode({
      "Get List of Identities" = {
        displayName   = "Get List of Identities"
        type          = "action"
        actionId      = "sp:get-identities"
        versionNumber = 2
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
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("identitynow_workflow.test", "id"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "name", "testWorkflow_changeNameUpd"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "description", "test workflow from terraformUpd"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "owner.type", "IDENTITY"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "owner.id", ownerIdentityId),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "enabled", "false"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "definition.start", "Get List of Identities"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "definition.steps", "{\"End Step\":{\"displayName\":\"\",\"type\":\"success\"},\"Get List of Identities\":{\"actionId\":\"sp:get-identities\",\"attributes\":{\"inputQuery\":\"*\",\"searchBy\":\"searchQuery\"},\"displayName\":\"Get List of Identities\",\"nextStep\":\"End Step\",\"type\":\"action\",\"versionNumber\":2}}"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestIntegration_WorkflowResource_Definition(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		CheckDestroy:             checkWorkflowIsDeleted(),
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerIntegrationConfig + `
resource "identitynow_workflow" "test" {
  name        = "testWorkflow_changeDefinition"
  description = "test workflow from terraform"
  owner = {
    id   = "` + ownerIdentityId + `"
  }
  enabled = false
  definition = {
    start = "Get List of Identities"
    steps = jsonencode({
      "Get List of Identities" = {
        displayName   = "Get List of Identities"
        type          = "action"
        actionId      = "sp:get-identities"
        versionNumber = 2
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
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("identitynow_workflow.test", "id"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "name", "testWorkflow_changeDefinition"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "description", "test workflow from terraform"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "owner.type", "IDENTITY"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "owner.id", ownerIdentityId),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "enabled", "false"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "definition.start", "Get List of Identities"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "definition.steps", "{\"End Step\":{\"displayName\":\"\",\"type\":\"success\"},\"Get List of Identities\":{\"actionId\":\"sp:get-identities\",\"attributes\":{\"inputQuery\":\"*\",\"searchBy\":\"searchQuery\"},\"displayName\":\"Get List of Identities\",\"nextStep\":\"End Step\",\"type\":\"action\",\"versionNumber\":2}}"),
				),
			},
			// Update and Read testing
			{
				Config: providerIntegrationConfig + `
resource "identitynow_workflow" "test" {
  name        = "testWorkflow_changeDefinition"
  description = "test workflow from terraformUpd"
  owner = {
    id   = "` + ownerIdentityId + `"
  }
  enabled = false
  definition = {
    start = "Get List of Identities"
    steps = jsonencode({
      "Get List of Identities" = {
        displayName   = "Get List of IdentitiesUpd"
        type          = "action"
        actionId      = "sp:get-identities"
        versionNumber = 2
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
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("identitynow_workflow.test", "id"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "name", "testWorkflow_changeDefinition"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "description", "test workflow from terraformUpd"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "owner.type", "IDENTITY"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "owner.id", ownerIdentityId),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "enabled", "false"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "definition.start", "Get List of Identities"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "definition.steps", "{\"End Step\":{\"displayName\":\"\",\"type\":\"success\"},\"Get List of Identities\":{\"actionId\":\"sp:get-identities\",\"attributes\":{\"inputQuery\":\"*\",\"searchBy\":\"searchQuery\"},\"displayName\":\"Get List of IdentitiesUpd\",\"nextStep\":\"End Step\",\"type\":\"action\",\"versionNumber\":2}}"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestIntegration_WorkflowResource_Identity(t *testing.T) {
	t.Skipf("Disabled test due to ISC bug that prevents deletion of subscriptions when deleting their associated workflows, causing them to pile up, reaching the tenant limit, and eventually block the creation of new workflows")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		CheckDestroy:             checkWorkflowIsDeleted(),
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerIntegrationConfig + `
resource "identitynow_workflow" "test" {
  name        = "testWorkflow_changeOwner"
  description = "test workflow from terraform"
  owner = {
    id   = "` + ownerIdentityId + `"
  }
  enabled = false
  definition = {
    start = "Get List of Identities"
    steps = jsonencode({
      "Get List of Identities" = {
        displayName   = "Get List of Identities"
        type          = "action"
        actionId      = "sp:get-identities"
        versionNumber = 2
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
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("identitynow_workflow.test", "id"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "name", "testWorkflow_changeOwner"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "description", "test workflow from terraform"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "owner.type", "IDENTITY"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "owner.id", ownerIdentityId),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "enabled", "false"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "definition.start", "Get List of Identities"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "definition.steps", "{\"End Step\":{\"displayName\":\"\",\"type\":\"success\"},\"Get List of Identities\":{\"actionId\":\"sp:get-identities\",\"attributes\":{\"inputQuery\":\"*\",\"searchBy\":\"searchQuery\"},\"displayName\":\"Get List of Identities\",\"nextStep\":\"End Step\",\"type\":\"action\",\"versionNumber\":2}}"),
				),
			},
			// Update and Read testing
			{
				Config: providerIntegrationConfig + `
resource "identitynow_workflow" "test" {
  name        = "testWorkflow_changeOwner"
  description = "test workflow from terraform"
  owner = {
    id   = "` + updatedOwnerIdentityId + `"
  }
  enabled = false
  definition = {
    start = "Get List of Identities"
    steps = jsonencode({
      "Get List of Identities" = {
        displayName   = "Get List of Identities"
        type          = "action"
        actionId      = "sp:get-identities"
        versionNumber = 2
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
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("identitynow_workflow.test", "id"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "name", "testWorkflow_changeOwner"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "description", "test workflow from terraform"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "owner.type", "IDENTITY"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "owner.id", updatedOwnerIdentityId),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "enabled", "false"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "definition.start", "Get List of Identities"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "definition.steps", "{\"End Step\":{\"displayName\":\"\",\"type\":\"success\"},\"Get List of Identities\":{\"actionId\":\"sp:get-identities\",\"attributes\":{\"inputQuery\":\"*\",\"searchBy\":\"searchQuery\"},\"displayName\":\"Get List of Identities\",\"nextStep\":\"End Step\",\"type\":\"action\",\"versionNumber\":2}}"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
func TestIntegration_WorkflowResource_Trigger(t *testing.T) {
	t.Skipf("Disabled test due to ISC bug that prevents deletion of subscriptions when deleting their associated workflows, causing them to pile up, reaching the tenant limit, and eventually block the creation of new workflows")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		CheckDestroy:             checkWorkflowIsDeleted(),
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerIntegrationConfig + `
resource "identitynow_workflow" "test" {
  name        = "testWorkflow_changeTrigger"
  description = "test workflow from terraform"
  owner = {
    id   = "` + ownerIdentityId + `"
  }
  enabled = false
  definition = {
    start = "Get List of Identities"
    steps = jsonencode({
      "Get List of Identities" = {
        displayName   = "Get List of Identities"
        type          = "action"
        actionId      = "sp:get-identities"
        versionNumber = 2
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
    }
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("identitynow_workflow.test", "id"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "name", "testWorkflow_changeTrigger"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "description", "test workflow from terraform"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "owner.type", "IDENTITY"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "owner.id", ownerIdentityId),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "enabled", "false"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "definition.start", "Get List of Identities"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "definition.steps", "{\"End Step\":{\"displayName\":\"\",\"type\":\"success\"},\"Get List of Identities\":{\"actionId\":\"sp:get-identities\",\"attributes\":{\"inputQuery\":\"*\",\"searchBy\":\"searchQuery\"},\"displayName\":\"Get List of Identities\",\"nextStep\":\"End Step\",\"type\":\"action\",\"versionNumber\":2}}"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "trigger.type", "EVENT"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "trigger.attributes.id", "idn:identity-created"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "trigger.attributes.filter", "filter"),
				),
			},
			// Update and Read testing
			{
				Config: providerIntegrationConfig + `
resource "identitynow_workflow" "test" {
  name        = "testWorkflow_changeTrigger"
  description = "test workflow from terraform"
  owner = {
    id   = "` + ownerIdentityId + `"
  }
  enabled = true
  definition = {
    start = "Get List of Identities"
    steps = jsonencode({
      "Get List of Identities" = {
        displayName   = "Get List of IdentitiesUpd"
        type          = "action"
        actionId      = "sp:get-identities"
        versionNumber = 2
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
      name = "xx"
      description = "xx"
    }
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("identitynow_workflow.test", "id"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "name", "testWorkflow_changeTrigger"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "description", "test workflow from terraform"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "owner.type", "IDENTITY"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "owner.id", ownerIdentityId),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "enabled", "true"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "definition.start", "Get List of Identities"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "definition.steps", "{\"End Step\":{\"displayName\":\"\",\"type\":\"success\"},\"Get List of Identities\":{\"actionId\":\"sp:get-identities\",\"attributes\":{\"inputQuery\":\"*\",\"searchBy\":\"searchQuery\"},\"displayName\":\"Get List of IdentitiesUpd\",\"nextStep\":\"End Step\",\"type\":\"action\",\"versionNumber\":2}}"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "trigger.type", "EXTERNAL"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "trigger.attributes.name", "xx"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "trigger.attributes.description", "xx"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestIntegration_WorkflowResource_DisabledToEnabled(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		CheckDestroy:             checkWorkflowIsDeleted(),
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerIntegrationConfig + `
resource "identitynow_workflow" "test" {
  name        = "testWorkflow_EnabledToDisabled"
  description = "test workflow from terraform"
  owner = {
    id   = "` + ownerIdentityId + `"
  }
  enabled = false
   definition = {
    start = "Get List of Identities"
    steps = jsonencode({
      "Get List of Identities" = {
        displayName   = "Get List of Identities"
        type          = "action"
        actionId      = "sp:get-identities"
        versionNumber = 2
        attributes = {
          inputQuery = "*"
          searchBy   = "searchQuery"
        }
        nextStep = "End Step - Success"
      }
      "End Step - Success" = {
        displayName = ""
        type        = "success"
      }
    })
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("identitynow_workflow.test", "id"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "name", "testWorkflow_EnabledToDisabled"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "description", "test workflow from terraform"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "owner.type", "IDENTITY"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "owner.id", ownerIdentityId),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "enabled", "false"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "definition.start", "Get List of Identities"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "definition.steps", "{\"End Step - Success\":{\"displayName\":\"\",\"type\":\"success\"},\"Get List of Identities\":{\"actionId\":\"sp:get-identities\",\"attributes\":{\"inputQuery\":\"*\",\"searchBy\":\"searchQuery\"},\"displayName\":\"Get List of Identities\",\"nextStep\":\"End Step - Success\",\"type\":\"action\",\"versionNumber\":2}}"),
				),
			},
			// Update and Read testing
			{
				Config: providerIntegrationConfig + `
resource "identitynow_workflow" "test" {
  name        = "testWorkflow_EnabledToDisabled"
  description = "test workflow from terraform"
  owner = {
    id   = "` + ownerIdentityId + `"
  }
  enabled = true
   definition = {
    start = "Get List of Identities"
    steps = jsonencode({
      "Get List of Identities" = {
        displayName   = "Get List of Identities"
        type          = "action"
        actionId      = "sp:get-identities"
        versionNumber = 2
        attributes = {
          inputQuery = "*"
          searchBy   = "searchQuery"
        }
        nextStep = "End Step - Success"
      }
      "End Step - Success" = {
        displayName = ""
        type        = "success"
      }
    })
  }
}
`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("identitynow_workflow.test", "id"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "name", "testWorkflow_EnabledToDisabled"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "description", "test workflow from terraform"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "owner.type", "IDENTITY"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "owner.id", ownerIdentityId),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "enabled", "true"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "definition.start", "Get List of Identities"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "definition.steps", "{\"End Step - Success\":{\"displayName\":\"\",\"type\":\"success\"},\"Get List of Identities\":{\"actionId\":\"sp:get-identities\",\"attributes\":{\"inputQuery\":\"*\",\"searchBy\":\"searchQuery\"},\"displayName\":\"Get List of Identities\",\"nextStep\":\"End Step - Success\",\"type\":\"action\",\"versionNumber\":2}}"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestIntegration_WorkflowResource_EventTrigger(t *testing.T) {
	t.Skipf("Disabled to avoid pile up of trigger subscriptions.")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		CheckDestroy:             checkWorkflowIsDeleted(),
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerIntegrationConfig + `
resource "identitynow_workflow" "test" {
  name        = "TestIntegration_WorkflowResource_EventTrigger"
  description = "test workflow from terraform"
  owner = {
    id   = "` + ownerIdentityId + `"
  }
  enabled = false
  definition = {
    start = "Get List of Identities"
    steps = jsonencode({
      "Get List of Identities" = {
        displayName   = "Get List of Identities"
        type          = "action"
        actionId      = "sp:get-identities"
        versionNumber = 2
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
      id = "idn:identity-attributes-changed"
      filter = "$.changes[?(@.attribute== \"cloudLifecycleState\")]"
    }
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("identitynow_workflow.test", "id"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "name", "TestIntegration_WorkflowResource_EventTrigger"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "description", "test workflow from terraform"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "owner.type", "IDENTITY"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "owner.id", ownerIdentityId),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "enabled", "false"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "definition.start", "Get List of Identities"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "definition.steps", "{\"End Step\":{\"displayName\":\"\",\"type\":\"success\"},\"Get List of Identities\":{\"actionId\":\"sp:get-identities\",\"attributes\":{\"inputQuery\":\"*\",\"searchBy\":\"searchQuery\"},\"displayName\":\"Get List of Identities\",\"nextStep\":\"End Step\",\"type\":\"action\",\"versionNumber\":2}}"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "trigger.type", "EVENT"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "trigger.attributes.id", "idn:identity-attributes-changed"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "trigger.attributes.filter", "$.changes[?(@.attribute== \"cloudLifecycleState\")]"),
				),
			},
			// Update and Read testing
			{
				Config: providerIntegrationConfig + `
			resource "identitynow_workflow" "test" {
			  name        = "TestIntegration_WorkflowResource_EventTrigger"
			  description = "test workflow from terraform"
			  owner = {
			    id   = "` + ownerIdentityId + `"
			  }
			  enabled = false
			  definition = {
			    start = "Get List of Identities"
			    steps = jsonencode({
			      "Get List of Identities" = {
			        displayName   = "Get List of IdentitiesUpd"
			        type          = "action"
			        actionId      = "sp:get-identities"
			        versionNumber = 2
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
          name = "externalName"
          description = "externalDesc"
        }
      }
    }
			`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("identitynow_workflow.test", "id"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "name", "TestIntegration_WorkflowResource_EventTrigger"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "description", "test workflow from terraform"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "owner.type", "IDENTITY"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "owner.id", ownerIdentityId),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "enabled", "false"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "definition.start", "Get List of Identities"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "definition.steps", "{\"End Step\":{\"displayName\":\"\",\"type\":\"success\"},\"Get List of Identities\":{\"actionId\":\"sp:get-identities\",\"attributes\":{\"inputQuery\":\"*\",\"searchBy\":\"searchQuery\"},\"displayName\":\"Get List of IdentitiesUpd\",\"nextStep\":\"End Step\",\"type\":\"action\",\"versionNumber\":2}}"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "trigger.type", "EXTERNAL"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "trigger.attributes.name", "externalName"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "trigger.attributes.description", "externalDesc"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestIntegration_WorkflowResource_ModifyEnabledWf(t *testing.T) {
	unixTimeStamp := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		CheckDestroy:             checkWorkflowIsDeleted(),
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerIntegrationConfig + `
resource "identitynow_workflow" "test" {
  name        = "TestIntegration_WorkflowResource_ModifyEnabledWf_` + unixTimeStamp + `"
  description = "Wf to test disabling and enabled wf to patch it and enable it again afterwards"
  owner = {
    id   = "` + ownerIdentityId + `"
  }
  enabled = true
  definition = {
    start = "Get List of Identities"
    steps = jsonencode({
      "Get List of Identities" = {
        displayName   = "Get List of Identities"
        type          = "action"
        actionId      = "sp:get-identities"
        versionNumber = 2
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
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("identitynow_workflow.test", "id"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "name", "TestIntegration_WorkflowResource_ModifyEnabledWf_"+unixTimeStamp),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "description", "Wf to test disabling and enabled wf to patch it and enable it again afterwards"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "enabled", "true"),
				),
			},
			// Update maintaining enabled = true, and Read testing
			{
				Config: providerIntegrationConfig + `
resource "identitynow_workflow" "test" {
  enabled = true
  name        = "TestIntegration_WorkflowResource_ModifyEnabledWf_` + unixTimeStamp + `_patched"
  description = "Wf to test disabling and enabled wf to patch it and enable it again afterwards - patched"
  owner = {
    id   = "` + ownerIdentityId + `"
  }
  definition = {
    start = "Get List of Identities"
    steps = jsonencode({
      "Get List of Identities" = {
        displayName   = "Get List of Identities"
        type          = "action"
        actionId      = "sp:get-identities"
        versionNumber = 2
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
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("identitynow_workflow.test", "id"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "name", "TestIntegration_WorkflowResource_ModifyEnabledWf_"+unixTimeStamp+"_patched"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "description", "Wf to test disabling and enabled wf to patch it and enable it again afterwards - patched"),
					resource.TestCheckResourceAttr("identitynow_workflow.test", "enabled", "true"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
