//go:build integration

package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	sailpointBeta "github.com/sailpoint-oss/golang-sdk/v2/api_beta"
	"terraform-provider-identitynow/internal/util"
	"testing"
)

func TestIntegration_EntitlementDataSource_SearchByValue(t *testing.T) {
	entitlement := getRandomEntitlement(t)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerIntegrationConfig + `
data "identitynow_entitlement" "test" {
    source_id = "` + *entitlement.Source.Id + `"
    value = "` + *entitlement.Value + `"
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.identitynow_entitlement.test", "id", *entitlement.Id),
					resource.TestCheckResourceAttr("data.identitynow_entitlement.test", "name", *entitlement.Name),
					resource.TestCheckResourceAttr("data.identitynow_entitlement.test", "value", *entitlement.Value),
					resource.TestCheckResourceAttr("data.identitynow_entitlement.test", "source_id", *entitlement.Source.Id),
				),
			},
		},
	})
}

func TestIntegration_EntitlementDataSource_SearchByName(t *testing.T) {
	entitlement := getRandomEntitlement(t)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerIntegrationConfig + `
data "identitynow_entitlement" "test" {
    source_id = "` + *entitlement.Source.Id + `"
    name = "` + *entitlement.Name + `"
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.identitynow_entitlement.test", "id", *entitlement.Id),
					resource.TestCheckResourceAttr("data.identitynow_entitlement.test", "name", *entitlement.Name),
					resource.TestCheckResourceAttr("data.identitynow_entitlement.test", "value", *entitlement.Value),
					resource.TestCheckResourceAttr("data.identitynow_entitlement.test", "source_id", *entitlement.Source.Id),
				),
			},
		},
	})
}

func getRandomEntitlement(t *testing.T) *sailpointBeta.Entitlement {
	entitlements, spResp, err := SPApiClient.Beta.EntitlementsAPI.ListEntitlements(context.Background()).Limit(1).Execute()
	if err != nil {
		t.Fatalf("Error fetching entitlements: %s\n%s", err, util.GetBody(spResp))
		return nil
	}
	if len(entitlements) == 0 {
		t.Fatalf("No entitlements found")
		return nil
	}
	return &entitlements[0]
}
