//go:build integration

package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"strconv"
	"terraform-provider-identitynow/internal/util"
	"testing"
	"time"
)

func checkSourceIsDeleted() func(state *terraform.State) error {
	return func(state *terraform.State) error {
		// Verify the source is destroyed
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "identitynow_source" {
				continue
			}
			id := rs.Primary.ID
			_, response, _ := SPApiClient.V3.SourcesAPI.GetSource(context.Background(), id).Execute()
			if response != nil && response.StatusCode == 404 {
				return nil
			}
			SPApiClient.V3.SourcesAPI.DeleteSource(context.Background(), id).Execute()
			return fmt.Errorf("source still exists: %s", id)
		}
		return nil
	}
}

func TestIntegration_SourceResourcePatch_ConnectorChangeRecreatesSource(t *testing.T) {
	checkForPendingCisTask(context.Background())

	managedClusters, spResp, err := SPApiClient.Beta.ManagedClustersAPI.GetManagedClusters(context.Background()).Limit(1).Execute()

	if err != nil {
		t.Fatalf("Error fetching managed clusters: %s\n%s", err, util.GetBody(spResp))
	}
	if len(managedClusters) == 0 {
		t.Fatalf("No managed clusters found")
	}
	managedCluster := managedClusters[0]

	unixTimeStamp := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		CheckDestroy:             checkSourceIsDeleted(),
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerIntegrationConfig + `
resource "identitynow_source" "test" {
  name        = "TestIntegrationSourcePatch_ConnectorChangeRecreatesSource` + unixTimeStamp + `"
  description = "TestIntegrationSourcePatch_ConnectorChangeRecreatesSource"
  owner = {
    id   = "` + ownerIdentityId + `"
    name = "` + ownerIdentityName + `"
  }
  cluster = {
    id   = "` + managedCluster.Id + `"
    name = "` + *managedCluster.Name + `"
  }
  connector_attributes = jsonencode({
    enableLCS           = "true"
  })
  connector_attributes_credentials = jsonencode({
    username = "test"
    password = "test"
  })
  connector        = "delimited-file-angularsc"
  connector_class  = "sailpoint.connector.DelimitedFileConnector"
  connection_type  = "file"
  delete_threshold = 10
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("identitynow_source.test", "id"),
					resource.TestCheckResourceAttrSet("identitynow_source.test", "cloud_external_id"),
					resource.TestCheckResourceAttr("identitynow_source.test", "name", "TestIntegrationSourcePatch_ConnectorChangeRecreatesSource"+unixTimeStamp),
					resource.TestCheckResourceAttr("identitynow_source.test", "description", "TestIntegrationSourcePatch_ConnectorChangeRecreatesSource"),
					resource.TestCheckResourceAttr("identitynow_source.test", "owner.type", "IDENTITY"),
					resource.TestCheckResourceAttr("identitynow_source.test", "owner.id", ownerIdentityId),
					resource.TestCheckResourceAttr("identitynow_source.test", "owner.name", ownerIdentityName),
					resource.TestCheckResourceAttr("identitynow_source.test", "cluster.id", managedCluster.Id),
					resource.TestCheckResourceAttr("identitynow_source.test", "cluster.name", *managedCluster.Name),
					resource.TestCheckResourceAttr("identitynow_source.test", "connector", "delimited-file-angularsc"),
					resource.TestCheckResourceAttr("identitynow_source.test", "connector_class", "sailpoint.connector.DelimitedFileConnector"),
					resource.TestCheckResourceAttr("identitynow_source.test", "connector_attributes_credentials", "{\"password\":\"test\",\"username\":\"test\"}"),
					resource.TestCheckResourceAttr("identitynow_source.test", "delete_threshold", "10"),
					resource.TestCheckResourceAttr("identitynow_source.test", "connection_type", "file"),
				),
			},
			// Update and Read testing
			{
				Config: providerIntegrationConfig + `
resource "identitynow_source" "test" {
  name        = "Recreated TestIntegrationSourcePatch_ConnectorChangeRecreatesSource` + unixTimeStamp + `"
  description = "TestIntegrationSourcePatch_ConnectorChangeRecreatesSource"
  owner = {
    id   = "` + updatedOwnerIdentityId + `"
    name = "` + updatedOwnerIdentityName + `"
  }
  cluster = {
    id   = "` + managedCluster.Id + `"
    name = "` + *managedCluster.Name + `"
  }
  connector_attributes = jsonencode({
    user         = "user1234"
    driverClass  = "oracle.jdbc.driver.OracleDriver"
    testConnSQL  = "select * from dual"
    url          = "jdbc:mysql://localhost:3306/mysql "
  })
  connector_attributes_credentials = jsonencode({
    password = "duishadiusahdiuhasiudh"
  })
  connector_files = ["ojdbc10-19.18.0.0.jar"]
  connector = "jdbc-angularsc"
  connector_class = "sailpoint.connector.JDBCConnector"
  delete_threshold = 10
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("identitynow_source.test", "id"),
					resource.TestCheckResourceAttrSet("identitynow_source.test", "cloud_external_id"),
					resource.TestCheckResourceAttr("identitynow_source.test", "name", "Recreated TestIntegrationSourcePatch_ConnectorChangeRecreatesSource"+unixTimeStamp),
					resource.TestCheckResourceAttr("identitynow_source.test", "description", "TestIntegrationSourcePatch_ConnectorChangeRecreatesSource"),
					resource.TestCheckResourceAttr("identitynow_source.test", "owner.type", "IDENTITY"),
					resource.TestCheckResourceAttr("identitynow_source.test", "owner.id", updatedOwnerIdentityId),
					resource.TestCheckResourceAttr("identitynow_source.test", "owner.name", updatedOwnerIdentityName),
					resource.TestCheckResourceAttr("identitynow_source.test", "cluster.id", managedCluster.Id),
					resource.TestCheckResourceAttr("identitynow_source.test", "cluster.name", *managedCluster.Name),
					resource.TestCheckResourceAttr("identitynow_source.test", "connector", "jdbc-angularsc"),
					resource.TestCheckResourceAttr("identitynow_source.test", "connector_class", "sailpoint.connector.JDBCConnector"),
					resource.TestCheckResourceAttr("identitynow_source.test", "connector_attributes_credentials", "{\"password\":\"duishadiusahdiuhasiudh\"}"),
					resource.TestCheckResourceAttr("identitynow_source.test", "delete_threshold", "10"),
					resource.TestCheckResourceAttr("identitynow_source.test", "connection_type", "direct"),
				),
			},
		},
	})
}

