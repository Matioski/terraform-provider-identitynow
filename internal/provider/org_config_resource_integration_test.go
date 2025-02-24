//go:build integration

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestIntegration_OrgConfig_CreateAndUpdate(t *testing.T) {
	t.Skip("Timezone changes were affecting the schedulerd jobs at the tenant. Skipping the test for now.")

	patchResource := providerIntegrationConfig + `
	resource "identitynow_org_config" "test" {
	  time_zone         = "Europe/Zurich"
	}
	`
	utcResource := providerIntegrationConfig + `
	resource "identitynow_org_config" "test" {
	  time_zone         = "` + defaultTimeZone + `"
	}
	`
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: patchResource,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("identitynow_org_config.test", "time_zone", "Europe/Zurich"),
				),
			},
			// Update and Read testing
			{
				Config: utcResource,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("identitynow_org_config.test", "time_zone", defaultTimeZone),
				),
			},
			// Delete does not happen in org_config
		},
	})
}
