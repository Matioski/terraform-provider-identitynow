package identity_attribute

import (
	"context"
	"fmt"
	"terraform-provider-identitynow/internal/sailpoint/custom"
	"terraform-provider-identitynow/internal/util"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	sailpoint "github.com/sailpoint-oss/golang-sdk/v2"
	sailpoint_beta "github.com/sailpoint-oss/golang-sdk/v2/api_beta"
)

var (
	_ resource.Resource              = &identityAttributeResource{}
	_ resource.ResourceWithConfigure = &identityAttributeResource{}
)

func NewIdentityAttributeResource() resource.Resource {
	return &identityAttributeResource{}
}

type identityAttributeResource struct {
	apiClient *sailpoint.APIClient
}

func (r *identityAttributeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *identityAttributeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_identity_attribute"
}

func (r *identityAttributeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The technical name of the identity attribute.",
				Required:    true,
			},
			"display_name": schema.StringAttribute{
				Description: "The business-friendly name of the identity attribute.",
				Required:    true,
			},
			"standard": schema.BoolAttribute{
				Description: "The business-friendly name of the identity attribute.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"type": schema.StringAttribute{
				Description: "The type of the identity attribute",
				Required:    true,
			},
			"multi": schema.BoolAttribute{
				Description: "Shows if the identity attribute is multi-valued",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"searchable": schema.BoolAttribute{
				Description: "Shows if the identity attribute is searchable",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"system": schema.BoolAttribute{
				Description: "Shows this is 'system' identity attribute that does not have a source and is not configurable.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"sources": schema.ListNestedAttribute{
				Description: "List of sources for an attribute, this specifies how the value of the rule is derived",
				Required:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Description: "The type of the source",
							Required:    true,
						},
						"properties": schema.StringAttribute{
							Description: "The source properties",
							CustomType:  jsontypes.ExactType{},
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

func (r *identityAttributeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan identityAttributeModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	idAttr := r.convertToAPIModel(&plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, fmt.Sprintf("Creating Identity Attribute: %s", util.PrettyPrint(idAttr)))
	identityAttribute, spResp, err := r.apiClient.Beta.IdentityAttributesAPI.CreateIdentityAttribute(ctx).IdentityAttribute(idAttr).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Identity Attribute",
			"Could not create Identity Attribute, unexpected error: "+err.Error()+"\n"+util.GetBody(spResp),
		)
		return
	}

	r.mapToTerraformModel(&plan, identityAttribute, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *identityAttributeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state identityAttributeModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := state.Name.ValueString()
	identityAttribute, spResp, err := r.apiClient.Beta.IdentityAttributesAPI.GetIdentityAttribute(ctx, name).Execute()
	if spResp.StatusCode == 404 {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Identity Attribute",
			"Could not read Identity Attribute '"+name+"': "+err.Error()+"\n"+util.GetBody(spResp),
		)
		return
	}
	tflog.Info(ctx, fmt.Sprintf("Identity Attribute: %s", util.PrettyPrint(identityAttribute)))

	r.mapToTerraformModel(&state, identityAttribute, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *identityAttributeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan identityAttributeModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	idAttr := r.convertToAPIModel(&plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	name := plan.Name.ValueString()
	identityAttribute, spResp, err := r.apiClient.Beta.IdentityAttributesAPI.PutIdentityAttribute(ctx, name).IdentityAttribute(idAttr).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Identity Attribute",
			"Could not update Identity Attribute '"+name+"': "+err.Error()+"\n"+util.GetBody(spResp),
		)
		return
	}

	r.mapToTerraformModel(&plan, identityAttribute, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *identityAttributeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state identityAttributeModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	spResp, err := r.apiClient.Beta.IdentityAttributesAPI.DeleteIdentityAttribute(ctx, state.Name.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Identity Attribute",
			"Could not delete Identity Attribute '"+state.Name.ValueString()+"': "+err.Error()+"\n"+util.GetBody(spResp),
		)
		return
	}
}

func (r *identityAttributeResource) convertToAPIModel(tfModel *identityAttributeModel, diagnostics *diag.Diagnostics) sailpoint_beta.IdentityAttribute {
	var sources []sailpoint_beta.Source1
	for _, source := range tfModel.Sources {
		properties := util.UnmarshalJsonType(source.Properties, diagnostics)
		if diagnostics.HasError() {
			return sailpoint_beta.IdentityAttribute{}
		}
		sources = append(sources, sailpoint_beta.Source1{
			Type:       source.Type.ValueStringPointer(),
			Properties: properties,
		})
	}

	idAttr := sailpoint_beta.IdentityAttribute{
		Name:        tfModel.Name.ValueString(),
		DisplayName: tfModel.DisplayName.ValueStringPointer(),
		Standard:    tfModel.Standard.ValueBoolPointer(),
		Type:        *sailpoint_beta.NewNullableString(tfModel.Type.ValueStringPointer()),
		Multi:       tfModel.Multi.ValueBoolPointer(),
		Searchable:  tfModel.Searchable.ValueBoolPointer(),
		System:      tfModel.System.ValueBoolPointer(),
		Sources:     sources,
	}
	return idAttr
}

func (r *identityAttributeResource) mapToTerraformModel(tfModel *identityAttributeModel, identityAttribute *sailpoint_beta.IdentityAttribute, diagnostics *diag.Diagnostics) {
	tfModel.Name = types.StringValue(identityAttribute.Name)
	tfModel.DisplayName = types.StringPointerValue(identityAttribute.DisplayName)
	tfModel.Standard = types.BoolPointerValue(identityAttribute.Standard)
	tfModel.Type = types.StringPointerValue(identityAttribute.Type.Get())
	tfModel.Multi = types.BoolPointerValue(identityAttribute.Multi)
	tfModel.Searchable = types.BoolPointerValue(identityAttribute.Searchable)
	tfModel.System = types.BoolPointerValue(identityAttribute.System)

	tfModel.Sources = make([]identityAttributeSourceModel, len(identityAttribute.Sources))
	for idx, source := range identityAttribute.Sources {
		properties := util.MarshalToJsonType(source.Properties, diagnostics)
		if diagnostics.HasError() {
			return
		}
		tfModel.Sources[idx] = identityAttributeSourceModel{
			Type:       types.StringPointerValue(source.Type),
			Properties: properties,
		}
	}
}
