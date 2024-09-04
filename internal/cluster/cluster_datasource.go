package cluster

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	sailpoint "github.com/sailpoint-oss/golang-sdk/v2"
	sailpoint_beta "github.com/sailpoint-oss/golang-sdk/v2/api_beta"
	"terraform-provider-identitynow/internal/sailpoint/custom"
	"terraform-provider-identitynow/internal/util"
)

var (
	_ datasource.DataSource              = &clusterDataSource{}
	_ datasource.DataSourceWithConfigure = &clusterDataSource{}
)

func NewClusterDataSource() datasource.DataSource {
	return &clusterDataSource{}
}

type clusterDataSource struct {
	apiClient *sailpoint.APIClient
}

func (d *clusterDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*custom.APIClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *sailpoint.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.apiClient = client.ApiClient
}

func (d *clusterDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cluster"
}

func (d *clusterDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Unique Cluster identifier attribute.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Cluster name.",
				Required:    true,
			},
			"pod": schema.StringAttribute{
				Description: "Cluster pod.",
				Computed:    true,
			},
			"org": schema.StringAttribute{
				Description: "Cluster organization.",
				Computed:    true,
			},
		},
	}
}

func (d *clusterDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config clusterModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterName := config.Name.ValueString()
	managedClusters, spResp, err := d.apiClient.Beta.ManagedClustersAPI.GetManagedClusters(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Cluster '"+clusterName+"'",
			"Could not read Cluster '"+clusterName+"': "+err.Error()+"\n"+util.GetBody(spResp),
		)
		return
	}
	var cluster *sailpoint_beta.ManagedCluster
	for _, item := range managedClusters {
		if item.Name != nil && *item.Name == clusterName {
			cluster = &item
			break
		}
	}
	if cluster == nil {
		resp.Diagnostics.AddError(
			"Unable to Read Cluster '"+clusterName+"'",
			"Could not find Cluster '"+clusterName+"':\n"+util.GetBody(spResp),
		)
		return
	}

	model := clusterModel{
		Id:   types.StringValue(cluster.Id),
		Name: types.StringPointerValue(cluster.Name),
		Pod:  types.StringPointerValue(cluster.Pod),
		Org:  types.StringPointerValue(cluster.Org),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}
