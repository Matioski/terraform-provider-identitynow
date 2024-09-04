package identity_profile

import (
	"context"
	"fmt"
	patch "terraform-provider-identitynow/internal/patch"
	"terraform-provider-identitynow/internal/sailpoint/custom"
	"terraform-provider-identitynow/internal/util"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	sailpoint "github.com/sailpoint-oss/golang-sdk/v2"
	sailpoint_beta "github.com/sailpoint-oss/golang-sdk/v2/api_beta"
)

// Implementation of IdentityNow Identity Profiles CRUD - https://developer.sailpoint.com/idn/api/beta/identity-profiles

var (
	_ resource.Resource              = &identityProfileResource{}
	_ resource.ResourceWithConfigure = &identityProfileResource{}
)

func NewIdentityProfileResource() resource.Resource {
	return &identityProfileResource{}
}

type identityProfileResource struct {
	apiClient *sailpoint.APIClient
}

func (r *identityProfileResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*custom.APIClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *sailpoint.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.apiClient = client.ApiClient
}

func (r *identityProfileResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_identity_profile"
}

func (r *identityProfileResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the Identity Profile",
				Optional:    true,
			},
			"owner": util.ResourceReferenceSchema("IDENTITY", false, "The owner of the Identity Profile"),
			"priority": schema.Int64Attribute{
				Description: "The priority for an Identity Profile",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"authoritative_source": util.ResourceReferenceSchema("SOURCE", true, "The authoritative source for this Identity Profile"),
			"identity_attribute_config": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						Description: "If the profile or mapping is enabled",
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(true),
					},
					"attribute_transforms": schema.ListNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"identity_attribute_name": schema.StringAttribute{
									Description: "Name of the identity attribute",
									Required:    true,
								},
								"transform_definition": schema.SingleNestedAttribute{
									Description: "The seaspray transformation definition",
									Optional:    true,
									Attributes: map[string]schema.Attribute{
										"type": schema.StringAttribute{
											Description: "The type of the transform definition",
											Required:    true,
										},
										"attributes": schema.StringAttribute{
											Description: "Arbitrary key-value pairs to store any metadata for the object",
											Optional:    true,
											CustomType:  jsontypes.ExactType{},
										},
									},
								},
							},
						},
					},
				},
			},
			"identity_exception_report_reference": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"task_result_id": schema.StringAttribute{
						Description: "The id of the task result",
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"report_name": schema.StringAttribute{
						Description: "The name of the report",
						Optional:    true,
						Computed:    true,
					},
				},
			},
		},
	}
}

func (r *identityProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan identityProfileModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	identityProfile := r.convertToAPIModel(&plan, &resp.Diagnostics)
	tflog.Info(ctx, "Creating Identity Profile: "+util.PrettyPrint(identityProfile))
	identityProfileResponse, spResp, err := r.apiClient.Beta.IdentityProfilesAPI.CreateIdentityProfile(ctx).IdentityProfile(identityProfile).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Identity Profile",
			"Could not create Identity Profile '"+plan.Name.ValueString()+"': "+err.Error()+"\n"+util.GetBody(spResp),
		)
		return
	}

	r.mapToTerraformModel(&plan, identityProfileResponse, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *identityProfileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state identityProfileModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	identityProfile, spResp, err := r.apiClient.Beta.IdentityProfilesAPI.GetIdentityProfile(ctx, state.Id.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Identity Profile",
			"Could not read Identity Profile '"+state.Name.ValueString()+"': "+err.Error()+"\n"+util.GetBody(spResp),
		)
		return
	}
	tflog.Info(ctx, fmt.Sprintf("Identity Profile: %s", util.PrettyPrint(identityProfile)))
	r.mapToTerraformModel(&state, identityProfile, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *identityProfileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state identityProfileModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	newModel := r.convertToAPIModel(&plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	oldModel := r.convertToAPIModel(&state, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	patchOperations, err := patch.NewIdentityProfilePatchBuilder(&newModel, &oldModel).GenerateJsonPatch()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Generating Update Patch",
			"Could not generate update patch for Identity Profile '"+plan.Name.ValueString()+"': "+err.Error(),
		)
		return
	}
	tflog.Info(ctx, fmt.Sprintf("Patch: %s", util.PrettyPrint(patchOperations)))
	identityProfileResp, spResp, err := r.apiClient.Beta.IdentityProfilesAPI.UpdateIdentityProfile(ctx, plan.Id.ValueString()).JsonPatchOperation(patchOperations).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Identity Profile",
			"Could not update Identity Profile '"+state.Name.ValueString()+"': "+err.Error()+"\n"+util.GetBody(spResp),
		)
		return
	}
	tflog.Info(ctx, fmt.Sprintf("Updated Identity Profile: %s", util.PrettyPrint(identityProfileResp)))

	r.mapToTerraformModel(&plan, identityProfileResp, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *identityProfileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state identityProfileModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, spResp, err := r.apiClient.Beta.IdentityProfilesAPI.DeleteIdentityProfile(ctx, state.Id.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Identity Profile",
			"Could not delete Identity Profile '"+state.Name.ValueString()+"': "+err.Error()+"\n"+util.GetBody(spResp),
		)
		return
	}
}

