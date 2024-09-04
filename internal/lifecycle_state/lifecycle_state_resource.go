package lifecycle_state

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	sailpoint "github.com/sailpoint-oss/golang-sdk/v2"
	sailpointV3 "github.com/sailpoint-oss/golang-sdk/v2/api_v3"
	"terraform-provider-identitynow/internal/patch"
	"terraform-provider-identitynow/internal/sailpoint/custom"
	"terraform-provider-identitynow/internal/util"
)

var (
	_ resource.Resource              = &lifeCycleResource{}
	_ resource.ResourceWithConfigure = &lifeCycleResource{}
)

func NewLifecycleStateResource() resource.Resource {
	return &lifeCycleResource{}
}

type lifeCycleResource struct {
	apiClient *sailpoint.APIClient
}

func (r *lifeCycleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *lifeCycleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_lifecycle_state"
}

func (r *lifeCycleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "System-generated unique ID of the Object",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"identity_profile_id": schema.StringAttribute{
				Description: "Identity Profile ID",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name of the Object",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether the lifecycle state is enabled or disabled",
				Required:    true,
			},
			"technical_name": schema.StringAttribute{
				Description: "The technical name for lifecycle state. This is for internal use",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Description: "Lifecycle state description.",
				Optional:    true,
			},
			"email_notification_option": schema.SingleNestedAttribute{
				Description: "This is used for representing email configuration for a lifecycle state",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"notify_managers": schema.BoolAttribute{
						Description: "If true, then the manager is notified of the lifecycle state change",
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
					},
					"notify_all_admins": schema.BoolAttribute{
						Description: "If true, then all the admins are notified of the lifecycle state change",
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
					},
					"notify_specific_users": schema.BoolAttribute{
						Description: "If true, then the users specified in \"email_address_list\" below are notified of lifecycle state change",
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
					},
					"email_address_list": schema.ListAttribute{
						Description: "List of user email addresses. If \"notify_specific_users\" option is true, then these users are notified of lifecycle state change",
						Optional:    true,
						Computed:    true,
						ElementType: types.StringType,
						Default:     listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{})),
					},
				},
			},
			"account_actions": schema.ListNestedAttribute{
				Description: "This is used for representing email configuration for a lifecycle state",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"action": schema.StringAttribute{
							Description: "Describes if action will be enabled or disabled",
							Required:    true,
							Validators: []validator.String{
								stringvalidator.OneOf("ENABLE", "DISABLE"),
							},
						},
						"source_ids": schema.ListAttribute{
							Description: "List of unique source IDs. The sources must have the ENABLE feature or flat file source. See \"/sources\" endpoint for source features",
							Required:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
			"access_profile_ids": schema.ListAttribute{
				Description: "List of unique access-profile IDs that are associated with the lifecycle state",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Default:     listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{})),
			},
		},
	}
}

