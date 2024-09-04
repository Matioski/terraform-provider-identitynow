package connector

import "github.com/hashicorp/terraform-plugin-framework/types"

type connectorModel struct {
	Name       types.String `tfsdk:"name"`
	Type       types.String `tfsdk:"type"`
	ScriptName types.String `tfsdk:"script_name"`
}
