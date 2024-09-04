package source_schema

import (
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-identitynow/internal/util"
)

type sourceSchemaModel struct {
	Id                 types.String     `tfsdk:"id"`
	SourceId           types.String     `tfsdk:"source_id"`
	Name               types.String     `tfsdk:"name"`
	NativeObjectType   types.String     `tfsdk:"native_object_type"`
	IdentityAttribute  types.String     `tfsdk:"identity_attribute"`
	DisplayAttribute   types.String     `tfsdk:"display_attribute"`
	HierarchyAttribute types.String     `tfsdk:"hierarchy_attribute"`
	IncludePermissions types.Bool       `tfsdk:"include_permissions"`
	Features           []types.String   `tfsdk:"features"`
	Configuration      jsontypes.Exact  `tfsdk:"configuration"`
	Attributes         []attributeModel `tfsdk:"attributes"`
}

type attributeModel struct {
	Name          types.String         `tfsdk:"name"`
	Type          types.String         `tfsdk:"type"`
	Schema        *util.ReferenceModel `tfsdk:"schema"`
	Description   types.String         `tfsdk:"description"`
	IsMulti       types.Bool           `tfsdk:"is_multi"`
	IsEntitlement types.Bool           `tfsdk:"is_entitlement"`
	IsGroup       types.Bool           `tfsdk:"is_group"`
}
