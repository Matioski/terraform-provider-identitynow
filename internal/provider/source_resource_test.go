//go:build !integration

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestSourceResource(t *testing.T) {
	t.Skipf("Unstable test with mockoon")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing (non DelimitedFile source)
			{
				Config: providerConfig + `
resource "identitynow_source" "test" {
  name         = "test source"
  description  = "Creating from integration tests"
  owner = {
    id   =  "ownerId"
  }
  cluster = {
    id   = "clusterId"
    name   = "clusterName"
  }
    account_correlation_config = {
        id = "accCorId"
    }
    account_correlation_rule = {
        id = "accCorRuleId"
    }
    manager_correlation_mapping = {
        account_attribute_name = "manAccAttr"
        identity_attribute_name = "manIdentAttr"
    }
    manager_correlation_rule = {
        id = "manCorRuleId"
    }
    before_provisioning_rule = {
        id = "befProvRuleId"
    }
  features = ["ENABLE"]
  connector_files = ["ojdbc10-19.18.0.0.txt"]
  connector_attributes = jsonencode({
    enableLCS           = true
    inherited = {
      first   = "1"
      seconds = "2"
    }
  })
  connector_attributes_credentials = jsonencode({
    inherited = {
      clientId     = "clientId"
      clientSecret = "clientSecret"
    }
  })
  connector        = "custom connector"
  delete_threshold = 10
    management_workgroup = {
        id = "manWorkId"
    }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("identitynow_source.test", "id"),
					resource.TestCheckResourceAttrSet("identitynow_source.test", "cloud_external_id"),
					resource.TestCheckResourceAttr("identitynow_source.test", "name", "test source"),
					resource.TestCheckResourceAttr("identitynow_source.test", "description", "Creating from integration tests"),
					resource.TestCheckResourceAttr("identitynow_source.test", "owner.type", "IDENTITY"),
					resource.TestCheckResourceAttr("identitynow_source.test", "owner.id", "ownerId"),
					resource.TestCheckResourceAttr("identitynow_source.test", "cluster.id", "clusterId"),
					resource.TestCheckResourceAttr("identitynow_source.test", "account_correlation_config.type", "ACCOUNT_CORRELATION_CONFIG"),
					resource.TestCheckResourceAttr("identitynow_source.test", "account_correlation_config.id", "accCorId"),
					resource.TestCheckResourceAttr("identitynow_source.test", "account_correlation_rule.type", "RULE"),
					resource.TestCheckResourceAttr("identitynow_source.test", "account_correlation_rule.id", "accCorRuleId"),
					resource.TestCheckResourceAttr("identitynow_source.test", "manager_correlation_mapping.account_attribute_name", "manAccAttr"),
					resource.TestCheckResourceAttr("identitynow_source.test", "manager_correlation_mapping.identity_attribute_name", "manIdentAttr"),
					resource.TestCheckResourceAttr("identitynow_source.test", "manager_correlation_rule.type", "RULE"),
					resource.TestCheckResourceAttr("identitynow_source.test", "manager_correlation_rule.id", "manCorRuleId"),
					resource.TestCheckResourceAttr("identitynow_source.test", "before_provisioning_rule.type", "RULE"),
					resource.TestCheckResourceAttr("identitynow_source.test", "before_provisioning_rule.id", "befProvRuleId"),
					resource.TestCheckResourceAttr("identitynow_source.test", "features.0", "ENABLE"),
					resource.TestCheckResourceAttr("identitynow_source.test", "type", "connectorType"),
					resource.TestCheckResourceAttr("identitynow_source.test", "connector", "custom connector"),
					resource.TestCheckResourceAttr("identitynow_source.test", "connector_attributes", "{\"enableLCS\":true,\"inherited\":{\"first\":\"1\",\"seconds\":\"2\"}}"),
					resource.TestCheckResourceAttr("identitynow_source.test", "connector_attributes_credentials", "{\"inherited\":{\"clientId\":\"clientId\",\"clientSecret\":\"clientSecret\"}}"),
					resource.TestCheckResourceAttr("identitynow_source.test", "delete_threshold", "10"),
					resource.TestCheckResourceAttr("identitynow_source.test", "management_workgroup.type", "GOVERNANCE_GROUP"),
					resource.TestCheckResourceAttr("identitynow_source.test", "management_workgroup.id", "manWorkId"),
					resource.TestCheckResourceAttr("identitynow_source.test", "connector_files.0", "ojdbc10-19.18.0.0.txt"),
				),
			},

			// Create and Read testing (DelimitedFile source)
			{
				Config: providerConfig + `
resource "identitynow_source" "test" {
  name         = "test csv source"
  description  = "Creating from integration tests"
  owner = {
    id   =  "ownerId"
  }
  cluster = {
    id   = "clusterId"
    name   = "clusterName"
  }
    account_correlation_config = {
        id = "accCorId"
    }
    account_correlation_rule = {
        id = "accCorRuleId"
    }
    manager_correlation_mapping = {
        account_attribute_name = "manAccAttr"
        identity_attribute_name = "manIdentAttr"
    }
    manager_correlation_rule = {
        id = "manCorRuleId"
    }
  features = ["ENABLE"]
  connection_type = "file"
  connector_attributes = jsonencode({
    enableLCS           = true
    inherited = {
      first   = "1"
      seconds = "2"
    }
  })
  connector        = "delimited-file-angularsc"
  delete_threshold = 10
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("identitynow_source.test", "id"),
					resource.TestCheckResourceAttrSet("identitynow_source.test", "cloud_external_id"),
					resource.TestCheckResourceAttr("identitynow_source.test", "name", "test csv source"),
					resource.TestCheckResourceAttr("identitynow_source.test", "description", "Creating from integration tests"),
					resource.TestCheckResourceAttr("identitynow_source.test", "owner.type", "IDENTITY"),
					resource.TestCheckResourceAttr("identitynow_source.test", "owner.id", "ownerId"),
					resource.TestCheckResourceAttr("identitynow_source.test", "cluster.id", "clusterId"),
					resource.TestCheckResourceAttr("identitynow_source.test", "account_correlation_config.type", "ACCOUNT_CORRELATION_CONFIG"),
					resource.TestCheckResourceAttr("identitynow_source.test", "account_correlation_config.id", "accCorId"),
					resource.TestCheckResourceAttr("identitynow_source.test", "account_correlation_rule.type", "RULE"),
					resource.TestCheckResourceAttr("identitynow_source.test", "account_correlation_rule.id", "accCorRuleId"),
					resource.TestCheckResourceAttr("identitynow_source.test", "manager_correlation_mapping.account_attribute_name", "manAccAttr"),
					resource.TestCheckResourceAttr("identitynow_source.test", "manager_correlation_mapping.identity_attribute_name", "manIdentAttr"),
					resource.TestCheckResourceAttr("identitynow_source.test", "manager_correlation_rule.type", "RULE"),
					resource.TestCheckResourceAttr("identitynow_source.test", "manager_correlation_rule.id", "manCorRuleId"),
					resource.TestCheckResourceAttr("identitynow_source.test", "features.0", "ENABLE"),
					resource.TestCheckResourceAttr("identitynow_source.test", "type", "DelimitedFile"),
					resource.TestCheckResourceAttr("identitynow_source.test", "connector", "custom connector"),
					resource.TestCheckResourceAttr("identitynow_source.test", "connector_attributes", "{\"enableLCS\":true,\"inherited\":{\"first\":\"1\",\"seconds\":\"2\"}}"),
					resource.TestCheckResourceAttr("identitynow_source.test", "delete_threshold", "10"),
				),
			},

			// Update and Read testing
			{
				Config: providerConfig + `
resource "identitynow_source" "test" {
  name         = "test source update"
  description  = "Creating from integration tests update"
  owner = {
    id   =  "ownerIdUpd"
  }
  cluster = {
    id   = "clusterIdUpd"
    name   = "clusterNameUpd"
  }
    account_correlation_config = {
        id = "accCorIdUpd"
    }
    account_correlation_rule = {
        id = "accCorRuleIdUpd"
    }
    manager_correlation_mapping = {
        account_attribute_name = "manAccAttrUpd"
        identity_attribute_name = "manIdentAttrUpd"
    }
    manager_correlation_rule = {
        id = "manCorRuleIdUpd"
    }
    before_provisioning_rule = {
        id = "befProvRuleIdUpd"
    }
  features = ["AUTHENTICATE", "ENABLE"]
connector_files = ["ojdbc10-19.18.0.0.txt", "ojdbc10-19.21.0.0.txt"]
  connector_attributes = jsonencode({
    enableLCS           = true
    inherited = {
      first   = "1Upd"
      seconds = "2"
    }
  })
  connector_attributes_credentials = jsonencode({
    inherited = {
      clientId     = "clientId"
      clientSecret = "clientSecretUpd"
    }
  })
  connector        = "custom connector"
  delete_threshold = 15
    management_workgroup = {
        id = "manWorkIdUpd"
    }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("identitynow_source.test", "id"),
					resource.TestCheckResourceAttrSet("identitynow_source.test", "cloud_external_id"),
					resource.TestCheckResourceAttr("identitynow_source.test", "name", "test source update"),
					resource.TestCheckResourceAttr("identitynow_source.test", "description", "Creating from integration tests update"),
					resource.TestCheckResourceAttr("identitynow_source.test", "owner.type", "IDENTITY"),
					resource.TestCheckResourceAttr("identitynow_source.test", "owner.id", "ownerIdUpd"),
					resource.TestCheckResourceAttr("identitynow_source.test", "cluster.id", "clusterIdUpd"),
					resource.TestCheckResourceAttr("identitynow_source.test", "account_correlation_config.type", "ACCOUNT_CORRELATION_CONFIG"),
					resource.TestCheckResourceAttr("identitynow_source.test", "account_correlation_config.id", "accCorIdUpd"),
					resource.TestCheckResourceAttr("identitynow_source.test", "account_correlation_rule.type", "RULE"),
					resource.TestCheckResourceAttr("identitynow_source.test", "account_correlation_rule.id", "accCorRuleIdUpd"),
					resource.TestCheckResourceAttr("identitynow_source.test", "manager_correlation_mapping.account_attribute_name", "manAccAttrUpd"),
					resource.TestCheckResourceAttr("identitynow_source.test", "manager_correlation_mapping.identity_attribute_name", "manIdentAttrUpd"),
					resource.TestCheckResourceAttr("identitynow_source.test", "manager_correlation_rule.type", "RULE"),
					resource.TestCheckResourceAttr("identitynow_source.test", "manager_correlation_rule.id", "manCorRuleIdUpd"),
					resource.TestCheckResourceAttr("identitynow_source.test", "before_provisioning_rule.type", "RULE"),
					resource.TestCheckResourceAttr("identitynow_source.test", "before_provisioning_rule.id", "befProvRuleIdUpd"),
					resource.TestCheckResourceAttr("identitynow_source.test", "features.0", "AUTHENTICATE"),
					resource.TestCheckResourceAttr("identitynow_source.test", "features.1", "ENABLE"),
					resource.TestCheckResourceAttr("identitynow_source.test", "type", "connectorType"),
					resource.TestCheckResourceAttr("identitynow_source.test", "connector", "custom connector"),
					resource.TestCheckResourceAttr("identitynow_source.test", "connector_attributes", "{\"enableLCS\":true,\"inherited\":{\"first\":\"1Upd\",\"seconds\":\"2\"}}"),
					resource.TestCheckResourceAttr("identitynow_source.test", "connector_attributes_credentials", "{\"inherited\":{\"clientId\":\"clientId\",\"clientSecret\":\"clientSecretUpd\"}}"),
					resource.TestCheckResourceAttr("identitynow_source.test", "delete_threshold", "15"),
					resource.TestCheckResourceAttr("identitynow_source.test", "management_workgroup.type", "GOVERNANCE_GROUP"),
					resource.TestCheckResourceAttr("identitynow_source.test", "management_workgroup.id", "manWorkIdUpd"),
					resource.TestCheckResourceAttr("identitynow_source.test", "connector_files.0", "ojdbc10-19.18.0.0.txt"),
					resource.TestCheckResourceAttr("identitynow_source.test", "connector_files.1", "ojdbc10-19.21.0.0.txt"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
