package source_schema

import (
	"context"
	"fmt"
	"terraform-provider-identitynow/internal/sailpoint/custom"
	"terraform-provider-identitynow/internal/util"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	sailpoint "github.com/sailpoint-oss/golang-sdk/v2"
	sailpointV3 "github.com/sailpoint-oss/golang-sdk/v2/api_v3"
)

var (
	_ resource.Resource              = &sourceSchemaResource{}
	_ resource.ResourceWithConfigure = &sourceSchemaResource{}
)

func NewSourceSchemaResource() resource.Resource {
	return &sourceSchemaResource{}
}

type sourceSchemaResource struct {
	apiClient *sailpoint.APIClient
}

func (r *sourceSchemaResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *sourceSchemaResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_source_schema"
}

func (r *sourceSchemaResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The id of the Schema",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"source_id": schema.StringAttribute{
				Description: "The Source id",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the Schema",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"native_object_type": schema.StringAttribute{
				Description: "The name of the object type on the native system that the schema represents",
				Optional:    true,
			},
			"identity_attribute": schema.StringAttribute{
				Description: "The name of the attribute used to calculate the unique identifier for an object in the schema",
				Optional:    true,
			},
			"display_attribute": schema.StringAttribute{
				Description: "The name of the attribute used to calculate the display value for an object in the schema",
				Optional:    true,
			},
			"hierarchy_attribute": schema.StringAttribute{
				Description: "The name of the attribute whose values represent other objects in a hierarchy. Only relevant to group schemas",
				Optional:    true,
			},
			"include_permissions": schema.BoolAttribute{
				Description: "Flag indicating whether or not the include permissions with the object data when aggregating the schema",
				Optional:    true,
			},
			"features": schema.SetAttribute{
				Description: "The features that the schema supports",
				Required:    true,
				ElementType: types.StringType,
			},
			"configuration": schema.StringAttribute{
				Description: "Holds any extra configuration data that the schema may require",
				CustomType:  jsontypes.ExactType{},
				Required:    true,
			},
			"attributes": schema.ListNestedAttribute{
				Description: "The attribute definitions which form the schema",
				Required:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "The name of the attribute",
							Required:    true,
						},
						"type": schema.StringAttribute{
							Description: "The type of the attribute. One of 'STRING', 'LONG', 'INT', 'BOOLEAN'",
							Required:    true,
						},
						"schema": util.ResourceReferenceSchema("CONNECTOR_SCHEMA", false, "A reference to the schema on the source to the attribute values map to"),
						"description": schema.StringAttribute{
							Description: "A human-readable description of the attribute",
							Optional:    true,
						},
						"is_multi": schema.BoolAttribute{
							Description: "Flag indicating whether or not the attribute is multi-valued",
							Required:    true,
						},
						"is_entitlement": schema.BoolAttribute{
							Description: "Flag indicating whether or not the attribute is an entitlement",
							Required:    true,
						},
						"is_group": schema.BoolAttribute{
							Description: "Flag indicating whether or not the attribute represents a group. This can only be true if isEntitlement is also true and there is a schema defined for the attribute",
							Required:    true,
						},
					},
				},
			},
		},
	}
}

func (r *sourceSchemaResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan sourceSchemaModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	sourceId := plan.SourceId.ValueString()
	schemaName := plan.Name.ValueString()
	schema := r.convertToAPIModel(&plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	existingSchema := r.findSchema(ctx, sourceId, schemaName, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	if existingSchema != nil {
		spResp, err := r.apiClient.V3.SourcesAPI.DeleteSourceSchema(ctx, sourceId, *existingSchema.Id).Execute()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Creating Source Schema",
				"Could not delete existing Source Schema '"+*existingSchema.Id+"': "+err.Error()+"\n"+util.GetBody(spResp),
			)
			return
		}
	}
	tflog.Info(ctx, fmt.Sprintf("Creating New Schema: %s", util.PrettyPrint(existingSchema)))
	schemaResp, spResp, err := r.apiClient.V3.SourcesAPI.CreateSourceSchema(ctx, sourceId).Schema(schema).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Source Schema",
			"Could not create Source Schema '"+schemaName+"': "+err.Error()+"\n"+util.GetBody(spResp),
		)
		return
	}

	r.mapToTerraformModel(&plan, schemaResp, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *sourceSchemaResource) findSchema(ctx context.Context, sourceId, schemaName string, diagnostics *diag.Diagnostics) *sailpointV3.Schema {
	schemas, spResp, err := r.apiClient.V3.SourcesAPI.GetSourceSchemas(ctx, sourceId).Execute()
	if err != nil {
		diagnostics.AddError(
			"Error Creating Source Schema",
			"Error during source schema lookup '"+schemaName+"': "+err.Error()+"\n"+util.GetBody(spResp),
		)
		return nil
	}
	for _, schema := range schemas {
		if schema.Name != nil && *schema.Name == schemaName {
			return &schema
		}
	}
	return nil
}