func (r *identityProfileResource) convertToAPIModel(tfModel *identityProfileModel, diagnostics *diag.Diagnostics) sailpoint_beta.IdentityProfile {
	owner := sailpoint_beta.NullableIdentityProfileAllOfOwner{}
	if tfModel.Owner != nil {
		owner.Set(&sailpoint_beta.IdentityProfileAllOfOwner{
			Type: util.GetTFStringPointer(tfModel.Owner.Type),
			Id:   util.GetTFStringPointer(tfModel.Owner.Id),
			Name: util.GetTFStringPointer(tfModel.Owner.Name),
		})
	}
	var identityAttributeConfig *sailpoint_beta.IdentityAttributeConfig
	if tfModel.IdentityAttributeConfig != nil {
		attributeTransforms := make([]sailpoint_beta.IdentityAttributeTransform, len(tfModel.IdentityAttributeConfig.AttributeTransforms))
		for key, value := range tfModel.IdentityAttributeConfig.AttributeTransforms {
			attributes := util.UnmarshalJsonType(value.TransformDefinition.Attributes, diagnostics)
			if diagnostics.HasError() {
				return sailpoint_beta.IdentityProfile{}
			}
			attributeTransforms[key] = sailpoint_beta.IdentityAttributeTransform{
				IdentityAttributeName: util.GetTFStringPointer(value.IdentityAttributeName),
				TransformDefinition: &sailpoint_beta.TransformDefinition{
					Type:       util.GetTFStringPointer(value.TransformDefinition.Type),
					Attributes: attributes,
				},
			}
		}
		identityAttributeConfig = &sailpoint_beta.IdentityAttributeConfig{
			Enabled:             tfModel.IdentityAttributeConfig.Enabled.ValueBoolPointer(),
			AttributeTransforms: attributeTransforms,
		}
	}
	identityExceptionReportReference := sailpoint_beta.NullableIdentityExceptionReportReference{}
	if tfModel.IdentityExceptionReportReference != nil {
		identityExceptionReportReference.Set(&sailpoint_beta.IdentityExceptionReportReference{
			TaskResultId: util.GetTFStringPointer(tfModel.IdentityExceptionReportReference.TaskResultId),
			ReportName:   util.GetTFStringPointer(tfModel.IdentityExceptionReportReference.ReportName),
		})
	}
	return sailpoint_beta.IdentityProfile{
		Name:        tfModel.Name.ValueString(),
		Description: *sailpoint_beta.NewNullableString(util.GetTFStringPointer(tfModel.Description)),
		Owner:       owner,
		Priority:    tfModel.Priority.ValueInt64Pointer(),
		AuthoritativeSource: sailpoint_beta.IdentityProfileAllOfAuthoritativeSource{
			Type: util.GetTFStringPointer(tfModel.AuthoritativeSource.Type),
			Id:   util.GetTFStringPointer(tfModel.AuthoritativeSource.Id),
			Name: util.GetTFStringPointer(tfModel.AuthoritativeSource.Name),
		},
		IdentityAttributeConfig:          identityAttributeConfig,
		IdentityExceptionReportReference: identityExceptionReportReference,
	}
}

func (r *identityProfileResource) mapToTerraformModel(tfModel *identityProfileModel, identityProfile *sailpoint_beta.IdentityProfile, diagnostics *diag.Diagnostics) {
	tfModel.Id = types.StringPointerValue(identityProfile.Id)
	tfModel.Name = types.StringValue(identityProfile.Name)
	tfModel.Description = types.StringPointerValue(identityProfile.Description.Get())
	if identityProfile.Owner.IsSet() {
		tfModel.Owner = util.NewPointerReferenceModel(identityProfile.Owner.Get().Type, identityProfile.Owner.Get().Id, identityProfile.Owner.Get().Name)
	}
	tfModel.Priority = types.Int64PointerValue(identityProfile.Priority)
	tfModel.AuthoritativeSource = *util.NewPointerReferenceModel(identityProfile.AuthoritativeSource.Type, identityProfile.AuthoritativeSource.Id, identityProfile.AuthoritativeSource.Name)

	if identityProfile.IdentityAttributeConfig != nil {
		attributeTransforms := make([]attributeTransformModel, len(identityProfile.IdentityAttributeConfig.AttributeTransforms))
		for key, value := range identityProfile.IdentityAttributeConfig.AttributeTransforms {
			attributeTransforms[key] = attributeTransformModel{
				IdentityAttributeName: types.StringPointerValue(value.IdentityAttributeName),
				TransformDefinition: &transformDefinitionModel{
					Type:       types.StringPointerValue(value.TransformDefinition.Type),
					Attributes: util.MarshalToJsonType(value.TransformDefinition.Attributes, diagnostics),
				},
			}
		}
		if diagnostics.HasError() {
			return
		}
		tfModel.IdentityAttributeConfig = &identityAttributeConfigModel{
			Enabled:             types.BoolPointerValue(identityProfile.IdentityAttributeConfig.Enabled),
			AttributeTransforms: attributeTransforms,
		}
	}
	if identityProfile.IdentityExceptionReportReference.IsSet() && identityProfile.IdentityExceptionReportReference.Get() != nil {
		tfModel.IdentityExceptionReportReference = &identityExceptionReportReferenceModel{
			TaskResultId: types.StringPointerValue(identityProfile.IdentityExceptionReportReference.Get().TaskResultId),
			ReportName:   types.StringPointerValue(identityProfile.IdentityExceptionReportReference.Get().ReportName),
		}
	}
}
