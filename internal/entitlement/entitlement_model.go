package entitlement

import "github.com/hashicorp/terraform-plugin-framework/types"

type entitlementModel struct {
	Id        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	Attribute types.String `tfsdk:"attribute"`
	Value     types.String `tfsdk:"value"`
	SourceId  types.String `tfsdk:"source_id"`
}
