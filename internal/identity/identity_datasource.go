package identity

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
	_ datasource.DataSource              = &identityDataSource{}
	_ datasource.DataSourceWithConfigure = &identityDataSource{}
)

func NewIdentityDataSource() datasource.DataSource {
	return &identityDataSource{}
}

type identityDataSource struct {
	apiClient *sailpoint.APIClient
}

func (d *identityDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *identityDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_identity"
}

func (d *identityDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Unique Identity identifier.",
				Computed:    true,
			},
			"alias": schema.StringAttribute{
				Description: "Unique Identity name.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "Identity display mame",
				Computed:    true,
			},
			"email": schema.StringAttribute{
				Description: "Identity email",
				Computed:    true,
			},
		},
	}
}

func (d *identityDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config identityModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	identityName := config.Alias.ValueString()
	identities, spResp, err := d.apiClient.Beta.IdentitiesAPI.ListIdentities(ctx).Limit(1).Filters("alias eq \"" + identityName + "\"").Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Identity '"+identityName+"'",
			"Could not read Identity '"+identityName+"': "+err.Error()+"\n"+util.GetBody(spResp),
		)
		return
	}
	if len(identities) != 1 {
		resp.Diagnostics.AddError(
			"Unable to Read Identity '"+identityName+"'",
			"Did not found with name '"+identityName+"'",
		)
		return
	}
	model := identityModel{
		Id:    types.StringPointerValue(identities[0].Id),
		Alias: types.StringPointerValue(identities[0].Alias),
		Name:  types.StringValue(identities[0].Name),
		Email: types.StringPointerValue(identities[0].EmailAddress.Get()),
	}

	// Set state
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}
