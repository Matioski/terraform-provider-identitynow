package connector_rule

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	sailpoint "github.com/sailpoint-oss/golang-sdk/v2"
	sailpointBeta "github.com/sailpoint-oss/golang-sdk/v2/api_beta"
	"terraform-provider-identitynow/internal/sailpoint/custom"
	"terraform-provider-identitynow/internal/util"
)

var (
	_ resource.Resource              = &connectorRuleResource{}
	_ resource.ResourceWithConfigure = &connectorRuleResource{}
)

func NewConnectorRuleResource() resource.Resource {
	return &connectorRuleResource{}
}

type connectorRuleResource struct {
	apiClient *sailpoint.APIClient
}

func (r *connectorRuleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *connectorRuleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connector_rule"
}

func (r *connectorRuleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the rule",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the rule",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Description: "The description of the rule's purpose",
				Optional:    true,
			},
			"type": schema.StringAttribute{
				Description: "The type of the rule",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"signature": schema.SingleNestedAttribute{
				Description: "The rule's function signature. Describes the rule's input arguments and output (if any)",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"input": schema.ListNestedAttribute{
						Description: "The input arguments of the rule",
						Required:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Description: "The name of the input argument",
									Required:    true,
								},
								"description": schema.StringAttribute{
									Description: "The description of the input argument",
									Optional:    true,
								},
								"type": schema.StringAttribute{
									Description: "The programmatic type of the input argument",
									Optional:    true,
								},
							},
						},
					},
					"output": schema.SingleNestedAttribute{
						Description: "The output of the rule",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"name": schema.StringAttribute{
								Description: "The name of the output argument",
								Required:    true,
							},
							"description": schema.StringAttribute{
								Description: "The description of the output argument",
								Optional:    true,
							},
							"type": schema.StringAttribute{
								Description: "The programmatic type of the output argument",
								Optional:    true,
							},
						},
					},
				},
			},
			"source_code": schema.SingleNestedAttribute{
				Description: "The rule's source code",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"version": schema.StringAttribute{
						Description: "The version of the source code",
						Required:    true,
					},
					"script": schema.StringAttribute{
						Description: "The source code of the rule",
						Required:    true,
					},
				},
			},
			"attributes": schema.StringAttribute{
				Description: "A map of string to objects",
				Optional:    true,
				CustomType:  jsontypes.ExactType{},
			},
		},
	}
}

func (r *connectorRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan connectorRuleModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	rule := r.convertToAPIModel(&plan, &resp.Diagnostics)
	ruleResp, spResp, err := r.apiClient.Beta.ConnectorRuleManagementAPI.CreateConnectorRule(ctx).ConnectorRuleCreateRequest(rule).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Connector Rule",
			"Could not create Connector Rule '"+plan.Name.ValueString()+"': "+err.Error()+"\n"+util.GetBody(spResp),
		)
		return
	}

	r.mapToTerraformModel(&plan, ruleResp, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *connectorRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state connectorRuleModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ruleResp, spResp, err := r.apiClient.Beta.ConnectorRuleManagementAPI.GetConnectorRule(ctx, state.Id.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Connector Rule",
			"Could not read Connector Rule '"+state.Name.ValueString()+"': "+err.Error()+"\n"+util.GetBody(spResp),
		)
		return
	}
	r.mapToTerraformModel(&state, ruleResp, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *connectorRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan connectorRuleModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	rule := r.convertToAPIUpdateModel(&plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	ruleResp, spResp, err := r.apiClient.Beta.ConnectorRuleManagementAPI.UpdateConnectorRule(ctx, plan.Id.ValueString()).ConnectorRuleUpdateRequest(rule).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Connector Rule",
			"Could not update Connector Rule '"+plan.Name.ValueString()+"': "+err.Error()+"\n"+util.GetBody(spResp),
		)
		return
	}

	r.mapToTerraformModel(&plan, ruleResp, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *connectorRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state connectorRuleModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	spResp, err := r.apiClient.Beta.ConnectorRuleManagementAPI.DeleteConnectorRule(ctx, state.Id.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Connector Rule",
			"Could not delete Connector Rule '"+state.Name.ValueString()+"': "+err.Error()+"\n"+util.GetBody(spResp),
		)
		return
	}
}

