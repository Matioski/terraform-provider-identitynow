package connector_rule

import (
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type connectorRuleModel struct {
	Id          types.String    `tfsdk:"id"`
	Name        types.String    `tfsdk:"name"`
	Description types.String    `tfsdk:"description"`
	Type        types.String    `tfsdk:"type"`
	Signature   *signature      `tfsdk:"signature"`
	SourceCode  sourceCode      `tfsdk:"source_code"`
	Attributes  jsontypes.Exact `tfsdk:"attributes"`
}

type signature struct {
	Input  []signatureData `tfsdk:"input"`
	Output *signatureData  `tfsdk:"output"`
}

type signatureData struct {
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Type        types.String `tfsdk:"type"`
}

type sourceCode struct {
	Version types.String `tfsdk:"version"`
	Script  types.String `tfsdk:"script"`
}
