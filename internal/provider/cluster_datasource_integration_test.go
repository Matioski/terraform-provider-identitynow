//go:build integration

package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"
)

func TestIntegration_ClusterDataSource(t *testing.T) {
	managedCluster := getManagedClusters(1, context.Background())[0]

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerIntegrationConfig + `data "identitynow_cluster" "test" { name = "` + *managedCluster.Name + `" }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.identitynow_cluster.test", "id", managedCluster.Id),
					resource.TestCheckResourceAttr("data.identitynow_cluster.test", "name", *managedCluster.Name),
					resource.TestCheckResourceAttr("data.identitynow_cluster.test", "pod", *managedCluster.Pod),
					resource.TestCheckResourceAttr("data.identitynow_cluster.test", "org", *managedCluster.Org),
				),
			},
		},
	})
}
