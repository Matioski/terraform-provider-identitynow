//go:build !integration

package provider

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"
)

func TestIdentityDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "identitynow_identity" "test" { alias = "John.Doe" }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.identitynow_identity.test", "id", "id12345"),
					resource.TestCheckResourceAttr("data.identitynow_identity.test", "name", "John Doe"),
					resource.TestCheckResourceAttr("data.identitynow_identity.test", "email", "john.doe@example.com"),
				),
			},
		},
	})
}
