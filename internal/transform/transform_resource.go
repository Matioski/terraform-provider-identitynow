package transform

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	sailpoint "github.com/sailpoint-oss/golang-sdk/v2"
	sailpoint_v3 "github.com/sailpoint-oss/golang-sdk/v2/api_v3"
	"terraform-provider-identitynow/internal/sailpoint/custom"
	"terraform-provider-identitynow/internal/util"
)

var (
	_ resource.Resource                = &transformResource{}
	_ resource.ResourceWithConfigure   = &transformResource{}
	_ resource.ResourceWithImportState = &transformResource{}
)

func NewTransformResource() resource.Resource {
	return &transformResource{}
}

type transformResource struct {
	apiClient *sailpoint.APIClient
}

func (r *transformResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *transformResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_transform"
}

func (r *transformResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(50),
				},
			},
			"type": schema.StringAttribute{
				Required: true,
			},
			"attributes": schema.StringAttribute{
				CustomType: jsontypes.ExactType{},
				Required:   true,
			},
		},
	}
}

func (r *transformResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan transformModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	transform := r.convertToAPIModel(&plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	transformResp, spResp, err := r.apiClient.V3.TransformsAPI.CreateTransform(ctx).Transform(transform).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Transform",
			"Could not create Transform '"+plan.Name.ValueString()+"': "+err.Error()+"\n"+util.GetBody(spResp),
		)
		return
	}

	plan.Id = types.StringValue(transformResp.Id)

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *transformResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state transformModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	transformRead, spResp, err := r.apiClient.V3.TransformsAPI.GetTransform(ctx, state.Id.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Transform",
			"Could not read Transform '"+state.Name.ValueString()+"': "+err.Error()+"\n"+util.GetBody(spResp),
		)
		return
	}
	tflog.Info(ctx, fmt.Sprintf("Transform: %s", util.PrettyPrint(transformRead)))

	r.mapToTerraformModel(&state, transformRead, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *transformResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan transformModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	transform := r.convertToAPIModel(&plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Updating IDN Transform '"+plan.Id.ValueString()+"'")
	transformUpdate, spResp, err := r.apiClient.V3.TransformsAPI.UpdateTransform(ctx, plan.Id.ValueString()).Transform(transform).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Transform",
			"Could not update Transform '"+plan.Name.ValueString()+"': "+err.Error()+"\n"+util.GetBody(spResp),
		)
		return
	}

	r.mapToTerraformModel(&plan, transformUpdate, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *transformResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state transformModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	spResp, err := r.apiClient.V3.TransformsAPI.DeleteTransform(ctx, state.Id.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Transform",
			"Could not delete Transform '"+state.Name.ValueString()+"': "+err.Error()+"\n"+util.GetBody(spResp),
		)
		return
	}
}

func (r *transformResource) convertToAPIModel(plan *transformModel, diagnostics *diag.Diagnostics) sailpoint_v3.Transform {
	attributes := util.UnmarshalJsonType(plan.Attributes, diagnostics)
	if diagnostics.HasError() {
		return sailpoint_v3.Transform{}
	}
	transform := sailpoint_v3.Transform{
		Name:       plan.Name.ValueString(),
		Type:       plan.Type.ValueString(),
		Attributes: attributes,
	}
	return transform
}

func (r *transformResource) mapToTerraformModel(tfModel *transformModel, transformRead *sailpoint_v3.TransformRead, diagnostics *diag.Diagnostics) {
	tfModel.Id = types.StringValue(transformRead.Id)
	tfModel.Name = types.StringValue(transformRead.Name)
	tfModel.Type = types.StringValue(transformRead.Type)
	tfModel.Attributes = util.MarshalToJsonType(transformRead.Attributes, diagnostics)
}

func (r *transformResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