func TestIntegration_SourceResourcePatch_UpdateConnectorFiles(t *testing.T) {
	checkForPendingCisTask(context.Background())
	managedCluster := getManagedClusters(1, context.Background())[0]
	unixTimeStamp := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		CheckDestroy:             checkSourceIsDeleted(),
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerIntegrationConfig + `
resource "identitynow_source" "test" {
  name        = "TestIntegrationSourcePatch_UpdateConnectorFiles` + unixTimeStamp + `"
  description = "TestIntegrationSourcePatch_ConnectorChangeRecreatesSource"
  owner = {
    id   = "` + ownerIdentityId + `"
    name = "` + ownerIdentityName + `"
  }
  cluster = {
    id   = "` + managedCluster.Id + `"
    name = "` + *managedCluster.Name + `"
  }
  connector_attributes = jsonencode({
    user         = "user1234"
    driverClass  = "oracle.jdbc.driver.OracleDriver"
    testConnSQL  = "select * from dual"
    url          = "jdbc:mysql://localhost:3306/mysql "
  })
  connector_attributes_credentials = jsonencode({
    password = "duishadiusahdiuhasiudh"
  })
  connector_files = ["ojdbc10-19.18.0.0.jar"]
  connector = "jdbc-angularsc"
  connector_class = "sailpoint.connector.JDBCConnector"
  delete_threshold = 10
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("identitynow_source.test", "id"),
					resource.TestCheckResourceAttrSet("identitynow_source.test", "cloud_external_id"),
					resource.TestCheckResourceAttr("identitynow_source.test", "name", "TestIntegrationSourcePatch_UpdateConnectorFiles"+unixTimeStamp),
					resource.TestCheckResourceAttr("identitynow_source.test", "description", "TestIntegrationSourcePatch_ConnectorChangeRecreatesSource"),
					resource.TestCheckResourceAttr("identitynow_source.test", "owner.type", "IDENTITY"),
					resource.TestCheckResourceAttr("identitynow_source.test", "owner.id", ownerIdentityId),
					resource.TestCheckResourceAttr("identitynow_source.test", "owner.name", ownerIdentityName),
					resource.TestCheckResourceAttr("identitynow_source.test", "cluster.id", managedCluster.Id),
					resource.TestCheckResourceAttr("identitynow_source.test", "cluster.name", *managedCluster.Name),
					resource.TestCheckResourceAttr("identitynow_source.test", "connector", "jdbc-angularsc"),
					resource.TestCheckResourceAttr("identitynow_source.test", "connector_class", "sailpoint.connector.JDBCConnector"),
					resource.TestCheckResourceAttr("identitynow_source.test", "connector_attributes_credentials", "{\"password\":\"duishadiusahdiuhasiudh\"}"),
					resource.TestCheckResourceAttr("identitynow_source.test", "delete_threshold", "10"),
					resource.TestCheckResourceAttr("identitynow_source.test", "connection_type", "direct"),
					resource.TestCheckResourceAttr("identitynow_source.test", "connector_files.0", "ojdbc10-19.18.0.0.jar"),
				),
			},
			// Update and Read testing
			{
				Config: providerIntegrationConfig + `
resource "identitynow_source" "test" {
  name        = "TestIntegrationSourcePatch_UpdateConnectorFiles` + unixTimeStamp + `"
  description = "TestIntegrationSourcePatch_ConnectorChangeRecreatesSource"
  owner = {
    id   = "` + ownerIdentityId + `"
    name = "` + ownerIdentityName + `"
  }
  cluster = {
    id   = "` + managedCluster.Id + `"
    name = "` + *managedCluster.Name + `"
  }
  connector_attributes = jsonencode({
    user         = "user1234"
    driverClass  = "oracle.jdbc.driver.OracleDriver"
    testConnSQL  = "select * from dual"
    url          = "jdbc:mysql://localhost:3306/mysql "
  })
  connector_attributes_credentials = jsonencode({
    password = "duishadiusahdiuhasiudh"
  })
  connector_files = ["ojdbc10-19.18.0.0.jar","ojdbc10-19.21.0.0.jar"]
  connector = "jdbc-angularsc"
  connector_class = "sailpoint.connector.JDBCConnector"
  delete_threshold = 10
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("identitynow_source.test", "id"),
					resource.TestCheckResourceAttrSet("identitynow_source.test", "cloud_external_id"),
					resource.TestCheckResourceAttr("identitynow_source.test", "name", "TestIntegrationSourcePatch_UpdateConnectorFiles"+unixTimeStamp),
					resource.TestCheckResourceAttr("identitynow_source.test", "description", "TestIntegrationSourcePatch_ConnectorChangeRecreatesSource"),
					resource.TestCheckResourceAttr("identitynow_source.test", "owner.type", "IDENTITY"),
					resource.TestCheckResourceAttr("identitynow_source.test", "owner.id", ownerIdentityId),
					resource.TestCheckResourceAttr("identitynow_source.test", "owner.name", ownerIdentityName),
					resource.TestCheckResourceAttr("identitynow_source.test", "cluster.id", managedCluster.Id),
					resource.TestCheckResourceAttr("identitynow_source.test", "cluster.name", *managedCluster.Name),
					resource.TestCheckResourceAttr("identitynow_source.test", "connector", "jdbc-angularsc"),
					resource.TestCheckResourceAttr("identitynow_source.test", "connector_class", "sailpoint.connector.JDBCConnector"),
					resource.TestCheckResourceAttr("identitynow_source.test", "connector_attributes_credentials", "{\"password\":\"duishadiusahdiuhasiudh\"}"),
					resource.TestCheckResourceAttr("identitynow_source.test", "delete_threshold", "10"),
					resource.TestCheckResourceAttr("identitynow_source.test", "connection_type", "direct"),
					resource.TestCheckResourceAttr("identitynow_source.test", "connector_files.0", "ojdbc10-19.18.0.0.jar"),
					resource.TestCheckResourceAttr("identitynow_source.test", "connector_files.1", "ojdbc10-19.21.0.0.jar"),
				),
			},
		},
	})
}

