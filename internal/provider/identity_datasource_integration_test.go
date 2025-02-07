//go:build integration

package provider

import (
	"context"
	"terraform-provider-identitynow/internal/util"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestIntegration_IdentityDataSource(t *testing.T) {
	identity, spResp, err := SPApiClient.Beta.IdentitiesAPI.ListIdentities(context.Background()).Limit(1).Execute()
	if err != nil {
		t.Fatalf("Error fetching sources: %s\n%s", err, util.GetBody(spResp))
	}
	if len(identity) == 0 {
		t.Fatalf("No sources found")
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerIntegrationConfig + `data "identitynow_identity" "test" { alias = "` + *identity[0].Alias + `" }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.identitynow_identity.test", "id", *identity[0].Id),
					resource.TestCheckResourceAttr("data.identitynow_identity.test", "name", identity[0].Name),
					//resource.TestCheckResourceAttr("data.identitynow_identity.test", "email", *identity[0].EmailAddress.Get()),
				),
			},
		},
	})
}
