//go:build !integration

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestEntitlementDataSource_SearchByValue(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `
data "identitynow_entitlement" "test" {
    source_id = "1234567890"
    value = "ROLE_ADMIN"
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.identitynow_entitlement.test", "id", "31d5e3d5-5f18421bb51f74e847767657"),
					resource.TestCheckResourceAttr("data.identitynow_entitlement.test", "name", "Role Administrator"),
					resource.TestCheckNoResourceAttr("data.identitynow_entitlement.test", "attribute"),
					resource.TestCheckResourceAttr("data.identitynow_entitlement.test", "value", "ROLE_ADMIN"),
					resource.TestCheckResourceAttr("data.identitynow_entitlement.test", "source_id", "1234567890"),
				),
			},
		},
	})
}