func TestIntegration_SourceResourcePatch_AddMgmtGroup(t *testing.T) {
	checkForPendingCisTask(context.Background())
	managedCluster := getManagedClusters(1, context.Background())[0]
	unixTimeStamp := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		CheckDestroy:             checkSourceIsDeleted(),
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerIntegrationConfig + `
resource "identitynow_source" "test" {
  name        = "TestIntegrationSourcePatch_AddMgmtGroup` + unixTimeStamp + `"
  description = "TestIntegrationSourcePatch_AddMgmtGroup"
  owner = {
    id   = "` + ownerIdentityId + `"
    name = "` + ownerIdentityName + `"
  }
  cluster = {
    id   = "` + managedCluster.Id + `"
    name = "` + *managedCluster.Name + `"
  }
  connector_attributes = jsonencode({
    enableLCS           = "true"
  })
  connector_attributes_credentials = jsonencode({
    username = "test"
    password = "test"
  })
  connector        = "delimited-file-angularsc"
  connector_class  = "sailpoint.connector.DelimitedFileConnector"
  connection_type  = "file"
  delete_threshold = 10
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("identitynow_source.test", "id"),
					resource.TestCheckResourceAttrSet("identitynow_source.test", "cloud_external_id"),
					resource.TestCheckResourceAttr("identitynow_source.test", "name", "TestIntegrationSourcePatch_AddMgmtGroup"+unixTimeStamp),
					resource.TestCheckResourceAttr("identitynow_source.test", "description", "TestIntegrationSourcePatch_AddMgmtGroup"),
					resource.TestCheckResourceAttr("identitynow_source.test", "owner.type", "IDENTITY"),
					resource.TestCheckResourceAttr("identitynow_source.test", "owner.id", ownerIdentityId),
					resource.TestCheckResourceAttr("identitynow_source.test", "owner.name", ownerIdentityName),
					resource.TestCheckResourceAttr("identitynow_source.test", "cluster.id", managedCluster.Id),
					resource.TestCheckResourceAttr("identitynow_source.test", "cluster.name", *managedCluster.Name),
					resource.TestCheckResourceAttr("identitynow_source.test", "connector", "delimited-file-angularsc"),
					resource.TestCheckResourceAttr("identitynow_source.test", "connector_class", "sailpoint.connector.DelimitedFileConnector"),
					resource.TestCheckResourceAttr("identitynow_source.test", "connector_attributes_credentials", "{\"password\":\"test\",\"username\":\"test\"}"),
					resource.TestCheckResourceAttr("identitynow_source.test", "delete_threshold", "10"),
					resource.TestCheckResourceAttr("identitynow_source.test", "connection_type", "file"),
				),
			},
			// Update and Read testing
			{
				Config: providerIntegrationConfig + `
resource "identitynow_source" "test" {
  name        = "TestIntegrationSourcePatch_AddMgmtGroup` + unixTimeStamp + `"
  description = "TestIntegrationSourcePatch_AddMgmtGroupUpd"
  owner = {
    id   = "` + updatedOwnerIdentityId + `"
    name = "` + updatedOwnerIdentityName + `"
  }
  cluster = {
    id   = "` + managedCluster.Id + `"
    name = "` + *managedCluster.Name + `"
  }
  connector_attributes = jsonencode({
    enableLCS           = "true"
  })
  connector_attributes_credentials = jsonencode({
    username = "testUpd"
    password = "testUpd"
  })
  connector        = "delimited-file-angularsc"
  connector_class  = "sailpoint.connector.DelimitedFileConnector"
  connection_type  = "file"
  delete_threshold = 15
  management_workgroup = {
    id = "786e45ee-3196-41d1-a7c1-d35aa0ebb4b0"
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("identitynow_source.test", "id"),
					resource.TestCheckResourceAttrSet("identitynow_source.test", "cloud_external_id"),
					resource.TestCheckResourceAttr("identitynow_source.test", "name", "TestIntegrationSourcePatch_AddMgmtGroup"+unixTimeStamp),
					resource.TestCheckResourceAttr("identitynow_source.test", "description", "TestIntegrationSourcePatch_AddMgmtGroupUpd"),
					resource.TestCheckResourceAttr("identitynow_source.test", "owner.type", "IDENTITY"),
					resource.TestCheckResourceAttr("identitynow_source.test", "owner.id", updatedOwnerIdentityId),
					resource.TestCheckResourceAttr("identitynow_source.test", "owner.name", updatedOwnerIdentityName),
					resource.TestCheckResourceAttr("identitynow_source.test", "cluster.id", managedCluster.Id),
					resource.TestCheckResourceAttr("identitynow_source.test", "cluster.name", *managedCluster.Name),
					resource.TestCheckResourceAttr("identitynow_source.test", "connector", "delimited-file-angularsc"),
					resource.TestCheckResourceAttr("identitynow_source.test", "connector_class", "sailpoint.connector.DelimitedFileConnector"),
					resource.TestCheckResourceAttr("identitynow_source.test", "connector_attributes_credentials", "{\"password\":\"testUpd\",\"username\":\"testUpd\"}"),
					resource.TestCheckResourceAttr("identitynow_source.test", "delete_threshold", "15"),
					resource.TestCheckResourceAttr("identitynow_source.test", "management_workgroup.id", "786e45ee-3196-41d1-a7c1-d35aa0ebb4b0"),
				),
			},
		},
	})
}

func TestIntegration_SourceResourcePatch_CreateWithAccountCorrelationConfig(t *testing.T) {
	checkForPendingCisTask(context.Background())
	managedCluster := getManagedClusters(1, context.Background())[0]
	unixTimeStamp := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		CheckDestroy:             checkSourceIsDeleted(),
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerIntegrationConfig + `
resource "identitynow_source" "test" {
  name        = "TestIntegrationSourcePatch_CreateWithAccountCorrelationConfig` + unixTimeStamp + `"
  description = "TestIntegrationSourcePatch_UpdatableValueChanged"
  owner = {
    id   = "` + ownerIdentityId + `"
    name = "` + ownerIdentityName + `"
  }
  cluster = {
    id   = "` + managedCluster.Id + `"
    name = "` + *managedCluster.Name + `"
  }
  connector_attributes = jsonencode({
    enableLCS           = "true"
  })
  connector_attributes_credentials = jsonencode({
    username = "test"
    password = "test"
  })
 account_correlation_config = {
        id = "ae8040c683fb4dfe8cea2770c8e7c741"
        type = "ACCOUNT_CORRELATION_CONFIG"
        name = "Active Directory Template Account Correlation Config"
    }
  connector        = "delimited-file-angularsc"
  connector_class  = "sailpoint.connector.DelimitedFileConnector"
  connection_type  = "file"
  delete_threshold = 10
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("identitynow_source.test", "id"),
					resource.TestCheckResourceAttrSet("identitynow_source.test", "cloud_external_id"),
					resource.TestCheckResourceAttr("identitynow_source.test", "name", "TestIntegrationSourcePatch_CreateWithAccountCorrelationConfig"+unixTimeStamp),
					resource.TestCheckResourceAttr("identitynow_source.test", "description", "TestIntegrationSourcePatch_UpdatableValueChanged"),
					resource.TestCheckResourceAttr("identitynow_source.test", "owner.type", "IDENTITY"),
					resource.TestCheckResourceAttr("identitynow_source.test", "owner.id", ownerIdentityId),
					resource.TestCheckResourceAttr("identitynow_source.test", "owner.name", ownerIdentityName),
					resource.TestCheckResourceAttr("identitynow_source.test", "cluster.id", managedCluster.Id),
					resource.TestCheckResourceAttr("identitynow_source.test", "cluster.name", *managedCluster.Name),
					resource.TestCheckResourceAttr("identitynow_source.test", "connector", "delimited-file-angularsc"),
					resource.TestCheckResourceAttr("identitynow_source.test", "connector_class", "sailpoint.connector.DelimitedFileConnector"),
					resource.TestCheckResourceAttr("identitynow_source.test", "connector_attributes_credentials", "{\"password\":\"test\",\"username\":\"test\"}"),
					resource.TestCheckResourceAttr("identitynow_source.test", "delete_threshold", "10"),
					resource.TestCheckResourceAttr("identitynow_source.test", "connection_type", "file"),
					resource.TestCheckResourceAttr("identitynow_source.test", "account_correlation_config.id", "ae8040c683fb4dfe8cea2770c8e7c741"),
					resource.TestCheckResourceAttr("identitynow_source.test", "account_correlation_config.type", "ACCOUNT_CORRELATION_CONFIG"),
					resource.TestCheckResourceAttr("identitynow_source.test", "account_correlation_config.name", "Active Directory Template Account Correlation Config"),
				),
			},
			// Update and Read testing
			{
				Config: providerIntegrationConfig + `
resource "identitynow_source" "test" {
  name        = "TestIntegrationSourcePatch_CreateWithAccountCorrelationConfigUpd` + unixTimeStamp + `"
  description = "TestIntegrationSourcePatch_UpdatableValueChangedUpd"
  owner = {
    id   = "` + ownerIdentityId + `"
    name = "` + ownerIdentityName + `"
  }
  cluster = {
    id   = "` + managedCluster.Id + `"
    name = "` + *managedCluster.Name + `"
  }
  connector_attributes = jsonencode({
    enableLCS           = "true"
  })
  connector_attributes_credentials = jsonencode({
    username = "testUpd"
    password = "testUpd"
  })
  connector        = "delimited-file-angularsc"
  connector_class  = "sailpoint.connector.DelimitedFileConnector"
  connection_type  = "file"
  delete_threshold = 15
 account_correlation_config = {
        id = "ae8040c683fb4dfe8cea2770c8e7c741"
        type = "ACCOUNT_CORRELATION_CONFIG"
        name = "Active Directory Template Account Correlation Config"
    }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("identitynow_source.test", "id"),
					resource.TestCheckResourceAttrSet("identitynow_source.test", "cloud_external_id"),
					resource.TestCheckResourceAttr("identitynow_source.test", "name", "TestIntegrationSourcePatch_CreateWithAccountCorrelationConfigUpd"+unixTimeStamp),
					resource.TestCheckResourceAttr("identitynow_source.test", "description", "TestIntegrationSourcePatch_UpdatableValueChangedUpd"),
					resource.TestCheckResourceAttr("identitynow_source.test", "owner.type", "IDENTITY"),
					resource.TestCheckResourceAttr("identitynow_source.test", "owner.id", ownerIdentityId),
					resource.TestCheckResourceAttr("identitynow_source.test", "owner.name", ownerIdentityName),
					resource.TestCheckResourceAttr("identitynow_source.test", "cluster.id", managedCluster.Id),
					resource.TestCheckResourceAttr("identitynow_source.test", "cluster.name", *managedCluster.Name),
					resource.TestCheckResourceAttr("identitynow_source.test", "connector", "delimited-file-angularsc"),
					resource.TestCheckResourceAttr("identitynow_source.test", "connector_class", "sailpoint.connector.DelimitedFileConnector"),
					resource.TestCheckResourceAttr("identitynow_source.test", "connector_attributes_credentials", "{\"password\":\"testUpd\",\"username\":\"testUpd\"}"),
					resource.TestCheckResourceAttr("identitynow_source.test", "delete_threshold", "15"),
					resource.TestCheckResourceAttr("identitynow_source.test", "account_correlation_config.id", "ae8040c683fb4dfe8cea2770c8e7c741"),
					resource.TestCheckResourceAttr("identitynow_source.test", "account_correlation_config.type", "ACCOUNT_CORRELATION_CONFIG"),
					resource.TestCheckResourceAttr("identitynow_source.test", "account_correlation_config.name", "Active Directory Template Account Correlation Config"),
				),
			},
		},
	})
}