func (r *lifeCycleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	diagnostics := resp.Diagnostics
	var plan lifecycleStateModel
	diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if diagnostics.HasError() {
		return
	}
	lifecycleState := r.convertToAPIModel(&plan, &diagnostics)
	tflog.Info(ctx, fmt.Sprintf("Creating LifeCycle State: %s", util.PrettyPrint(lifecycleState)))
	if diagnostics.HasError() {
		return
	}
	identityProfileId := plan.IdentityProfileId.ValueString()

	existing := r.findExisting(ctx, &plan, &diagnostics)
	if existing == nil {
		lifecycleStateResp, spResp, err := r.apiClient.V3.LifecycleStatesAPI.CreateLifecycleState(ctx, identityProfileId).LifecycleState(lifecycleState).Execute()
		if err != nil {
			diagnostics.AddError(
				"Error Creating Lifecycle State",
				"Could not create Lifecycle State '"+plan.Name.ValueString()+"': "+err.Error()+"\n"+util.GetBody(spResp),
			)
			return
		}
		r.mapToTerraformModel(&plan, lifecycleStateResp, &diagnostics)
	} else {
		if diagnostics.HasError() {
			return
		}
		jsonPatch, err := patch.NewLifecycleStatePatchBuilder(&lifecycleState, existing).GenerateJsonPatch()
		if err != nil {
			diagnostics.AddError(
				"Error Generating Update Patch",
				"Could not generate update patch for Lifecycle State '"+plan.Name.ValueString()+"': "+err.Error(),
			)
			return
		}
		v3JsonPatch, err := patch.ConvertFromBetaToV3(jsonPatch)
		if err != nil {
			diagnostics.AddError(
				"Error Generating Update Patch",
				"Could not convert patch to V3 for Lifecycle State '"+plan.Name.ValueString()+"': "+err.Error(),
			)
			return
		}
		lifecycleStateResp, spResp, err := r.apiClient.V3.LifecycleStatesAPI.UpdateLifecycleStates(ctx, plan.IdentityProfileId.ValueString(), *existing.Id).JsonPatchOperation(v3JsonPatch).Execute()
		if err != nil {
			diagnostics.AddError(
				"Error Updating Lifecycle State",
				"Could not update Lifecycle State '"+plan.Name.ValueString()+"': "+err.Error()+"\n"+util.GetBody(spResp),
			)
			return
		}
		r.mapToTerraformModel(&plan, lifecycleStateResp, &diagnostics)
	}

	if diagnostics.HasError() {
		return
	}

	// Set state to fully populated data
	diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *lifeCycleResource) findExisting(ctx context.Context, plan *lifecycleStateModel, diagnostics *diag.Diagnostics) *sailpointV3.LifecycleState {
	identityProfileId := plan.IdentityProfileId.ValueString()
	lifecycleStateResp, spResp, err := r.apiClient.V3.LifecycleStatesAPI.ListLifecycleStates(ctx, identityProfileId).Execute()
	if err != nil {
		diagnostics.AddError(
			"Error Creating Lifecycle State",
			"Error during listing of Lifecycle State of '"+identityProfileId+"': "+err.Error()+"\n"+util.GetBody(spResp),
		)
		return nil
	}
	for _, lifecycleState := range lifecycleStateResp {
		if lifecycleState.Name == plan.Name.ValueString() {
			return &lifecycleState
		}
	}
	return nil
}

func (r *lifeCycleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state lifecycleStateModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	lifecycleState, spResp, err := r.apiClient.V3.LifecycleStatesAPI.GetLifecycleState(ctx, state.IdentityProfileId.ValueString(), state.Id.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Lifecycle State",
			"Could not read Lifecycle State '"+state.Name.ValueString()+"': "+err.Error()+"\n"+util.GetBody(spResp),
		)
		return
	}

	r.mapToTerraformModel(&state, lifecycleState, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *lifeCycleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state lifecycleStateModel
	diagnostics := resp.Diagnostics
	diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	diagnostics.Append(req.State.Get(ctx, &state)...)
	if diagnostics.HasError() {
		return
	}

	newModel := r.convertToAPIModel(&plan, &diagnostics)
	if diagnostics.HasError() {
		return
	}
	oldModel := r.convertToAPIModel(&state, &diagnostics)
	if diagnostics.HasError() {
		return
	}
	jsonPatch, err := patch.NewLifecycleStatePatchBuilder(&newModel, &oldModel).GenerateJsonPatch()
	if err != nil {
		diagnostics.AddError(
			"Error Generating Update Patch",
			"Could not generate update patch for Lifecycle State '"+plan.Name.ValueString()+"': "+err.Error(),
		)
		return
	}
	v3JsonPatch, err := patch.ConvertFromBetaToV3(jsonPatch)
	if err != nil {
		diagnostics.AddError(
			"Error Generating Update Patch",
			"Could not convert patch to V3 for Lifecycle State '"+plan.Name.ValueString()+"': "+err.Error(),
		)
		return
	}
	tflog.Info(ctx, fmt.Sprintf("Updating LifeCycle State '%s': %s", state.Name.ValueString(), util.PrettyPrint(v3JsonPatch)))
	lifecycleStateResp, spResp, err := r.apiClient.V3.LifecycleStatesAPI.UpdateLifecycleStates(ctx, state.IdentityProfileId.ValueString(), state.Id.ValueString()).JsonPatchOperation(v3JsonPatch).Execute()
	if err != nil {
		diagnostics.AddError(
			"Error Updating Lifecycle State",
			"Could not update Lifecycle State '"+plan.Name.ValueString()+"': "+err.Error()+"\n"+util.GetBody(spResp),
		)
		return
	}

	r.mapToTerraformModel(&plan, lifecycleStateResp, &diagnostics)
	if diagnostics.HasError() {
		return
	}

	// Set refreshed state
	diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *lifeCycleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state lifecycleStateModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, spResp, err := r.apiClient.V3.LifecycleStatesAPI.DeleteLifecycleState(ctx, state.IdentityProfileId.ValueString(), state.Id.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Lifecycle State",
			"Could not delete Lifecycle State '"+state.Name.ValueString()+"': "+err.Error()+"\n"+util.GetBody(spResp),
		)
		return
	}
}

