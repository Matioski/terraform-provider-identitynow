package workflow

import (
	"terraform-provider-identitynow/internal/util"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type workflowModel struct {
	Id          types.String        `tfsdk:"id"`
	Name        types.String        `tfsdk:"name"`
	Owner       util.ReferenceModel `tfsdk:"owner"`
	Description types.String        `tfsdk:"description"`
	Enabled     types.Bool          `tfsdk:"enabled"`
	Definition  *definition         `tfsdk:"definition"`
	Trigger     *trigger            `tfsdk:"trigger"`
}

type definition struct {
	Start types.String    `tfsdk:"start"`
	Steps jsontypes.Exact `tfsdk:"steps"`
}

type trigger struct {
	Type       types.String      `tfsdk:"type"`
	Attributes triggerAttributes `tfsdk:"attributes"`
}

// One of values sets is applicable
type triggerAttributes struct {
	// Fields for EVENT type of trigger
	Id                types.String `tfsdk:"id" json:"id"`
	Filter            types.String `tfsdk:"filter" json:"filter.$"`
	AttributeToFilter types.String `tfsdk:"attribute_to_filter" json:"attributeToFilter"`

	// Fields for EXTERNAL type of trigger
	Name        types.String `tfsdk:"name" json:"name"`
	Description types.String `tfsdk:"description" json:"description"`

	// Fields for SCHEDULED type of trigger
	CronString types.String `tfsdk:"cron_string" json:"cronString"`
}
