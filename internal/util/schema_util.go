package util

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ResourceReferenceSchema(allowedType string, required bool, description string) schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description: description,
		Required:    required,
		Optional:    !required,
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Computed: true,
				Optional: true,
				Default:  stringdefault.StaticString(allowedType),
				Validators: []validator.String{
					stringvalidator.OneOf(allowedType),
				},
			},
			"id": schema.StringAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func ResourceReferenceListSchema(allowedType string, required bool, description string) schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Description:  description,
		Required:     required,
		Optional:     !required,
		NestedObject: ResourceReferenceNestedObject(allowedType),
	}
}

func ResourceReferenceSetSchema(allowedType string, required bool, description string) schema.SetNestedAttribute {
	return schema.SetNestedAttribute{
		Description:  description,
		Required:     required,
		Optional:     !required,
		NestedObject: ResourceReferenceNestedObject(allowedType),
	}
}

func ResourceReferenceNestedObject(allowedType string) schema.NestedAttributeObject {
	return schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Computed: true,
				Optional: true,
				Default:  stringdefault.StaticString(allowedType),
				Validators: []validator.String{
					stringvalidator.OneOf(allowedType),
				},
			},
			"id": schema.StringAttribute{
				Required: true,
			},
			"name": schema.StringAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func ResourceNameReferenceSchema(allowedType string, required bool, description string) schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description: description,
		Required:    required,
		Optional:    !required,
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Computed: true,
				Optional: true,
				Default:  stringdefault.StaticString(allowedType),
				Validators: []validator.String{
					stringvalidator.OneOf(allowedType),
				},
			},
			"id": schema.StringAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

func ResourceListNameReferenceSchema(allowedType string) schema.NestedAttributeObject {
	return schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Computed: true,
				Optional: true,
				Default:  stringdefault.StaticString(allowedType),
				Validators: []validator.String{
					stringvalidator.OneOf(allowedType),
				},
			},
			"id": schema.StringAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

func NewPointerReferenceModel(theType *string, id *string, name *string) *ReferenceModel {
	return &ReferenceModel{
		Type: types.StringPointerValue(theType),
		Id:   types.StringPointerValue(id),
		Name: types.StringPointerValue(name),
	}
}

func NewReferenceModel(theType string, id string, name string) *ReferenceModel {
	return &ReferenceModel{
		Type: types.StringValue(theType),
		Id:   types.StringValue(id),
		Name: types.StringValue(name),
	}
}

type ReferenceModel struct {
	Type types.String `tfsdk:"type"`
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func ConvertReferenceModelToMap(ctx context.Context, value, target interface{}) diag.Diagnostics {
	return tfsdk.ValueFrom(ctx, value, types.SetType{
		ElemType: types.ObjectType{
			AttrTypes: ReferenceModelAttrTypes(),
		},
	}, target)
}

func ReferenceModelAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"type": types.StringType,
		"id":   types.StringType,
		"name": types.StringType,
	}
}