func (r *sourceSchemaResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state sourceSchemaModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	sourceId := state.SourceId.ValueString()
	schemaId := state.Id.ValueString()
	schema, spResp, err := r.apiClient.V3.SourcesAPI.GetSourceSchema(ctx, sourceId, schemaId).Execute()
	if spResp.StatusCode == 404 {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Source Schema",
			"Could not read Source Schema '"+state.Name.ValueString()+"': "+err.Error()+"\n"+util.GetBody(spResp),
		)
		return
	}
	tflog.Info(ctx, fmt.Sprintf("Source Schema: %s", util.PrettyPrint(schema)))

	r.mapToTerraformModel(&state, schema, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *sourceSchemaResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan sourceSchemaModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	schema := r.convertToAPIModel(&plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, fmt.Sprintf("Updating Schema: %s", util.PrettyPrint(schema)))
	schemaResp, spResp, err := r.apiClient.V3.SourcesAPI.PutSourceSchema(ctx, plan.SourceId.ValueString(), plan.Id.ValueString()).Schema(schema).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Source Schema",
			"Could not update Source Schema '"+plan.Name.ValueString()+"': "+err.Error()+"\n"+util.GetBody(spResp),
		)
		return
	}

	r.mapToTerraformModel(&plan, schemaResp, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *sourceSchemaResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state sourceSchemaModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	spResp, err := r.apiClient.V3.SourcesAPI.DeleteSourceSchema(ctx, state.SourceId.ValueString(), state.Id.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Source Schema",
			"Could not delete Source Schema '"+state.Name.ValueString()+"': "+err.Error()+"\n"+util.GetBody(spResp),
		)
		return
	}
}

func (r *sourceSchemaResource) convertToAPIModel(model *sourceSchemaModel, diagnostics *diag.Diagnostics) sailpointV3.Schema {
	features := make([]string, len(model.Features))
	for index, item := range model.Features {
		features[index] = item.ValueString()
	}
	attributes := make([]sailpointV3.AttributeDefinition, len(model.Attributes))
	for i, attribute := range model.Attributes {
		attrType, err := sailpointV3.NewAttributeDefinitionTypeFromValue(attribute.Type.ValueString())
		if err != nil {
			diagnostics.AddError(
				"Error Processing Schema Attribute Type",
				fmt.Sprintf("Could not process Schema Attribute Type '%s': %s", attribute.Type.ValueString(), err.Error()),
			)
			return sailpointV3.Schema{}
		}
		var attrSchema *sailpointV3.AttributeDefinitionSchema
		if attribute.Schema != nil {
			attrSchema = &sailpointV3.AttributeDefinitionSchema{
				Type: attribute.Schema.Type.ValueStringPointer(),
				Id:   attribute.Schema.Id.ValueStringPointer(),
				Name: attribute.Schema.Name.ValueStringPointer(),
			}
		}
		attributes[i] = sailpointV3.AttributeDefinition{
			Name:          attribute.Name.ValueStringPointer(),
			Type:          attrType,
			Description:   attribute.Description.ValueStringPointer(),
			IsMulti:       attribute.IsMulti.ValueBoolPointer(),
			IsEntitlement: attribute.IsEntitlement.ValueBoolPointer(),
			IsGroup:       attribute.IsGroup.ValueBoolPointer(),
			Schema:        attrSchema,
		}
	}
	return sailpointV3.Schema{
		Id:                 model.Id.ValueStringPointer(),
		Name:               model.Name.ValueStringPointer(),
		NativeObjectType:   model.NativeObjectType.ValueStringPointer(),
		IdentityAttribute:  model.IdentityAttribute.ValueStringPointer(),
		DisplayAttribute:   model.DisplayAttribute.ValueStringPointer(),
		HierarchyAttribute: model.HierarchyAttribute.ValueStringPointer(),
		IncludePermissions: model.IncludePermissions.ValueBoolPointer(),
		Features:           features,
		Configuration:      util.UnmarshalJsonType(model.Configuration, diagnostics),
		Attributes:         attributes,
	}
}

func (r *sourceSchemaResource) mapToTerraformModel(tfModel *sourceSchemaModel, schema *sailpointV3.Schema, diagnostics *diag.Diagnostics) {
	tfModel.Id = types.StringPointerValue(schema.Id)
	tfModel.Name = types.StringPointerValue(schema.Name)
	tfModel.NativeObjectType = types.StringPointerValue(schema.NativeObjectType)
	tfModel.IdentityAttribute = types.StringPointerValue(schema.IdentityAttribute)
	tfModel.DisplayAttribute = types.StringPointerValue(schema.DisplayAttribute)
	tfModel.HierarchyAttribute = types.StringPointerValue(schema.HierarchyAttribute)
	tfModel.IncludePermissions = types.BoolPointerValue(schema.IncludePermissions)
	tfModel.Features = make([]types.String, len(schema.Features))
	for i, item := range schema.Features {
		tfModel.Features[i] = types.StringValue(item)
	}
	tfModel.Configuration = util.MarshalToJsonType(schema.Configuration, diagnostics)
	attributes := make([]attributeModel, len(schema.Attributes))
	for i, item := range schema.Attributes {
		var attrType string
		if item.Type != nil {
			attrType = string(*item.Type)
		}
		var attrSchema *util.ReferenceModel
		if item.Schema != nil {
			attrSchema = util.NewPointerReferenceModel(item.Schema.Type, item.Schema.Id, item.Schema.Name)
		}
		attributes[i] = attributeModel{
			Name:          types.StringPointerValue(item.Name),
			Description:   types.StringPointerValue(item.Description),
			Type:          types.StringValue(attrType),
			Schema:        attrSchema,
			IsMulti:       types.BoolPointerValue(item.IsMulti),
			IsEntitlement: types.BoolPointerValue(item.IsEntitlement),
			IsGroup:       types.BoolPointerValue(item.IsGroup),
		}
	}
	tfModel.Attributes = attributes
}
