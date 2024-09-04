package identity_profile

import (
	"terraform-provider-identitynow/internal/util"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type identityProfileModel struct {
	Id                               types.String                           `tfsdk:"id"`
	Name                             types.String                           `tfsdk:"name"`
	Description                      types.String                           `tfsdk:"description"`
	Owner                            *util.ReferenceModel                   `tfsdk:"owner"`
	Priority                         types.Int64                            `tfsdk:"priority"`
	AuthoritativeSource              util.ReferenceModel                    `tfsdk:"authoritative_source"`
	IdentityAttributeConfig          *identityAttributeConfigModel          `tfsdk:"identity_attribute_config"`
	IdentityExceptionReportReference *identityExceptionReportReferenceModel `tfsdk:"identity_exception_report_reference"`
}

type identityAttributeConfigModel struct {
	Enabled             types.Bool                `tfsdk:"enabled"`
	AttributeTransforms []attributeTransformModel `tfsdk:"attribute_transforms"`
}
type attributeTransformModel struct {
	IdentityAttributeName types.String              `tfsdk:"identity_attribute_name"`
	TransformDefinition   *transformDefinitionModel `tfsdk:"transform_definition"`
}

type transformDefinitionModel struct {
	Type       types.String    `tfsdk:"type"`
	Attributes jsontypes.Exact `tfsdk:"attributes"`
}

type identityExceptionReportReferenceModel struct {
	TaskResultId types.String `tfsdk:"task_result_id"`
	ReportName   types.String `tfsdk:"report_name"`
}
