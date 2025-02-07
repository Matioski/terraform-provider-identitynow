package org_config

import (
	"context"
	"fmt"
	"strings"
	"terraform-provider-identitynow/internal/patch"
	"terraform-provider-identitynow/internal/sailpoint/custom"
	"terraform-provider-identitynow/internal/util"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	sailpoint "github.com/sailpoint-oss/golang-sdk/v2"
	sailpointBeta "github.com/sailpoint-oss/golang-sdk/v2/api_beta"
)

var (
	_ resource.Resource              = &orgConfigResource{}
	_ resource.ResourceWithConfigure = &orgConfigResource{}
)

const (
	defaultTimeZone = "UTC"
)

func NewOrgConfigResource() resource.Resource {
	return &orgConfigResource{}
}

type orgConfigResource struct {
	apiClient *sailpoint.APIClient
}

func (r *orgConfigResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *orgConfigResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_org_config"
}

func (r *orgConfigResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"time_zone": schema.StringAttribute{
				Description: "The selected time zone which is to be used for the org. This directly affects when scheduled tasks are executed. Valid options can be found at /beta/org-config/valid-time-zones",
				Required:    true,
			},
		},
	}
}

func (r *orgConfigResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state orgConfigModel

	req.Plan.Get(ctx, &plan)
	state = r.doRead(ctx, &resp.Diagnostics)

	newModel := r.convertToAPIModel(&plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	oldModel := r.convertToAPIModel(&state, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	result := r.doUpdate(ctx, &newModel, &oldModel, &resp.Diagnostics)

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, result)...)
}

func (r *orgConfigResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	resp.Diagnostics.Append(resp.State.Set(ctx, r.doRead(ctx, &resp.Diagnostics))...)
}

func (r *orgConfigResource) doRead(ctx context.Context, diagnostics *diag.Diagnostics) orgConfigModel {
	orgConfig, spResp, err := r.apiClient.Beta.OrgConfigAPI.GetOrgConfig(ctx).Execute()
	if err != nil {
		diagnostics.AddError(
			"Error Reading Organization Configuration",
			"Could not read Organization Configuration: "+err.Error()+"\n"+util.GetBody(spResp),
		)
		return orgConfigModel{}
	}

	state := orgConfigModel{
		TimeZone: types.StringPointerValue(orgConfig.TimeZone),
	}

	return state
}

func (r *orgConfigResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state orgConfigModel

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

	result := r.doUpdate(ctx, &newModel, &oldModel, &resp.Diagnostics)

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, result)...)
}

func (r *orgConfigResource) doUpdate(ctx context.Context, newModel *sailpointBeta.OrgConfig, oldModel *sailpointBeta.OrgConfig, diagnostics *diag.Diagnostics) orgConfigModel {
	jsonPatch, err := patch.NewOrgConfigPatchBuilder(newModel, oldModel).GenerateJsonPatch()
	// if there is nothing to patch, return the new model back
	if jsonPatch == nil && err == nil {
		return orgConfigModel{
			TimeZone: types.StringPointerValue(newModel.TimeZone),
		}
	}

	if err != nil {
		diagnostics.AddError(
			"Patch error: Error Generating Update Patch",
			"Could not generate update patch for Organization Configuration ': "+err.Error(),
		)
		return orgConfigModel{}
	}
	orgConfigPatchResp, spResp, err := r.apiClient.Beta.OrgConfigAPI.PatchOrgConfig(ctx).JsonPatchOperation(jsonPatch).Execute()
	if err != nil {
		// Skip error from the SDK if the error is due to the armSapSystemIdMappings field
		if !strings.Contains(err.Error(), "OrgConfig.armSapSystemIdMappings") {
			diagnostics.AddError(
				"Error Updating Organization Configuration",
				"Could not update Organization Configuration ': "+err.Error()+"\n"+util.GetBody(spResp),
			)
			return orgConfigModel{}
		}
	}

	// It's also been observed that the response from the API, even when the error is nil, does not always contain the updated value, but instead returns nil or "" after a successful patch, so we need to check for that
	if orgConfigPatchResp == nil || orgConfigPatchResp.TimeZone == nil || *orgConfigPatchResp.TimeZone == "" {
		return orgConfigModel{
			TimeZone: types.StringPointerValue(newModel.TimeZone),
		}
	}

	result := orgConfigModel{
		TimeZone: types.StringPointerValue(orgConfigPatchResp.TimeZone),
	}
	return result
}

func (r *orgConfigResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *orgConfigResource) convertToAPIModel(tfModelToConvert *orgConfigModel, _ *diag.Diagnostics) sailpointBeta.OrgConfig {
	return sailpointBeta.OrgConfig{
		TimeZone: tfModelToConvert.TimeZone.ValueStringPointer(),
	}
}
