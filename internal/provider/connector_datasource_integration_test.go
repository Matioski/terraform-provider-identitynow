//go:build integration

package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"terraform-provider-identitynow/internal/util"
	"testing"
)

func TestIntegration_ConnectorDataSource(t *testing.T) {
	connector, spResp, err := SPApiClient.Beta.ConnectorsAPI.GetConnectorList(context.Background()).Limit(1).Execute()
	if err != nil {
		t.Fatalf("Error fetching sources: %s\n%s", err, util.GetBody(spResp))
	}
	if len(connector) == 0 {
		t.Fatalf("No sources found")
	}
	connectorName := *connector[0].Name

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerIntegrationConfig + `data "identitynow_connector" "test" { name = "` + connectorName + `" }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.identitynow_connector.test", "name", connectorName),
					resource.TestCheckResourceAttr("data.identitynow_connector.test", "type", *connector[0].Type),
					resource.TestCheckResourceAttr("data.identitynow_connector.test", "script_name", *connector[0].ScriptName),
				),
			},
		},
	})
}
