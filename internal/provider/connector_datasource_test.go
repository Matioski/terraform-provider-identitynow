//go:build !integration

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestConnectorDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "identitynow_connector" "test" { name = "ADAM" }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.identitynow_connector.test", "name", "ADAM"),
					resource.TestCheckResourceAttr("data.identitynow_connector.test", "type", "ADAM - Direct"),
					resource.TestCheckResourceAttr("data.identitynow_connector.test", "script_name", "adam-angularsc"),
				),
			},
		},
	})
}