func (r *lifeCycleResource) convertToAPIModel(model *lifecycleStateModel, _ *diag.Diagnostics) sailpointV3.LifecycleState {
	var emailNotificationOption *sailpointV3.EmailNotificationOption
	if model.EmailNotificationOption != nil {
		emailAddressList := make([]string, len(model.EmailNotificationOption.EmailAddressList))
		for i, v := range model.EmailNotificationOption.EmailAddressList {
			emailAddressList[i] = v.ValueString()
		}
		emailNotificationOption = &sailpointV3.EmailNotificationOption{
			NotifyManagers:      model.EmailNotificationOption.NotifyManagers.ValueBoolPointer(),
			NotifyAllAdmins:     model.EmailNotificationOption.NotifyAllAdmins.ValueBoolPointer(),
			NotifySpecificUsers: model.EmailNotificationOption.NotifySpecificUsers.ValueBoolPointer(),
			EmailAddressList:    emailAddressList,
		}
	}
	var accountActions []sailpointV3.AccountAction
	if model.AccountActions != nil {
		accountActions = make([]sailpointV3.AccountAction, len(model.AccountActions))
		for i, v := range model.AccountActions {
			sourceIds := make([]string, len(v.SourceIds))
			for j, w := range v.SourceIds {
				sourceIds[j] = w.ValueString()
			}
			accountActions[i] = sailpointV3.AccountAction{
				Action:    v.Action.ValueStringPointer(),
				SourceIds: sourceIds,
			}
		}
	}
	var accessProfileIds []string
	if model.AccessProfileIds != nil {
		accessProfileIds = make([]string, len(model.AccessProfileIds))
		for i, v := range model.AccessProfileIds {
			accessProfileIds[i] = v.ValueString()
		}
	}
	return sailpointV3.LifecycleState{
		Name:                    model.Name.ValueString(),
		Enabled:                 model.Enabled.ValueBoolPointer(),
		TechnicalName:           model.TechnicalName.ValueString(),
		Description:             model.Description.ValueStringPointer(),
		EmailNotificationOption: emailNotificationOption,
		AccountActions:          accountActions,
		AccessProfileIds:        accessProfileIds,
	}
}

func (r *lifeCycleResource) mapToTerraformModel(model *lifecycleStateModel, lifecycleState *sailpointV3.LifecycleState, _ *diag.Diagnostics) {
	model.Id = types.StringPointerValue(lifecycleState.Id)
	model.Name = types.StringValue(lifecycleState.Name)
	model.Enabled = types.BoolPointerValue(lifecycleState.Enabled)
	model.TechnicalName = types.StringValue(lifecycleState.TechnicalName)
	model.Description = types.StringPointerValue(lifecycleState.Description)
	if lifecycleState.EmailNotificationOption != nil {
		model.EmailNotificationOption = &emailNotificationOption{
			NotifyManagers:      types.BoolPointerValue(lifecycleState.EmailNotificationOption.NotifyManagers),
			NotifyAllAdmins:     types.BoolPointerValue(lifecycleState.EmailNotificationOption.NotifyAllAdmins),
			NotifySpecificUsers: types.BoolPointerValue(lifecycleState.EmailNotificationOption.NotifySpecificUsers),
			EmailAddressList:    make([]types.String, len(lifecycleState.EmailNotificationOption.EmailAddressList)),
		}
		for i, v := range lifecycleState.EmailNotificationOption.EmailAddressList {
			model.EmailNotificationOption.EmailAddressList[i] = types.StringValue(v)
		}
	}
	if lifecycleState.AccountActions != nil && len(lifecycleState.AccountActions) > 0 {
		model.AccountActions = make([]accountAction, len(lifecycleState.AccountActions))
		for i, v := range lifecycleState.AccountActions {
			model.AccountActions[i] = accountAction{
				Action:    types.StringPointerValue(v.Action),
				SourceIds: make([]types.String, len(v.SourceIds)),
			}
			for j, w := range v.SourceIds {
				model.AccountActions[i].SourceIds[j] = types.StringValue(w)
			}
		}
	}
	if lifecycleState.AccessProfileIds != nil && len(lifecycleState.AccessProfileIds) > 0 {
		model.AccessProfileIds = make([]types.String, len(lifecycleState.AccessProfileIds))
		for i, v := range lifecycleState.AccessProfileIds {
			model.AccessProfileIds[i] = types.StringValue(v)
		}
	}
}
