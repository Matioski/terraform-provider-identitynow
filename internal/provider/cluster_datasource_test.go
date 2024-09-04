//go:build !integration

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestClusterDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "identitynow_cluster" "test" { name = "clusterName" }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.identitynow_cluster.test", "id", "clusterId"),
					resource.TestCheckResourceAttr("data.identitynow_cluster.test", "name", "clusterName"),
					resource.TestCheckResourceAttr("data.identitynow_cluster.test", "pod", "eu-pod"),
					resource.TestCheckResourceAttr("data.identitynow_cluster.test", "org", "test-org"),
				),
			},
		},
	})
}
