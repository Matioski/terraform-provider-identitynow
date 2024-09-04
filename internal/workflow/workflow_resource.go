package workflow

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	sailpoint "github.com/sailpoint-oss/golang-sdk/v2"
	sailpointBeta "github.com/sailpoint-oss/golang-sdk/v2/api_beta"
	"net/http"
	"terraform-provider-identitynow/internal/patch"
	"terraform-provider-identitynow/internal/sailpoint/custom"
	"terraform-provider-identitynow/internal/util"
)

var (
	_ resource.Resource              = &workflowResource{}
	_ resource.ResourceWithConfigure = &workflowResource{}
)

func NewWorkflowResource() resource.Resource {
	return &workflowResource{}
}

type workflowResource struct {
	apiClient *sailpoint.APIClient
}

func (r *workflowResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *workflowResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workflow"
}

func (r *workflowResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the workflow",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the workflow",
				Required:    true,
			},
			"owner": schema.SingleNestedAttribute{
				Description: "The identity that owns the workflow. The owner's permissions in IDN will determine what actions the workflow is allowed to perform. Ownership can be changed by updating the owner in a PUT or PATCH request.",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Computed: true,
						Optional: true,
						Default:  stringdefault.StaticString("IDENTITY"),
						Validators: []validator.String{
							stringvalidator.OneOf("IDENTITY"),
						},
					},
					"id": schema.StringAttribute{
						Computed: true,
						Optional: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
							stringplanmodifier.RequiresReplace(),
						},
					},
					"name": schema.StringAttribute{
						Computed: true,
						Optional: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
							stringplanmodifier.RequiresReplace(),
						},
					},
				},
			},
			"description": schema.StringAttribute{
				Description: "The description of the rule's purpose",
				Optional:    true,
			},
			"definition": schema.SingleNestedAttribute{
				Description: "The map of steps that the workflow will execute",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"start": schema.StringAttribute{
						Description: "The name of the starting step",
						Optional:    true,
					},
					"steps": schema.StringAttribute{
						Description: "One or more step objects that comprise this workflow. Please see the Workflow documentation to see the JSON schema for each step type",
						Optional:    true,
						CustomType:  jsontypes.ExactType{},
					},
				},
			},
			"enabled": schema.BoolAttribute{
				Description: "Enable or disable the workflow. Workflows cannot be created in an enabled state",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"trigger": schema.SingleNestedAttribute{
				Description: "The trigger that starts the workflow",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Description: "The type of trigger",
						Required:    true,
						Validators: []validator.String{
							stringvalidator.OneOf("EVENT", "EXTERNAL", "SCHEDULED"),
						},
					},
					"attributes": schema.SingleNestedAttribute{
						Description: "Workflow Trigger Attributes. One of the following sets of attributes is required, depending on the trigger type",
						Required:    true,
						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								Description: "The ID of the trigger. EVENT trigger type",
								Optional:    true,
							},
							"filter": schema.StringAttribute{
								Description: "JSON path expression that will limit which events the trigger will fire on. EVENT trigger type",
								Optional:    true,
							},

							"name": schema.StringAttribute{
								Description: "A unique name for the external trigger. EXTERNAL trigger type",
								Optional:    true,
							},
							"description": schema.StringAttribute{
								Description: "Additonal context about the external trigger. EXTERNAL trigger type",
								Optional:    true,
							},

							"cron_string": schema.StringAttribute{
								Description: "A valid CRON expression. SCHEDULED trigger type",
								Optional:    true,
							},
						},
					},
				},
			},
		},
	}
}

