//go:build !integration

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestOrgConfigResource(t *testing.T) {
	t.Skipf("Unstable test with mockoon")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "identitynow_org_config" "test" {
  time_zone         = "Europe/Zurich"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("identitynow_org_config.test", "time_zone", "Europe/Zurich"),
				),
			},
			// Update and Read testing
			{
				Config: providerConfig + `
			resource "identitynow_org_config" "test" {
			  time_zone         = "UTC"
			}
			`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("identitynow_org_config.test", "time_zone", "UTC"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
