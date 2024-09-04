//go:build integration

package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"
)

func TestIntegration_SourceAggregationScheduleResource_AccountToEntitlement(t *testing.T) {
	t.Skip("aggregation API was disabled by SailPoint")
	sourceCloud := getSources(1, context.Background())[0].ConnectorAttributes["cloudExternalId"].(string)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerIntegrationConfig + `
resource "identitynow_source_aggregation_schedule" "test_account_schedule" {
  source_cloud_id  = ` + sourceCloud + `
  cron_expression  = "0 5 0 * * ?"
  aggregation_type = "account"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("identitynow_source_aggregation_schedule.test_account_schedule", "source_cloud_id", sourceCloud),
					resource.TestCheckResourceAttr("identitynow_source_aggregation_schedule.test_account_schedule", "cron_expression", "0 5 0 * * ?"),
					resource.TestCheckResourceAttr("identitynow_source_aggregation_schedule.test_account_schedule", "aggregation_type", "account"),
				),
			},
			// Update and Read testing
			{
				Config: providerIntegrationConfig + `
resource "identitynow_source_aggregation_schedule" "test_account_schedule" {
  source_cloud_id  = ` + sourceCloud + `
  cron_expression  = "0 4 0 * * ?"
  aggregation_type = "entitlement"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("identitynow_source_aggregation_schedule.test_account_schedule", "source_cloud_id", sourceCloud),
					resource.TestCheckResourceAttr("identitynow_source_aggregation_schedule.test_account_schedule", "cron_expression", "0 4 0 * * ?"),
					resource.TestCheckResourceAttr("identitynow_source_aggregation_schedule.test_account_schedule", "aggregation_type", "entitlement"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestIntegration_SourceAggregationScheduleResource_EntitlementToAccount(t *testing.T) {
	t.Skip("aggregation API was disabled by SailPoint")
	sourceCloud := getSources(1, context.Background())[0].ConnectorAttributes["cloudExternalId"].(string)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerIntegrationConfig + `
resource "identitynow_source_aggregation_schedule" "test_account_schedule" {
  source_cloud_id  = ` + sourceCloud + `
  cron_expression  = "0 5 0 * * ?"
  aggregation_type = "entitlement"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("identitynow_source_aggregation_schedule.test_account_schedule", "source_cloud_id", sourceCloud),
					resource.TestCheckResourceAttr("identitynow_source_aggregation_schedule.test_account_schedule", "cron_expression", "0 5 0 * * ?"),
					resource.TestCheckResourceAttr("identitynow_source_aggregation_schedule.test_account_schedule", "aggregation_type", "entitlement"),
				),
			},
			// Update and Read testing
			{
				Config: providerIntegrationConfig + `
resource "identitynow_source_aggregation_schedule" "test_account_schedule" {
  source_cloud_id  = ` + sourceCloud + `
  cron_expression  = "0 4 0 * * ?"
  aggregation_type = "account"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("identitynow_source_aggregation_schedule.test_account_schedule", "source_cloud_id", sourceCloud),
					resource.TestCheckResourceAttr("identitynow_source_aggregation_schedule.test_account_schedule", "cron_expression", "0 4 0 * * ?"),
					resource.TestCheckResourceAttr("identitynow_source_aggregation_schedule.test_account_schedule", "aggregation_type", "account"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestIntegration_SourceAggregationScheduleResource_AccountSourceAndCronUpdate(t *testing.T) {
	t.Skip("aggregation API was disabled by SailPoint")
	sourceCloud := getSources(2, context.Background())
	originalSourceCloudId := sourceCloud[0].ConnectorAttributes["cloudExternalId"].(string)
	updatedSourceCloudId := sourceCloud[1].ConnectorAttributes["cloudExternalId"].(string)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerIntegrationConfig + `
resource "identitynow_source_aggregation_schedule" "test_account_schedule" {
  source_cloud_id  = ` + originalSourceCloudId + `
  cron_expression  = "0 5 0 * * ?"
  aggregation_type = "account"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("identitynow_source_aggregation_schedule.test_account_schedule", "source_cloud_id", originalSourceCloudId),
					resource.TestCheckResourceAttr("identitynow_source_aggregation_schedule.test_account_schedule", "cron_expression", "0 5 0 * * ?"),
					resource.TestCheckResourceAttr("identitynow_source_aggregation_schedule.test_account_schedule", "aggregation_type", "account"),
				),
			},
			// Update and Read testing
			{
				Config: providerIntegrationConfig + `
resource "identitynow_source_aggregation_schedule" "test_account_schedule" {
  source_cloud_id  = ` + updatedSourceCloudId + `
  cron_expression  = "0 4 0 * * ?"
  aggregation_type = "account"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("identitynow_source_aggregation_schedule.test_account_schedule", "source_cloud_id", updatedSourceCloudId),
					resource.TestCheckResourceAttr("identitynow_source_aggregation_schedule.test_account_schedule", "cron_expression", "0 4 0 * * ?"),
					resource.TestCheckResourceAttr("identitynow_source_aggregation_schedule.test_account_schedule", "aggregation_type", "account"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestIntegration_SourceAggregationScheduleResource_EntitlementSourceAndCronUpdate(t *testing.T) {
	t.Skip("aggregation API was disabled by SailPoint")
	sourceCloud := getSources(2, context.Background())
	originalSourceCloudId := sourceCloud[0].ConnectorAttributes["cloudExternalId"].(string)
	updatedSourceCloudId := sourceCloud[1].ConnectorAttributes["cloudExternalId"].(string)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerIntegrationConfig + `
resource "identitynow_source_aggregation_schedule" "test_account_schedule" {
  source_cloud_id  = ` + originalSourceCloudId + `
  cron_expression  = "0 5 0 * * ?"
  aggregation_type = "entitlement"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("identitynow_source_aggregation_schedule.test_account_schedule", "source_cloud_id", originalSourceCloudId),
					resource.TestCheckResourceAttr("identitynow_source_aggregation_schedule.test_account_schedule", "cron_expression", "0 5 0 * * ?"),
					resource.TestCheckResourceAttr("identitynow_source_aggregation_schedule.test_account_schedule", "aggregation_type", "entitlement"),
				),
			},
			// Update and Read testing
			{
				Config: providerIntegrationConfig + `
resource "identitynow_source_aggregation_schedule" "test_account_schedule" {
  source_cloud_id  = ` + updatedSourceCloudId + `
  cron_expression  = "0 4 0 * * ?"
  aggregation_type = "entitlement"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("identitynow_source_aggregation_schedule.test_account_schedule", "source_cloud_id", updatedSourceCloudId),
					resource.TestCheckResourceAttr("identitynow_source_aggregation_schedule.test_account_schedule", "cron_expression", "0 4 0 * * ?"),
					resource.TestCheckResourceAttr("identitynow_source_aggregation_schedule.test_account_schedule", "aggregation_type", "entitlement"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
