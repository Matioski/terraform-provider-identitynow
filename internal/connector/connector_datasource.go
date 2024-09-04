package connector

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	sailpoint "github.com/sailpoint-oss/golang-sdk/v2"
	"terraform-provider-identitynow/internal/sailpoint/custom"
	"terraform-provider-identitynow/internal/util"
)

var (
	_ datasource.DataSource              = &connectorDataSource{}
	_ datasource.DataSourceWithConfigure = &connectorDataSource{}
)

func NewConnectorDataSource() datasource.DataSource {
	return &connectorDataSource{}
}

type connectorDataSource struct {
	apiClient *sailpoint.APIClient
}

func (d *connectorDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *connectorDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connector"
}

func (d *connectorDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The connector name",
				Required:    true,
			},
			"type": schema.StringAttribute{
				Description: "The connector type",
				Computed:    true,
			},
			"script_name": schema.StringAttribute{
				Description: "The connector script name",
				Computed:    true,
			},
		},
	}
}

func (d *connectorDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config connectorModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	name := config.Name.ValueString()
	connectors, spResp, err := d.apiClient.Beta.ConnectorsAPI.GetConnectorList(ctx).Limit(1).Filters("name sw \"" + name + "\"").Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Connector '"+name+"'",
			"Could not read Connector '"+name+"': "+err.Error()+"\n"+util.GetBody(spResp),
		)
		return
	}
	if len(connectors) != 1 {
		resp.Diagnostics.AddError(
			"Unable to Read Connectors '"+name+"'",
			"Did not found Connector with name '"+name+"'",
		)
		return
	}
	model := connectorModel{
		Name:       types.StringPointerValue(connectors[0].Name),
		Type:       types.StringPointerValue(connectors[0].Type),
		ScriptName: types.StringPointerValue(connectors[0].ScriptName),
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}
