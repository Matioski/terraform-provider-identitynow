package identity_attribute

import (
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type identityAttributeModel struct {
	Name        types.String                   `tfsdk:"name"`
	DisplayName types.String                   `tfsdk:"display_name"`
	Standard    types.Bool                     `tfsdk:"standard"`
	Type        types.String                   `tfsdk:"type"`
	Multi       types.Bool                     `tfsdk:"multi"`
	Searchable  types.Bool                     `tfsdk:"searchable"`
	System      types.Bool                     `tfsdk:"system"`
	Sources     []identityAttributeSourceModel `tfsdk:"sources"`
}

type identityAttributeSourceModel struct {
	Type       types.String    `tfsdk:"type"`
	Properties jsontypes.Exact `tfsdk:"properties"`
}