func (r *connectorRuleResource) convertToAPIModel(model *connectorRuleModel, diagnostics *diag.Diagnostics) sailpointBeta.ConnectorRuleCreateRequest {
	var signature *sailpointBeta.ConnectorRuleCreateRequestSignature
	if model.Signature != nil {
		inputs := make([]sailpointBeta.Argument, len(model.Signature.Input))
		for i, input := range model.Signature.Input {
			inputs[i] = sailpointBeta.Argument{
				Name:        input.Name.ValueString(),
				Description: input.Description.ValueStringPointer(),
				Type:        *sailpointBeta.NewNullableString(input.Type.ValueStringPointer()),
			}
		}
		var output *sailpointBeta.Argument
		if model.Signature.Output != nil {
			output = &sailpointBeta.Argument{
				Name:        model.Signature.Output.Name.ValueString(),
				Description: model.Signature.Output.Description.ValueStringPointer(),
				Type:        *sailpointBeta.NewNullableString(model.Signature.Output.Type.ValueStringPointer()),
			}
		}
		signature = &sailpointBeta.ConnectorRuleCreateRequestSignature{
			Input:  inputs,
			Output: *sailpointBeta.NewNullableArgument(output),
		}
	}

	return sailpointBeta.ConnectorRuleCreateRequest{
		Name:        model.Name.ValueString(),
		Description: model.Description.ValueStringPointer(),
		Type:        model.Type.ValueString(),
		Signature:   signature,
		SourceCode: sailpointBeta.SourceCode{
			Version: model.SourceCode.Version.ValueString(),
			Script:  model.SourceCode.Script.ValueString(),
		},
		Attributes: util.UnmarshalJsonType(model.Attributes, diagnostics),
	}
}

func (r *connectorRuleResource) convertToAPIUpdateModel(model *connectorRuleModel, diagnostics *diag.Diagnostics) sailpointBeta.ConnectorRuleUpdateRequest {
	var signature *sailpointBeta.ConnectorRuleCreateRequestSignature
	if model.Signature != nil {
		inputs := make([]sailpointBeta.Argument, len(model.Signature.Input))
		for i, input := range model.Signature.Input {
			inputs[i] = sailpointBeta.Argument{
				Name:        input.Name.ValueString(),
				Description: input.Description.ValueStringPointer(),
				Type:        *sailpointBeta.NewNullableString(input.Type.ValueStringPointer()),
			}
		}
		var output *sailpointBeta.Argument
		if model.Signature.Output != nil {
			output = &sailpointBeta.Argument{
				Name:        model.Signature.Output.Name.ValueString(),
				Description: model.Signature.Output.Description.ValueStringPointer(),
				Type:        *sailpointBeta.NewNullableString(model.Signature.Output.Type.ValueStringPointer()),
			}
		}
		signature = &sailpointBeta.ConnectorRuleCreateRequestSignature{
			Input:  inputs,
			Output: *sailpointBeta.NewNullableArgument(output),
		}
	}

	return sailpointBeta.ConnectorRuleUpdateRequest{
		Id:          model.Id.ValueString(),
		Name:        model.Name.ValueString(),
		Description: model.Description.ValueStringPointer(),
		Type:        model.Type.ValueString(),
		Signature:   signature,
		SourceCode: sailpointBeta.SourceCode{
			Version: model.SourceCode.Version.ValueString(),
			Script:  model.SourceCode.Script.ValueString(),
		},
		Attributes: util.UnmarshalJsonType(model.Attributes, diagnostics),
	}
}

func (r *connectorRuleResource) mapToTerraformModel(model *connectorRuleModel, resp *sailpointBeta.ConnectorRuleResponse, diagnostics *diag.Diagnostics) {
	model.Id = types.StringValue(resp.Id)
	model.Name = types.StringValue(resp.Name)
	model.Description = types.StringPointerValue(resp.Description)
	model.Type = types.StringValue(resp.Type)
	if resp.Signature != nil {
		inputs := make([]signatureData, len(resp.Signature.Input))
		for i, input := range resp.Signature.Input {
			inputs[i] = signatureData{
				Name:        types.StringValue(input.Name),
				Description: types.StringPointerValue(input.Description),
				Type:        types.StringPointerValue(input.Type.Get()),
			}
		}
		var output *signatureData
		if resp.Signature.Output.IsSet() {
			output = &signatureData{
				Name:        types.StringValue(resp.Signature.Output.Get().Name),
				Description: types.StringPointerValue(resp.Signature.Output.Get().Description),
				Type:        types.StringPointerValue(resp.Signature.Output.Get().Type.Get()),
			}
		}
		model.Signature = &signature{
			Input:  inputs,
			Output: output,
		}
	}
	model.SourceCode = sourceCode{
		Version: types.StringValue(resp.SourceCode.Version),
		Script:  types.StringValue(resp.SourceCode.Script),
	}
	model.Attributes = util.MarshalToJsonTypeWithDefinedSchema(resp.Attributes, model.Attributes, diagnostics)
}
