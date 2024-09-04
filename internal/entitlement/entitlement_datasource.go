package entitlement

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
	_ datasource.DataSource              = &entitlementDataSource{}
	_ datasource.DataSourceWithConfigure = &entitlementDataSource{}
)

func NewEntitlementDataSource() datasource.DataSource {
	return &entitlementDataSource{}
}

type entitlementDataSource struct {
	apiClient *sailpoint.APIClient
}

func (d *entitlementDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *entitlementDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_entitlement"
}

func (d *entitlementDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The entitlement id",
				Computed:    true,
			},
			"source_id": schema.StringAttribute{
				Description: "Source ID of the entitlement",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "The entitlement name",
				Optional:    true,
				Computed:    true,
			},
			"attribute": schema.StringAttribute{
				Description: "The entitlement attribute name",
				Optional:    true,
				Computed:    true,
			},
			"value": schema.StringAttribute{
				Description: "The value of the entitlement",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func (d *entitlementDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var model entitlementModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}
	filters, err := d.generateFilters(&model)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Entitlement",
			"Could not generate filters: "+err.Error(),
		)
		return
	}
	entitlements, spResp, err := d.apiClient.Beta.EntitlementsAPI.ListEntitlements(ctx).Filters(filters).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Entitlement",
			"Could not read entitlement '"+filters+"': "+err.Error()+"\n"+util.GetBody(spResp),
		)
		return
	}
	if len(entitlements) != 1 {
		resp.Diagnostics.AddError(
			"Unable to Read Entitlement",
			"List with filter '"+filters+"' returned "+fmt.Sprint(len(entitlements))+" entitlements, expected 1",
		)
		return
	}
	model = entitlementModel{
		Id:        types.StringPointerValue(entitlements[0].Id),
		Name:      types.StringPointerValue(entitlements[0].Name),
		Attribute: types.StringPointerValue(entitlements[0].Attribute.Get()),
		Value:     types.StringPointerValue(entitlements[0].Value),
		SourceId:  types.StringPointerValue(entitlements[0].Source.Id),
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (d *entitlementDataSource) generateFilters(model *entitlementModel) (string, error) {
	if model.Name.ValueString() == "" && model.Value.ValueString() == "" {
		return "", fmt.Errorf("either 'name' or 'value' must be set")
	}
	filter := "source.id eq \"" + model.SourceId.ValueString() + "\""
	if model.Name.ValueString() != "" {
		filter += " and name eq \"" + model.Name.ValueString() + "\""
	}
	if model.Value.ValueString() != "" {
		filter += " and value eq \"" + model.Value.ValueString() + "\""
	}
	if model.Attribute.ValueString() != "" {
		filter += " and attribute eq \"" + model.Attribute.ValueString() + "\""
	}
	return filter, nil
}