func TestIntegration_SourceResourcePatch_CreateWithAllFields(t *testing.T) {
	checkForPendingCisTask(context.Background())
	managedCluster := getManagedClusters(1, context.Background())[0]
	unixTimeStamp := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		CheckDestroy:             checkSourceIsDeleted(),
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerIntegrationConfig + `
resource "identitynow_source" "test" {
  name        = "TestIntegrationSourcePatch_CreateWithAllFields` + unixTimeStamp + `"
  description = "TestIntegrationSourcePatch_UpdatableValueChanged"
  owner = {
    id   = "` + ownerIdentityId + `"
    name = "` + ownerIdentityName + `"
  }
  cluster = {
    id   = "` + managedCluster.Id + `"
    name = "` + *managedCluster.Name + `"
  }
  features = [
      "DISCOVER_SCHEMA",
      "DIRECT_PERMISSIONS",
      "NO_RANDOM_ACCESS"
  ]
  connector_attributes = jsonencode({
    enableLCS           = "true"
    innerObject = {
        innerKey = "innerValue"
    }
    arrayOfObjects = [{
        key1 = "value1"
    }]
    array = ["value1", "value2"]
  })
  connector_attributes_credentials = jsonencode({
    username = "test"
    password = "test"
    innerObject = {
        password = "test"
    }
    arrayOfObjects = [{
        password = "test"
    }]
    array = ["value3"]
  })
  management_workgroup = {
    id = "786e45ee-3196-41d1-a7c1-d35aa0ebb4b0"
  }
  account_correlation_config = {
        id = "ae8040c683fb4dfe8cea2770c8e7c741"
        type = "ACCOUNT_CORRELATION_CONFIG"
        name = "Active Directory Template Account Correlation Config"
  }
  manager_correlation_mapping = {
       account_attribute_name = "name"
       identity_attribute_name = "manager"
  }
  manager_correlation_rule = {
      id = "c291ffafb40045399e7fd4f44b8b3c40"
      type = "RULE"
      name = "Cloud Correlate Manager by AccountId"
  }
  connector        = "delimited-file-angularsc"
  connector_class  = "sailpoint.connector.DelimitedFileConnector"
  connection_type  = "file"
  delete_threshold = 10
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("identitynow_source.test", "id"),
					resource.TestCheckResourceAttrSet("identitynow_source.test", "cloud_external_id"),
					resource.TestCheckResourceAttr("identitynow_source.test", "name", "TestIntegrationSourcePatch_CreateWithAllFields"+unixTimeStamp),
					resource.TestCheckResourceAttr("identitynow_source.test", "description", "TestIntegrationSourcePatch_UpdatableValueChanged"),
					resource.TestCheckResourceAttr("identitynow_source.test", "owner.type", "IDENTITY"),
					resource.TestCheckResourceAttr("identitynow_source.test", "owner.id", ownerIdentityId),
					resource.TestCheckResourceAttr("identitynow_source.test", "owner.name", ownerIdentityName),
					resource.TestCheckResourceAttr("identitynow_source.test", "cluster.id", managedCluster.Id),
					resource.TestCheckResourceAttr("identitynow_source.test", "cluster.name", *managedCluster.Name),
					resource.TestCheckResourceAttr("identitynow_source.test", "connector", "delimited-file-angularsc"),
					resource.TestCheckResourceAttr("identitynow_source.test", "connector_class", "sailpoint.connector.DelimitedFileConnector"),
					resource.TestCheckResourceAttr("identitynow_source.test", "connector_attributes_credentials", "{\"array\":[\"value3\"],\"arrayOfObjects\":[{\"password\":\"test\"}],\"innerObject\":{\"password\":\"test\"},\"password\":\"test\",\"username\":\"test\"}"),
					resource.TestCheckResourceAttr("identitynow_source.test", "connector_attributes", "{\"array\":[\"value1\",\"value2\"],\"arrayOfObjects\":[{\"key1\":\"value1\"}],\"enableLCS\":\"true\",\"innerObject\":{\"innerKey\":\"innerValue\"}}"),
					resource.TestCheckResourceAttr("identitynow_source.test", "delete_threshold", "10"),
					resource.TestCheckResourceAttr("identitynow_source.test", "connection_type", "file"),
					resource.TestCheckResourceAttr("identitynow_source.test", "manager_correlation_mapping.account_attribute_name", "name"),
					resource.TestCheckResourceAttr("identitynow_source.test", "manager_correlation_mapping.identity_attribute_name", "manager"),
					resource.TestCheckResourceAttr("identitynow_source.test", "manager_correlation_rule.id", "c291ffafb40045399e7fd4f44b8b3c40"),
					resource.TestCheckResourceAttr("identitynow_source.test", "manager_correlation_rule.name", "Cloud Correlate Manager by AccountId"),
					resource.TestCheckResourceAttr("identitynow_source.test", "manager_correlation_rule.type", "RULE"),
					resource.TestCheckResourceAttr("identitynow_source.test", "management_workgroup.id", "786e45ee-3196-41d1-a7c1-d35aa0ebb4b0"),
					resource.TestCheckResourceAttr("identitynow_source.test", "account_correlation_config.id", "ae8040c683fb4dfe8cea2770c8e7c741"),
					resource.TestCheckResourceAttr("identitynow_source.test", "account_correlation_config.type", "ACCOUNT_CORRELATION_CONFIG"),
					resource.TestCheckResourceAttr("identitynow_source.test", "account_correlation_config.name", "Active Directory Template Account Correlation Config"),
					resource.TestCheckResourceAttr("identitynow_source.test", "features.#", "3"),
					resource.TestCheckResourceAttr("identitynow_source.test", "features.0", "DIRECT_PERMISSIONS"),
					resource.TestCheckResourceAttr("identitynow_source.test", "features.1", "DISCOVER_SCHEMA"),
					resource.TestCheckResourceAttr("identitynow_source.test", "features.2", "NO_RANDOM_ACCESS"),
				),
			},
			// Update and Read testing
			{
				Config: providerIntegrationConfig + `
resource "identitynow_source" "test" {
  name        = "TestIntegrationSourcePatch_CreateWithAllFieldsUpd` + unixTimeStamp + `"
  description = "TestIntegrationSourcePatch_UpdatableValueChangedUpd"
  owner = {
    id   = "` + updatedOwnerIdentityId + `"
    name = "` + updatedOwnerIdentityName + `"
  }
  cluster = {
    id   = "` + managedCluster.Id + `"
    name = "` + *managedCluster.Name + `"
  }
  features = [
      "DISCOVER_SCHEMA"
  ]
  connector_attributes = jsonencode({
    enableLCS           = "true"
    innerObject = {
        innerKey = "innerValueUpd"
    }
    arrayOfObjects = [{
        key1 = "value1Upd"
    }]
    array = ["value1Upd", "value2Upd"]
  })
  connector_attributes_credentials = jsonencode({
    username = "testUpd"
    password = "testUpd"
    innerObject = {
        password = "testUpd"
    }
    arrayOfObjects = [{
        password = "testUpd"
    }]
    array = ["value3Upd"]
  })
  management_workgroup = {
    id = "786e45ee-3196-41d1-a7c1-d35aa0ebb4b0"
  }
  account_correlation_config = {
        type = "ACCOUNT_CORRELATION_CONFIG"
        id   = "ae8040c683fb4dfe8cea2770c8e7c741"
        name = "Active Directory Template Account Correlation Config"
  }
  manager_correlation_mapping = {
      account_attribute_name = "name"
      identity_attribute_name = "manager"
  }
  manager_correlation_rule = {
      id = "c291ffafb40045399e7fd4f44b8b3c40"
      type = "RULE"
      name = "Cloud Correlate Manager by AccountId"
  }
  connector        = "delimited-file-angularsc"
  connector_class  = "sailpoint.connector.DelimitedFileConnector"
  connection_type  = "file"
  delete_threshold = 15
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("identitynow_source.test", "id"),
					resource.TestCheckResourceAttrSet("identitynow_source.test", "cloud_external_id"),
					resource.TestCheckResourceAttr("identitynow_source.test", "name", "TestIntegrationSourcePatch_CreateWithAllFieldsUpd"+unixTimeStamp),
					resource.TestCheckResourceAttr("identitynow_source.test", "description", "TestIntegrationSourcePatch_UpdatableValueChangedUpd"),
					resource.TestCheckResourceAttr("identitynow_source.test", "owner.type", "IDENTITY"),
					resource.TestCheckResourceAttr("identitynow_source.test", "owner.id", updatedOwnerIdentityId),
					resource.TestCheckResourceAttr("identitynow_source.test", "owner.name", updatedOwnerIdentityName),
					resource.TestCheckResourceAttr("identitynow_source.test", "cluster.id", managedCluster.Id),
					resource.TestCheckResourceAttr("identitynow_source.test", "cluster.name", *managedCluster.Name),
					resource.TestCheckResourceAttr("identitynow_source.test", "connector", "delimited-file-angularsc"),
					resource.TestCheckResourceAttr("identitynow_source.test", "connector_class", "sailpoint.connector.DelimitedFileConnector"),
					resource.TestCheckResourceAttr("identitynow_source.test", "connector_attributes_credentials", "{\"array\":[\"value3Upd\"],\"arrayOfObjects\":[{\"password\":\"testUpd\"}],\"innerObject\":{\"password\":\"testUpd\"},\"password\":\"testUpd\",\"username\":\"testUpd\"}"),
					resource.TestCheckResourceAttr("identitynow_source.test", "connector_attributes", "{\"array\":[\"value1Upd\",\"value2Upd\"],\"arrayOfObjects\":[{\"key1\":\"value1Upd\"}],\"enableLCS\":\"true\",\"innerObject\":{\"innerKey\":\"innerValueUpd\"}}"),
					resource.TestCheckResourceAttr("identitynow_source.test", "delete_threshold", "15"),
					resource.TestCheckResourceAttr("identitynow_source.test", "manager_correlation_mapping.account_attribute_name", "name"),
					resource.TestCheckResourceAttr("identitynow_source.test", "manager_correlation_mapping.identity_attribute_name", "manager"),
					resource.TestCheckResourceAttr("identitynow_source.test", "manager_correlation_rule.id", "c291ffafb40045399e7fd4f44b8b3c40"),
					resource.TestCheckResourceAttr("identitynow_source.test", "manager_correlation_rule.name", "Cloud Correlate Manager by AccountId"),
					resource.TestCheckResourceAttr("identitynow_source.test", "manager_correlation_rule.type", "RULE"),
					resource.TestCheckResourceAttr("identitynow_source.test", "management_workgroup.id", "786e45ee-3196-41d1-a7c1-d35aa0ebb4b0"),
					resource.TestCheckResourceAttr("identitynow_source.test", "features.#", "1"),
					resource.TestCheckResourceAttr("identitynow_source.test", "features.0", "DISCOVER_SCHEMA"),
					resource.TestCheckResourceAttr("identitynow_source.test", "account_correlation_config.id", "ae8040c683fb4dfe8cea2770c8e7c741"),
					resource.TestCheckResourceAttr("identitynow_source.test", "account_correlation_config.type", "ACCOUNT_CORRELATION_CONFIG"),
					resource.TestCheckResourceAttr("identitynow_source.test", "account_correlation_config.name", "Active Directory Template Account Correlation Config"),
				),
			},
		},
	})
}

