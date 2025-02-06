//go:build integration

package provider

import (
	"context"
	"fmt"
	"os"
	"strings"
	"terraform-provider-identitynow/internal/util"
	"time"

	"github.com/sailpoint-oss/golang-sdk/v2/api_v2024"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	sailpoint "github.com/sailpoint-oss/golang-sdk/v2"
)

const (
	// using env variables for the provider configuration
	providerIntegrationConfig = `
provider "identitynow" {

}
`
	ownerIdentityId           = "TODO"
	ownerIdentityName         = "TODO"
	ownerIdentityEmail        = "TODO"
	updatedOwnerIdentityId    = "TODO"
	updatedOwnerIdentityName  = "TODO"
	updatedOwnerIdentityEmail = "TODO"
	tenSeconds                = 10 * time.Second
	defaultTimeZone           = "UTC"
)

var (
	TestAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"identitynow": providerserver.NewProtocol6WithError(New("test")()),
	}

	configuration = sailpoint.NewConfiguration(sailpoint.ClientConfiguration{
		ClientId:     os.Getenv("IDN_CLIENT_ID"),
		ClientSecret: os.Getenv("IDN_CLIENT_SECRET"),
		BaseURL:      os.Getenv("IDN_HOST"),
		TokenURL:     os.Getenv("IDN_HOST") + "/oauth/token",
	})

	SPApiClient = sailpoint.NewAPIClient(configuration)
	_           = enableLogging()
	cisTasks    = []string{"Periodic Cloud Connector Updater", "Periodic workgroup Connection Sync Task", "Perform Nightly Sync", "Periodic Attribute Synchronization Provisioning", "Cleanup Orphan Indices", "IRS Housekeeper Task", "Perform provisioning activity search delete synchronization", "Periodic certification remediation scans", "Perform maintenance", "CMS Housekeeper tasks", "Periodic Cloud Orphan Task Result Reaper", "Purge Expired Indices from ES", "Check expired work items daily", "Periodic Entitlements Bulk Reconcile Delete", "Publish CIS and ES counts", "Periodic Access Profiles Bulk Synchronizer", "Perform identity search reconciliation", "Periodic Access Profiles Bulk Reconcile Delete", "Nightly Roles Bulk Synchronizer", "Periodic Entitlements Bulk Synchronizer", "Periodic Source Health Checker", "Nightly Cloud Housekeeping Task", "Nightly Entitlements Bulk Synchronizer", "Nightly Prune Identity Cubes", "Periodic Identity Bulk Delete Synchronizer", "Periodic Roles Bulk Reconcile Delete", "Periodic Consolidated Refresh", "Nightly Access Profiles Bulk Synchronizer", "Periodic Roles Bulk Synchronizer", "Prune audit events by type", "Nightly Manual Work Item Summary Generator", "Periodic Accounts Reconciliation", "Periodic Accounts Bulk Synchronize Deletes", "Background Object Termination"}
)

func enableLogging() bool {
	os.Setenv("TF_LOG", "INFO")
	return true
}

/*
*

	Checks for pending cis tasks based on a list provided by SailPoint. Pending cis tasks while running tests cause deletion of source to fail.
*/
func checkForPendingCisTask(ctx context.Context) {
	for true {
		pendingTasks, spResp, err := SPApiClient.Beta.TaskManagementAPI.GetPendingTasks(ctx).Execute()
		if err != nil {
			fmt.Printf("Error fetching pending tasks: %s\n%s\n", util.PrettyPrint(err), util.PrettyPrint(util.GetBody(spResp)))
		}
		if len(pendingTasks) == 0 {
			fmt.Printf("No pending tasks found\n")
			break
		}

		isCisTask := false
		for _, pendingTask := range pendingTasks {
			fmt.Printf("Evaluating pending task\n")
			for _, cisTask := range cisTasks {
				if strings.Contains(cisTask, pendingTask.UniqueName) {
					fmt.Printf("Pending Cis task found: %s\n", util.PrettyPrint(cisTask))
					time.Sleep(tenSeconds)
					isCisTask = true
					break
				}
			}
		}

		if !isCisTask {
			break
		}
	}
}

func getSources(limit int32, filters string) []api_v2024.Source {
	sources, spResp, err := SPApiClient.V2024.SourcesAPI.ListSources(context.Background()).Filters(filters).Limit(limit).Execute()
	if err != nil {
		fmt.Printf("Error fetching sources: %s\n%s\n", err, util.GetBody(spResp))
	}
	if len(sources) < int(limit) {
		fmt.Printf("Unable to provide %d sources; %d found\n", limit, len(sources))
	}
	return sources
}

func getAccessProfiles(limit int32) []api_v2024.AccessProfile {
	accessProfiles, spResp, err := SPApiClient.V2024.AccessProfilesAPI.ListAccessProfiles(context.Background()).Limit(limit).Execute()
	if err != nil {
		fmt.Printf("Error fetching access profiles: %s\n%s", err, util.GetBody(spResp))
	}
	if len(accessProfiles) < int(limit) {
		fmt.Printf("Unable to provide %d access profile(s); %d found\n", limit, len(accessProfiles))
	}
	return accessProfiles

}

func getManagedClusters(limit int32) []api_v2024.ManagedCluster {
	managedClusters, spResp, err := SPApiClient.V2024.ManagedClustersAPI.GetManagedClusters(context.Background()).Limit(limit).Execute()
	if err != nil {
		fmt.Printf("Error fetching managed cluster(s): %s\n%s", err, util.GetBody(spResp))
	}
	if len(managedClusters) < int(limit) {
		fmt.Printf("Unable to provide %d managed cluster(s); %d found\n", limit, len(managedClusters))
	}
	return managedClusters
}
