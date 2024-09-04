//go:build integration

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// not possible to update name and type
func TestIntegration_TransformResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerIntegrationConfig + `
resource "identitynow_transform" "test" {
  name         = "test"
  type         = "substring"
  attributes = jsonencode({
    regex       = "^(PRP_)*"
    replacement = ""
  })
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("identitynow_transform.test", "id"),
					resource.TestCheckResourceAttr("identitynow_transform.test", "name", "test"),
					resource.TestCheckResourceAttr("identitynow_transform.test", "type", "substring"),
					resource.TestCheckResourceAttr("identitynow_transform.test", "attributes", "{\"regex\":\"^(PRP_)*\",\"replacement\":\"\"}"),
				),
			},
			// Update and Read testing
			{
				Config: providerIntegrationConfig + `
resource "identitynow_transform" "test" {
  name         = "test"
  type         = "substring"
  attributes = jsonencode({
    end = 5
    begin = {
      type = "indexOf"
    }
    beginOffset = 1
  })
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("identitynow_transform.test", "id"),
					resource.TestCheckResourceAttr("identitynow_transform.test", "name", "test"),
					resource.TestCheckResourceAttr("identitynow_transform.test", "type", "substring"),
					resource.TestCheckResourceAttr("identitynow_transform.test", "attributes", "{\"begin\":{\"type\":\"indexOf\"},\"beginOffset\":1,\"end\":5}"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