func TestIntegration_SourceResourcePatch_UpdatableValueChanged(t *testing.T) {
	checkForPendingCisTask(context.Background())
	managedCluster := getManagedClusters(1, context.Background())[0]
	unixTimeStamp := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		CheckDestroy:             checkSourceIsDeleted(),
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerIntegrationConfig + `
resource "identitynow_source" "test" {
  name        = "TestIntegrationSourcePatch_UpdatableValueChanged` + unixTimeStamp + `"
  description = "TestIntegrationSourcePatch_UpdatableValueChanged"
  owner = {
    id   = "` + ownerIdentityId + `"
    name = "` + ownerIdentityName + `"
  }
  cluster = {
    id   = "` + managedCluster.Id + `"
    name = "` + *managedCluster.Name + `"
  }
  connector_attributes = jsonencode({
    enableLCS           = "true"
  })
  connector_attributes_credentials = jsonencode({
    username = "test"
    password = "test"
  })
  connector        = "delimited-file-angularsc"
  connector_class  = "sailpoint.connector.DelimitedFileConnector"
  connection_type  = "file"
  delete_threshold = 10
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("identitynow_source.test", "id"),
					resource.TestCheckResourceAttrSet("identitynow_source.test", "cloud_external_id"),
					resource.TestCheckResourceAttr("identitynow_source.test", "name", "TestIntegrationSourcePatch_UpdatableValueChanged"+unixTimeStamp),
					resource.TestCheckResourceAttr("identitynow_source.test", "description", "TestIntegrationSourcePatch_UpdatableValueChanged"),
					resource.TestCheckResourceAttr("identitynow_source.test", "owner.type", "IDENTITY"),
					resource.TestCheckResourceAttr("identitynow_source.test", "owner.id", ownerIdentityId),
					resource.TestCheckResourceAttr("identitynow_source.test", "owner.name", ownerIdentityName),
					resource.TestCheckResourceAttr("identitynow_source.test", "cluster.id", managedCluster.Id),
					resource.TestCheckResourceAttr("identitynow_source.test", "cluster.name", *managedCluster.Name),
					resource.TestCheckResourceAttr("identitynow_source.test", "connector", "delimited-file-angularsc"),
					resource.TestCheckResourceAttr("identitynow_source.test", "connector_class", "sailpoint.connector.DelimitedFileConnector"),
					resource.TestCheckResourceAttr("identitynow_source.test", "connector_attributes_credentials", "{\"password\":\"test\",\"username\":\"test\"}"),
					resource.TestCheckResourceAttr("identitynow_source.test", "delete_threshold", "10"),
					resource.TestCheckResourceAttr("identitynow_source.test", "connection_type", "file"),
				),
			},
			// Update and Read testing
			{
				Config: providerIntegrationConfig + `
resource "identitynow_source" "test" {
  name        = "TestIntegrationSourcePatch_UpdatableValueChangedUpd` + unixTimeStamp + `"
  description = "TestIntegrationSourcePatch_UpdatableValueChangedUpd"
  owner = {
    id   = "` + updatedOwnerIdentityId + `"
    name = "` + updatedOwnerIdentityName + `"
  }
  cluster = {
    id   = "` + managedCluster.Id + `"
    name = "` + *managedCluster.Name + `"
  }
  connector_attributes = jsonencode({
    enableLCS           = "true"
  })
  connector_attributes_credentials = jsonencode({
    username = "testUpd"
    password = "testUpd"
  })
  connector        = "delimited-file-angularsc"
  connector_class  = "sailpoint.connector.DelimitedFileConnector"
  connection_type  = "file"
  delete_threshold = 15
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("identitynow_source.test", "id"),
					resource.TestCheckResourceAttrSet("identitynow_source.test", "cloud_external_id"),
					resource.TestCheckResourceAttr("identitynow_source.test", "name", "TestIntegrationSourcePatch_UpdatableValueChangedUpd"+unixTimeStamp),
					resource.TestCheckResourceAttr("identitynow_source.test", "description", "TestIntegrationSourcePatch_UpdatableValueChangedUpd"),
					resource.TestCheckResourceAttr("identitynow_source.test", "owner.type", "IDENTITY"),
					resource.TestCheckResourceAttr("identitynow_source.test", "owner.id", updatedOwnerIdentityId),
					resource.TestCheckResourceAttr("identitynow_source.test", "owner.name", updatedOwnerIdentityName),
					resource.TestCheckResourceAttr("identitynow_source.test", "cluster.id", managedCluster.Id),
					resource.TestCheckResourceAttr("identitynow_source.test", "cluster.name", *managedCluster.Name),
					resource.TestCheckResourceAttr("identitynow_source.test", "connector", "delimited-file-angularsc"),
					resource.TestCheckResourceAttr("identitynow_source.test", "connector_class", "sailpoint.connector.DelimitedFileConnector"),
					resource.TestCheckResourceAttr("identitynow_source.test", "connector_attributes_credentials", "{\"password\":\"testUpd\",\"username\":\"testUpd\"}"),
					resource.TestCheckResourceAttr("identitynow_source.test", "delete_threshold", "15"),
					resource.TestCheckResourceAttr("identitynow_source.test", "connection_type", "file"),
				),
			},
		},
	})
}

