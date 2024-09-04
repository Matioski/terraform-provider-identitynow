package transform

import (
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type transformModel struct {
	Id         types.String    `tfsdk:"id"`
	Name       types.String    `tfsdk:"name"`
	Type       types.String    `tfsdk:"type"`
	Attributes jsontypes.Exact `tfsdk:"attributes"`
}