func (r *workflowResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan workflowModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	workflow := r.convertToAPIModel(&plan, &resp.Diagnostics)
	tflog.Info(ctx, "Creating Workflow "+util.PrettyPrint(workflow))
	workflowResp, spResp, err := r.apiClient.Beta.WorkflowsAPI.CreateWorkflow(ctx).CreateWorkflowRequest(workflow).Execute()
	if err != nil && spResp.StatusCode != http.StatusCreated {
		resp.Diagnostics.AddError(
			"Error Creating Workflow",
			"Could not create Workflow '"+plan.Name.ValueString()+"': "+err.Error()+"\n"+util.GetBody(spResp),
		)
		return
	}

	r.mapToTerraformModel(&plan, workflowResp, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *workflowResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state workflowModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	workflowResp, spResp, err := r.apiClient.Beta.WorkflowsAPI.GetWorkflow(ctx, state.Id.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Workflow",
			"Could not read Workflow '"+state.Name.ValueString()+"': "+err.Error()+"\n"+util.GetBody(spResp),
		)
		return
	}

	r.mapToTerraformModel(&state, workflowResp, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *workflowResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state workflowModel
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

	jsonPatch, err := patch.NewWorkflowPatchBuilder(&newModel, &oldModel).GenerateJsonPatch()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Generating Update Patch",
			"Could not generate update patch for Workflow '"+plan.Name.ValueString()+"': "+err.Error(),
		)
		return
	}
	tflog.Info(ctx, "Updating Workflow "+util.PrettyPrint(jsonPatch))

	workflowResp, spResp, err := r.apiClient.Beta.WorkflowsAPI.PatchWorkflow(ctx, state.Id.ValueString()).JsonPatchOperation(jsonPatch).Execute()
	if err != nil || (spResp != nil && spResp.StatusCode != http.StatusOK) {
		resp.Diagnostics.AddError(
			"Error Updating Workflow",
			"Could not update Workflow '"+plan.Name.ValueString()+"': "+err.Error()+"\n"+util.GetBody(spResp),
		)
		return
	}

	r.mapToTerraformModel(&plan, workflowResp, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *workflowResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state workflowModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueString()
	if state.Enabled.ValueBool() == true {
		value := sailpointBeta.BoolAsJsonPatchOperationValue(sailpointBeta.PtrBool(false))
		tflog.Info(ctx, "Disabling Workflow '"+state.Name.ValueString()+"' before deletion")
		_, spResp, err := r.apiClient.Beta.WorkflowsAPI.PatchWorkflow(context.Background(), id).JsonPatchOperation([]sailpointBeta.JsonPatchOperation{
			{
				Op:    "replace",
				Path:  "/enabled",
				Value: &value,
			},
		}).Execute()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Disabling Workflow",
				"Could not disable Workflow '"+state.Name.ValueString()+"': "+err.Error()+"\n"+util.GetBody(spResp),
			)
			return
		}
	}
	spResp, err := r.apiClient.Beta.WorkflowsAPI.DeleteWorkflow(ctx, id).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Workflow",
			"Could not delete Workflow '"+state.Name.ValueString()+"': "+err.Error()+"\n"+util.GetBody(spResp),
		)
		return
	}
}

func (r *workflowResource) convertToAPIModel(model *workflowModel, diagnostics *diag.Diagnostics) sailpointBeta.CreateWorkflowRequest {
	var modelDef *sailpointBeta.WorkflowDefinition
	if model.Definition != nil {
		steps := util.UnmarshalJsonType(model.Definition.Steps, diagnostics)
		if diagnostics.HasError() {
			return sailpointBeta.CreateWorkflowRequest{}
		}
		modelDef = &sailpointBeta.WorkflowDefinition{
			Start: model.Definition.Start.ValueStringPointer(),
			Steps: steps,
		}
	}
	var wfTrigger *sailpointBeta.WorkflowTrigger
	if model.Trigger != nil {
		attributes := util.ConvertTFModelToMap(model.Trigger.Attributes)
		wfTrigger = &sailpointBeta.WorkflowTrigger{
			Type:       model.Trigger.Type.ValueString(),
			Attributes: attributes,
		}
	}
	return sailpointBeta.CreateWorkflowRequest{
		Name: model.Name.ValueString(),
		Owner: sailpointBeta.WorkflowBodyOwner{
			Id:   model.Owner.Id.ValueStringPointer(),
			Name: model.Owner.Name.ValueStringPointer(),
			Type: model.Owner.Type.ValueStringPointer(),
		},
		Description: model.Description.ValueStringPointer(),
		Definition:  modelDef,
		Enabled:     model.Enabled.ValueBoolPointer(),
		Trigger:     wfTrigger,
	}
}

func (r *workflowResource) mapToTerraformModel(model *workflowModel, workflow *sailpointBeta.Workflow, diagnostic *diag.Diagnostics) {
	model.Id = types.StringPointerValue(workflow.Id)
	model.Name = types.StringPointerValue(workflow.Name)
	model.Owner = *util.NewPointerReferenceModel(workflow.Owner.Type, workflow.Owner.Id, workflow.Owner.Name)
	model.Description = types.StringPointerValue(workflow.Description)
	model.Enabled = types.BoolPointerValue(workflow.Enabled)
	if workflow.Definition != nil {
		model.Definition = &definition{
			Start: types.StringPointerValue(workflow.Definition.Start),
			Steps: util.MarshalToJsonType(workflow.Definition.Steps, diagnostic),
		}
	}
	if workflow.Trigger != nil && workflow.Trigger.Type != "" {
		model.Trigger = &trigger{
			Type:       types.StringValue(workflow.Trigger.Type),
			Attributes: r.convertToTriggerAttributes(workflow.Trigger.Attributes),
		}
	}
}

func (r *workflowResource) convertToTriggerAttributes(attributes map[string]interface{}) triggerAttributes {
	return triggerAttributes{
		Id:          r.getTFString(attributes, "id"),
		Filter:      r.getTFString(attributes, "filter.$"),
		Name:        r.getTFString(attributes, "name"),
		Description: r.getTFString(attributes, "description"),
		CronString:  r.getTFString(attributes, "cronString"),
	}
}

func (r *workflowResource) getTFString(attributes map[string]interface{}, key string) types.String {
	if val, ok := attributes[key]; ok {
		return types.StringValue(val.(string))
	}
	return types.StringNull()
}