func TestIntegration_SourceResourcePatch_DefaultValues_DontBreak(t *testing.T) {
	checkForPendingCisTask(context.Background())
	managedCluster := getManagedClusters(1, context.Background())[0]
	unixTimeStamp := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		CheckDestroy:             checkSourceIsDeleted(),
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerIntegrationConfig + `
resource "identitynow_source" "test" {
  name        = "TestIntegration_SourceResourcePatch_DefaultValues_DontBreak` + unixTimeStamp + `"
  description = "TestIntegration_SourceResourcePatch_DefaultValues_DontBreak"
  owner = {
    id   = "` + ownerIdentityId + `"
    name = "` + ownerIdentityName + `"
  }
  cluster = {
    id   = "` + managedCluster.Id + `"
    name = "` + *managedCluster.Name + `"
  }
  connector_attributes = jsonencode({
    enableLCS           = "true"
  })
  connector        = "active-directory-angularsc"
  connection_type  = "direct"
  delete_threshold = 10
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("identitynow_source.test", "id"),
					resource.TestCheckResourceAttrSet("identitynow_source.test", "cloud_external_id"),
					resource.TestCheckResourceAttr("identitynow_source.test", "name", "TestIntegration_SourceResourcePatch_DefaultValues_DontBreak"+unixTimeStamp),
					resource.TestCheckResourceAttr("identitynow_source.test", "description", "TestIntegration_SourceResourcePatch_DefaultValues_DontBreak"),
					resource.TestCheckResourceAttr("identitynow_source.test", "owner.type", "IDENTITY"),
					resource.TestCheckResourceAttr("identitynow_source.test", "owner.id", ownerIdentityId),
					resource.TestCheckResourceAttr("identitynow_source.test", "owner.name", ownerIdentityName),
					resource.TestCheckResourceAttr("identitynow_source.test", "cluster.id", managedCluster.Id),
					resource.TestCheckResourceAttr("identitynow_source.test", "cluster.name", *managedCluster.Name),
					resource.TestCheckResourceAttr("identitynow_source.test", "connector", "active-directory-angularsc"),
					resource.TestCheckResourceAttr("identitynow_source.test", "delete_threshold", "10"),
					resource.TestCheckResourceAttr("identitynow_source.test", "connection_type", "direct"),
				),
			},
			// Update and Read testing
			{
				Config: providerIntegrationConfig + `
resource "identitynow_source" "test" {
  name        = "TestIntegration_SourceResourcePatch_DefaultValues_DontBreak` + unixTimeStamp + `"
  description = "TestIntegration_SourceResourcePatch_DefaultValues_DontBreak"
  owner = {
    id   = "` + ownerIdentityId + `"
    name = "` + ownerIdentityName + `"
  }
  cluster = {
    id   = "` + managedCluster.Id + `"
    name = "` + *managedCluster.Name + `"
  }
  connector_attributes = jsonencode({
    enableLCS           = "true"
  })
  connector        = "active-directory-angularsc"
  connection_type  = "direct"
  delete_threshold = 10
  account_correlation_config = {
        type = "ACCOUNT_CORRELATION_CONFIG"
        id   = "ae8040c683fb4dfe8cea2770c8e7c741"
        name = "Active Directory Template Account Correlation Config"
  }
  manager_correlation_rule = {
        type = "RULE"
        id = "c291ffafb40045399e7fd4f44b8b3c40"
        name = "Cloud Correlate Manager by AccountId"
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("identitynow_source.test", "id"),
					resource.TestCheckResourceAttrSet("identitynow_source.test", "cloud_external_id"),
					resource.TestCheckResourceAttr("identitynow_source.test", "name", "TestIntegration_SourceResourcePatch_DefaultValues_DontBreak"+unixTimeStamp),
					resource.TestCheckResourceAttr("identitynow_source.test", "description", "TestIntegration_SourceResourcePatch_DefaultValues_DontBreak"),
					resource.TestCheckResourceAttr("identitynow_source.test", "owner.type", "IDENTITY"),
					resource.TestCheckResourceAttr("identitynow_source.test", "owner.id", ownerIdentityId),
					resource.TestCheckResourceAttr("identitynow_source.test", "owner.name", ownerIdentityName),
					resource.TestCheckResourceAttr("identitynow_source.test", "cluster.id", managedCluster.Id),
					resource.TestCheckResourceAttr("identitynow_source.test", "cluster.name", *managedCluster.Name),
					resource.TestCheckResourceAttr("identitynow_source.test", "connector", "active-directory-angularsc"),
					resource.TestCheckResourceAttr("identitynow_source.test", "delete_threshold", "10"),
					resource.TestCheckResourceAttr("identitynow_source.test", "connection_type", "direct"),
					resource.TestCheckResourceAttr("identitynow_source.test", "account_correlation_config.id", "ae8040c683fb4dfe8cea2770c8e7c741"),
					resource.TestCheckResourceAttr("identitynow_source.test", "account_correlation_config.type", "ACCOUNT_CORRELATION_CONFIG"),
					resource.TestCheckResourceAttr("identitynow_source.test", "account_correlation_config.name", "Active Directory Template Account Correlation Config"),
					resource.TestCheckResourceAttr("identitynow_source.test", "manager_correlation_rule.id", "c291ffafb40045399e7fd4f44b8b3c40"),
					resource.TestCheckResourceAttr("identitynow_source.test", "manager_correlation_rule.type", "RULE"),
					resource.TestCheckResourceAttr("identitynow_source.test", "manager_correlation_rule.name", "Cloud Correlate Manager by AccountId"),
				),
			},
		},
	})
}

