package identity

import "github.com/hashicorp/terraform-plugin-framework/types"

type identityModel struct {
	Id    types.String `tfsdk:"id"`
	Alias types.String `tfsdk:"alias"`
	Name  types.String `tfsdk:"name"`
	Email types.String `tfsdk:"email"`
}