func TestIntegration_SourceResourcePatch_RemoveOptionalFields(t *testing.T) {
	checkForPendingCisTask(context.Background())
	managedCluster := getManagedClusters(1, context.Background())[0]
	unixTimeStamp := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		CheckDestroy:             checkSourceIsDeleted(),
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerIntegrationConfig + `
resource "identitynow_source" "test" {
  name        = "TestIntegrationSourcePatch_CreateWithAllFields` + unixTimeStamp + `"
  description = "TestIntegrationSourcePatch_UpdatableValueChanged"
  owner = {
    id   = "` + ownerIdentityId + `"
    name = "` + ownerIdentityName + `"
  }
  cluster = {
    id   = "` + managedCluster.Id + `"
    name = "` + *managedCluster.Name + `"
  }
  features = [
      "DISCOVER_SCHEMA",
      "DIRECT_PERMISSIONS",
      "NO_RANDOM_ACCESS"
  ]
  connector_attributes = jsonencode({
    enableLCS           = "true"
    innerObject = {
        innerKey = "innerValue"
    }
    arrayOfObjects = [{
        key1 = "value1"
    }]
    array = ["value1", "value2"]
  })
  connector_attributes_credentials = jsonencode({
    username = "test"
    password = "test"
    innerObject = {
        password = "test"
    }
    arrayOfObjects = [{
        password = "test"
    }]
    array = ["value3"]
  })
  management_workgroup = {
    id = "786e45ee-3196-41d1-a7c1-d35aa0ebb4b0"
  }
  account_correlation_config = {
        id = "ae8040c683fb4dfe8cea2770c8e7c741"
        type = "ACCOUNT_CORRELATION_CONFIG"
        name = "Active Directory Template Account Correlation Config"
  }
  manager_correlation_mapping = {
       account_attribute_name = "name"
       identity_attribute_name = "manager"
  }
  manager_correlation_rule = {
      id = "c291ffafb40045399e7fd4f44b8b3c40"
      type = "RULE"
      name = "Cloud Correlate Manager by AccountId"
  }
  connector        = "delimited-file-angularsc"
  connector_class  = "sailpoint.connector.DelimitedFileConnector"
  connection_type  = "file"
  delete_threshold = 10
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("identitynow_source.test", "id"),
					resource.TestCheckResourceAttrSet("identitynow_source.test", "cloud_external_id"),
					resource.TestCheckResourceAttr("identitynow_source.test", "name", "TestIntegrationSourcePatch_CreateWithAllFields"+unixTimeStamp),
					resource.TestCheckResourceAttr("identitynow_source.test", "description", "TestIntegrationSourcePatch_UpdatableValueChanged"),
					resource.TestCheckResourceAttr("identitynow_source.test", "owner.type", "IDENTITY"),
					resource.TestCheckResourceAttr("identitynow_source.test", "owner.id", ownerIdentityId),
					resource.TestCheckResourceAttr("identitynow_source.test", "owner.name", ownerIdentityName),
					resource.TestCheckResourceAttr("identitynow_source.test", "cluster.id", managedCluster.Id),
					resource.TestCheckResourceAttr("identitynow_source.test", "cluster.name", *managedCluster.Name),
					resource.TestCheckResourceAttr("identitynow_source.test", "connector", "delimited-file-angularsc"),
					resource.TestCheckResourceAttr("identitynow_source.test", "connector_class", "sailpoint.connector.DelimitedFileConnector"),
					resource.TestCheckResourceAttr("identitynow_source.test", "connector_attributes_credentials", "{\"array\":[\"value3\"],\"arrayOfObjects\":[{\"password\":\"test\"}],\"innerObject\":{\"password\":\"test\"},\"password\":\"test\",\"username\":\"test\"}"),
					resource.TestCheckResourceAttr("identitynow_source.test", "connector_attributes", "{\"array\":[\"value1\",\"value2\"],\"arrayOfObjects\":[{\"key1\":\"value1\"}],\"enableLCS\":\"true\",\"innerObject\":{\"innerKey\":\"innerValue\"}}"),
					resource.TestCheckResourceAttr("identitynow_source.test", "delete_threshold", "10"),
					resource.TestCheckResourceAttr("identitynow_source.test", "connection_type", "file"),
					resource.TestCheckResourceAttr("identitynow_source.test", "manager_correlation_mapping.account_attribute_name", "name"),
					resource.TestCheckResourceAttr("identitynow_source.test", "manager_correlation_mapping.identity_attribute_name", "manager"),
					resource.TestCheckResourceAttr("identitynow_source.test", "manager_correlation_rule.id", "c291ffafb40045399e7fd4f44b8b3c40"),
					resource.TestCheckResourceAttr("identitynow_source.test", "manager_correlation_rule.name", "Cloud Correlate Manager by AccountId"),
					resource.TestCheckResourceAttr("identitynow_source.test", "manager_correlation_rule.type", "RULE"),
					resource.TestCheckResourceAttr("identitynow_source.test", "management_workgroup.id", "786e45ee-3196-41d1-a7c1-d35aa0ebb4b0"),
					resource.TestCheckResourceAttr("identitynow_source.test", "account_correlation_config.id", "ae8040c683fb4dfe8cea2770c8e7c741"),
					resource.TestCheckResourceAttr("identitynow_source.test", "account_correlation_config.type", "ACCOUNT_CORRELATION_CONFIG"),
					resource.TestCheckResourceAttr("identitynow_source.test", "account_correlation_config.name", "Active Directory Template Account Correlation Config"),
					resource.TestCheckResourceAttr("identitynow_source.test", "features.#", "3"),
					resource.TestCheckResourceAttr("identitynow_source.test", "features.0", "DIRECT_PERMISSIONS"),
					resource.TestCheckResourceAttr("identitynow_source.test", "features.1", "DISCOVER_SCHEMA"),
					resource.TestCheckResourceAttr("identitynow_source.test", "features.2", "NO_RANDOM_ACCESS"),
				),
			},
			// Update and Read testing
			{
				Config: providerIntegrationConfig + `
resource "identitynow_source" "test" {
  name        = "TestIntegrationSourcePatch_CreateWithAllFieldsUpd` + unixTimeStamp + `"
  description = "TestIntegrationSourcePatch_UpdatableValueChangedUpd"
  owner = {
    id   = "` + updatedOwnerIdentityId + `"
    name = "` + updatedOwnerIdentityName + `"
  }
  cluster = {
    id   = "` + managedCluster.Id + `"
    name = "` + *managedCluster.Name + `"
  }
  features = [
      "DISCOVER_SCHEMA"
  ]
  connector_attributes = jsonencode({
    enableLCS           = "true"
    innerObject = {
        innerKey = "innerValueUpd"
    }
    arrayOfObjects = [{
        key1 = "value1Upd"
    }]
    array = ["value1Upd", "value2Upd"]
  })
  connector_attributes_credentials = jsonencode({
    username = "testUpd"
    password = "testUpd"
    innerObject = {
        password = "testUpd"
    }
    arrayOfObjects = [{
        password = "testUpd"
    }]
    array = ["value3Upd"]
  })
  management_workgroup = {
    id = "786e45ee-3196-41d1-a7c1-d35aa0ebb4b0"
  }
  connector        = "delimited-file-angularsc"
  connector_class  = "sailpoint.connector.DelimitedFileConnector"
  connection_type  = "file"
  delete_threshold = 15
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("identitynow_source.test", "id"),
					resource.TestCheckResourceAttrSet("identitynow_source.test", "cloud_external_id"),
					resource.TestCheckResourceAttr("identitynow_source.test", "name", "TestIntegrationSourcePatch_CreateWithAllFieldsUpd"+unixTimeStamp),
					resource.TestCheckResourceAttr("identitynow_source.test", "description", "TestIntegrationSourcePatch_UpdatableValueChangedUpd"),
					resource.TestCheckResourceAttr("identitynow_source.test", "owner.type", "IDENTITY"),
					resource.TestCheckResourceAttr("identitynow_source.test", "owner.id", updatedOwnerIdentityId),
					resource.TestCheckResourceAttr("identitynow_source.test", "owner.name", updatedOwnerIdentityName),
					resource.TestCheckResourceAttr("identitynow_source.test", "cluster.id", managedCluster.Id),
					resource.TestCheckResourceAttr("identitynow_source.test", "cluster.name", *managedCluster.Name),
					resource.TestCheckResourceAttr("identitynow_source.test", "connector", "delimited-file-angularsc"),
					resource.TestCheckResourceAttr("identitynow_source.test", "connector_class", "sailpoint.connector.DelimitedFileConnector"),
					resource.TestCheckResourceAttr("identitynow_source.test", "connector_attributes_credentials", "{\"array\":[\"value3Upd\"],\"arrayOfObjects\":[{\"password\":\"testUpd\"}],\"innerObject\":{\"password\":\"testUpd\"},\"password\":\"testUpd\",\"username\":\"testUpd\"}"),
					resource.TestCheckResourceAttr("identitynow_source.test", "connector_attributes", "{\"array\":[\"value1Upd\",\"value2Upd\"],\"arrayOfObjects\":[{\"key1\":\"value1Upd\"}],\"enableLCS\":\"true\",\"innerObject\":{\"innerKey\":\"innerValueUpd\"}}"),
					resource.TestCheckResourceAttr("identitynow_source.test", "delete_threshold", "15"),
					resource.TestCheckResourceAttr("identitynow_source.test", "features.#", "1"),
					resource.TestCheckResourceAttr("identitynow_source.test", "features.0", "DISCOVER_SCHEMA"),
					resource.TestCheckNoResourceAttr("identitynow_source.test", "account_correlation_config"),
					resource.TestCheckNoResourceAttr("identitynow_source.test", "manager_correlation_mapping"),
					resource.TestCheckNoResourceAttr("identitynow_source.test", "manager_correlation_rule"),
				),
			},
		},
	})
}
