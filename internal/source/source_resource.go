package source

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	sailpoint "github.com/sailpoint-oss/golang-sdk/v2"
	sailpoint_v3 "github.com/sailpoint-oss/golang-sdk/v2/api_v3"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"terraform-provider-identitynow/internal/patch"
	"terraform-provider-identitynow/internal/sailpoint/custom"
	"terraform-provider-identitynow/internal/util"
)

// Implementation of IdentityNow Source CRUD - https://developer.sailpoint.com/idn/api/v3/create-source
var (
	_ resource.Resource              = &sourceResource{}
	_ resource.ResourceWithConfigure = &sourceResource{}
)

const FILE_FOLDER = "files"

func NewSourceResource() resource.Resource {
	return &sourceResource{}
}

type sourceResource struct {
	apiClient *sailpoint.APIClient
}

func (r *sourceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *sourceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_source"
}

func (r *sourceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"cloud_external_id": schema.StringAttribute{
				Description: "Legacy Source ID for interacting with CC API",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Human-readable name of the source",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Required: true,
			},
			"owner": util.ResourceReferenceSchema("IDENTITY", true, "Reference to an owning Identity Object"),
			"cluster": schema.ObjectAttribute{
				Description: "Reference to the associated Cluster",
				Optional:    true,
				Computed:    true,
				AttributeTypes: map[string]attr.Type{
					"id":   types.StringType,
					"name": types.StringType,
				},
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
			},
			"account_correlation_config": util.ResourceReferenceSchema("ACCOUNT_CORRELATION_CONFIG", false, "Reference to an Account Correlation Config object"),
			"account_correlation_rule":   util.ResourceReferenceSchema("RULE", false, "Reference to a Rule that can do COMPLEX correlation, should only be used when accountCorrelationConfig can't be used"),
			"manager_correlation_mapping": schema.SingleNestedAttribute{
				Description: "Filter Object used during manager correlation to match incoming manager values to an existing manager's Account/Identity",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"account_attribute_name": schema.StringAttribute{
						Description: "Name of the attribute to use for manager correlation. The value found on the account attribute will be used to lookup the manager's identity",
						Required:    true,
					},
					"identity_attribute_name": schema.StringAttribute{
						Description: "Name of the identity attribute to search when trying to find a manager using the value from the accountAttribute",
						Required:    true,
					},
				},
			},
			"manager_correlation_rule": util.ResourceReferenceSchema("RULE", false, "Reference to the ManagerCorrelationRule, only used when a simple filter isn't sufficient"),
			"before_provisioning_rule": util.ResourceReferenceSchema("RULE", false, "Rule that runs on the CCG and allows for customization of provisioning plans before the connector is called"),
			"password_policies": schema.ListNestedAttribute{
				Description:  "List of references to the associated PasswordPolicy objects",
				Optional:     true,
				NestedObject: util.ResourceReferenceNestedObject("PASSWORD_POLICY"),
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			"features": schema.SetAttribute{
				Description: "Optional features that can be supported by a source.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"type": schema.StringAttribute{
				Description: "Specifies the type of system being managed e.g. Active Directory, Workday, etc.. If you are creating a Delimited File source, you must set the provisionasCsv query parameter to true",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"connector": schema.StringAttribute{
				Description: "Connector script name",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"connector_class": schema.StringAttribute{
				Description: "The fully qualified name of the Java class that implements the connector interface",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"connector_attributes": schema.StringAttribute{
				Description: "Connector specific configuration; will differ from type to type",
				CustomType:  jsontypes.NormalizedType{},
				Optional:    true,
				Computed:    true,
			},
			"connector_attributes_credentials": schema.StringAttribute{
				Description: "Connector specific configuration for storing credentials; will differ from type to type; will be merged with `connector_attributes`",
				CustomType:  jsontypes.ExactType{},
				Optional:    true,
				Sensitive:   true,
			},
			"delete_threshold": schema.Int64Attribute{
				Description: "Number from 0 to 100 that specifies when to skip the delete phase",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(0),
			},
			"authoritative": schema.BoolAttribute{
				Description: "When true indicates the source is referenced by an IdentityProfile.",
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"management_workgroup": util.ResourceReferenceSchema("GOVERNANCE_GROUP", false, "Reference to Management Workgroup for this Source"),
			"status": schema.StringAttribute{
				Description: "A status identifier, giving specific information on why a source is healthy or not",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"connector_id": schema.StringAttribute{
				Description: "The id of connector",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"connector_name": schema.StringAttribute{
				Description: "The name of the connector that was chosen on source creation",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"connection_type": schema.StringAttribute{
				Description: "The type of connection (direct or file)",
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString("direct"),
				Validators: []validator.String{
					stringvalidator.OneOf("direct", "file"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"connector_implementation_id": schema.StringAttribute{
				Description: "The connector implementation id",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"connector_files": schema.SetAttribute{
				Description: "This uploads a supplemental source connector file (like jdbc driver jars) to a source's S3 bucket. Files must be located in the same folder or in folder 'files'.",
				Optional:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func (r *sourceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan sourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	source := r.convertToCreateAPIModel(&plan)
	if resp.Diagnostics.HasError() {
		return
	}
	newModel := r.convertToAPIModel(&plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Creating source '%s': %s", source.Name, util.PrettyPrint(source)))
	sourceResponse, spResp, err := r.apiClient.V3.SourcesAPI.CreateSource(ctx).Source(source).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Source",
			"Could not create Source '"+plan.Name.ValueString()+"': "+err.Error()+"\n"+util.GetBody(spResp),
		)
		return
	}

	source.ConnectorAttributes = make(map[string]interface{})
	jsonPatch := r.generateJsonPatch(&newModel, &source, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		r.apiClient.V3.SourcesAPI.DeleteSource(ctx, *sourceResponse.Id).Execute()
		return
	}
	tflog.Info(ctx, fmt.Sprintf("Modifying source '%s': %s", source.Name, util.PrettyPrint(jsonPatch)))
	sourceResponseAfterPatch, spResp, err := r.apiClient.V3.SourcesAPI.UpdateSource(ctx, *sourceResponse.Id).JsonPatchOperation(jsonPatch).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Source",
			"Could not update Source '"+plan.Name.ValueString()+"': "+err.Error()+"\n"+util.GetBody(spResp),
		)
		r.apiClient.V3.SourcesAPI.DeleteSource(ctx, *sourceResponse.Id).Execute()
		return
	}
	if !plan.ConnectorFiles.IsUnknown() && !plan.ConnectorFiles.IsNull() && len(plan.ConnectorFiles.Elements()) > 0 {
		var sourceFilUploadResponse *sailpoint_v3.Source = nil
		for _, element := range plan.ConnectorFiles.Elements() {
			filePath := element.(basetypes.StringValue)
			sourceFilUploadResponse = r.uploadConnectorFiles(ctx, *sourceResponse.Id, filePath.ValueString(), &resp.Diagnostics)
		}
		if sourceFilUploadResponse != nil {
			r.mapToTerraformModel(&plan, sourceFilUploadResponse, &resp.Diagnostics)
		}
	} else {
		r.mapToTerraformModel(&plan, sourceResponseAfterPatch, &resp.Diagnostics)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *sourceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state sourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	source, spResp, err := r.apiClient.V3.SourcesAPI.GetSource(ctx, state.Id.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Source",
			"Could not read Source '"+state.Name.ValueString()+"': "+err.Error()+"\n"+util.GetBody(spResp),
		)
		return
	}

	r.mapToTerraformModel(&state, source, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *sourceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state sourceModel
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
	jsonPatch := r.generateJsonPatch(&newModel, &oldModel, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, fmt.Sprintf("Updating source '%s' with json patch: %s", state.Id.ValueString(), util.PrettyPrint(jsonPatch)))
	sourceResponse, spResp, err := r.apiClient.V3.SourcesAPI.UpdateSource(ctx, state.Id.ValueString()).JsonPatchOperation(jsonPatch).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Source",
			"Could not update Source '"+plan.Name.ValueString()+"': "+err.Error()+"\n"+util.GetBody(spResp),
		)
		return
	}

	connectorFileDifference := r.getDifference(state.ConnectorFiles, plan.ConnectorFiles)
	if len(connectorFileDifference) > 0 {
		var uploadConnectorFiles *sailpoint_v3.Source
		for _, file := range connectorFileDifference {
			uploadConnectorFiles = r.uploadConnectorFiles(ctx, state.Id.ValueString(), file, &resp.Diagnostics)

		}
		if uploadConnectorFiles != nil {
			r.mapToTerraformModel(&plan, uploadConnectorFiles, &resp.Diagnostics)
		}
	} else {
		r.mapToTerraformModel(&plan, sourceResponse, &resp.Diagnostics)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *sourceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state sourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Deleting source '%s'", state.Id.ValueString()))
	taskResult, spResp, err := r.apiClient.V3.SourcesAPI.DeleteSource(ctx, state.Id.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Source",
			"Could not delete Source '"+state.Name.ValueString()+"': "+err.Error()+"\n"+util.GetBody(spResp),
		)
		return
	}
	err = util.WaitUntilCompletedOrFailAfter(ctx, r.apiClient, *taskResult.Id, 60)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Source",
			"Could not delete Source '"+state.Name.ValueString()+"': "+err.Error(),
		)
		return
	}
}

func (r *sourceResource) convertToCreateAPIModel(model *sourceModel) sailpoint_v3.Source {
	cluster := sailpoint_v3.NullableSourceCluster{}
	if !model.Cluster.IsNull() && !model.Cluster.IsUnknown() {
		attributes := model.Cluster.Attributes()
		cluster = *sailpoint_v3.NewNullableSourceCluster(&sailpoint_v3.SourceCluster{
			Type: "CLUSTER",
			Id:   attributes["id"].(basetypes.StringValue).ValueString(),
			Name: attributes["name"].(basetypes.StringValue).ValueString(),
		})
	}
	return sailpoint_v3.Source{
		Name:        model.Name.ValueString(),
		Description: util.GetTFStringPointer(model.Description),
		Connector:   model.Connector.ValueString(),
		Owner: sailpoint_v3.SourceOwner{
			Type: util.GetTFStringPointer(model.Owner.Type),
			Id:   util.GetTFStringPointer(model.Owner.Id),
			Name: util.GetTFStringPointer(model.Owner.Name),
		},
		Cluster: cluster,
	}
}

func (r *sourceResource) convertToAPIModel(model *sourceModel, diagnostics *diag.Diagnostics) sailpoint_v3.Source {
	connectorAttributes := r.getConnectorAttributes(model, diagnostics)
	if diagnostics.HasError() {
		return sailpoint_v3.Source{}
	}
	var passwordPolicies []sailpoint_v3.SourcePasswordPoliciesInner = nil
	if len(model.PasswordPolicies) > 0 {
		passwordPolicies = make([]sailpoint_v3.SourcePasswordPoliciesInner, len(model.PasswordPolicies))
		for index, item := range model.PasswordPolicies {
			passwordPolicies[index] = sailpoint_v3.SourcePasswordPoliciesInner{
				Type: util.GetTFStringPointer(item.Type),
				Id:   util.GetTFStringPointer(item.Id),
				Name: util.GetTFStringPointer(item.Name),
			}
		}
	}
	cluster := sailpoint_v3.NullableSourceCluster{}
	if !model.Cluster.IsNull() && !model.Cluster.IsUnknown() {
		attributes := model.Cluster.Attributes()
		cluster = *sailpoint_v3.NewNullableSourceCluster(&sailpoint_v3.SourceCluster{
			Type: "CLUSTER",
			Id:   attributes["id"].(basetypes.StringValue).ValueString(),
			Name: attributes["name"].(basetypes.StringValue).ValueString(),
		})
	}
	accountCorrelationConfig := sailpoint_v3.NullableSourceAccountCorrelationConfig{}
	if model.AccountCorrelationConfig != nil {
		accountCorrelationConfig = *sailpoint_v3.NewNullableSourceAccountCorrelationConfig(
			&sailpoint_v3.SourceAccountCorrelationConfig{
				Type: util.GetTFStringPointer(model.AccountCorrelationConfig.Type),
				Id:   util.GetTFStringPointer(model.AccountCorrelationConfig.Id),
				Name: util.GetTFStringPointer(model.AccountCorrelationConfig.Name),
			})
	}
	accountCorrelationRule := sailpoint_v3.NullableSourceAccountCorrelationRule{}
	if model.AccountCorrelationRule != nil {
		accountCorrelationRule = *sailpoint_v3.NewNullableSourceAccountCorrelationRule(
			&sailpoint_v3.SourceAccountCorrelationRule{
				Type: util.GetTFStringPointer(model.AccountCorrelationRule.Type),
				Id:   util.GetTFStringPointer(model.AccountCorrelationRule.Id),
				Name: util.GetTFStringPointer(model.AccountCorrelationRule.Name),
			})
	}
	var managerCorrelationMapping *sailpoint_v3.SourceManagerCorrelationMapping
	if model.ManagerCorrelationMapping != nil {
		managerCorrelationMapping = &sailpoint_v3.SourceManagerCorrelationMapping{
			AccountAttributeName:  util.GetTFStringPointer(model.ManagerCorrelationMapping.AccountAttributeName),
			IdentityAttributeName: util.GetTFStringPointer(model.ManagerCorrelationMapping.IdentityAttributeName),
		}
	}
	managerCorrelationRule := sailpoint_v3.NullableSourceManagerCorrelationRule{}
	if model.ManagerCorrelationRule != nil {
		managerCorrelationRule = *sailpoint_v3.NewNullableSourceManagerCorrelationRule(
			&sailpoint_v3.SourceManagerCorrelationRule{
				Type: util.GetTFStringPointer(model.ManagerCorrelationRule.Type),
				Id:   util.GetTFStringPointer(model.ManagerCorrelationRule.Id),
				Name: util.GetTFStringPointer(model.ManagerCorrelationRule.Name),
			})
	}
	deleteThreshold := int32(model.DeleteThreshold.ValueInt64())
	beforeProvisioningRule := sailpoint_v3.NullableSourceBeforeProvisioningRule{}
	if model.BeforeProvisioningRule != nil {
		beforeProvisioningRule = *sailpoint_v3.NewNullableSourceBeforeProvisioningRule(
			&sailpoint_v3.SourceBeforeProvisioningRule{
				Type: util.GetTFStringPointer(model.BeforeProvisioningRule.Type),
				Id:   util.GetTFStringPointer(model.BeforeProvisioningRule.Id),
				Name: util.GetTFStringPointer(model.BeforeProvisioningRule.Name),
			})
	}
	managementWorkgroup := sailpoint_v3.NullableSourceManagementWorkgroup{}
	if model.ManagementWorkgroup != nil {
		managementWorkgroup = *sailpoint_v3.NewNullableSourceManagementWorkgroup(
			&sailpoint_v3.SourceManagementWorkgroup{
				Type: util.GetTFStringPointer(model.ManagementWorkgroup.Type),
				Id:   util.GetTFStringPointer(model.ManagementWorkgroup.Id),
				Name: util.GetTFStringPointer(model.ManagementWorkgroup.Name),
			})
	}
	features := make([]string, len(model.Features.Elements()))
	if !model.Features.IsNull() && !model.Features.IsUnknown() {
		for index, item := range model.Features.Elements() {
			stringItem := item.(basetypes.StringValue)
			features[index] = stringItem.ValueString()
		}
	}
	return sailpoint_v3.Source{
		Name:        model.Name.ValueString(),
		Description: util.GetTFStringPointer(model.Description),
		Owner: sailpoint_v3.SourceOwner{
			Type: util.GetTFStringPointer(model.Owner.Type),
			Id:   util.GetTFStringPointer(model.Owner.Id),
			Name: util.GetTFStringPointer(model.Owner.Name),
		},
		Cluster:                   cluster,
		AccountCorrelationConfig:  accountCorrelationConfig,
		AccountCorrelationRule:    accountCorrelationRule,
		ManagerCorrelationMapping: managerCorrelationMapping,
		ManagerCorrelationRule:    managerCorrelationRule,
		BeforeProvisioningRule:    beforeProvisioningRule,
		PasswordPolicies:          passwordPolicies,
		Features:                  features,
		Type:                      util.GetTFStringPointer(model.Type),
		Connector:                 model.Connector.ValueString(),
		ConnectorClass:            util.GetTFStringPointer(model.ConnectorClass),
		ConnectorAttributes:       connectorAttributes,
		DeleteThreshold:           &deleteThreshold,
		Authoritative:             model.Authoritative.ValueBoolPointer(),
		ManagementWorkgroup:       managementWorkgroup,
		Status:                    util.GetTFStringPointer(model.Status),
		ConnectorId:               util.GetTFStringPointer(model.ConnectorId),
		ConnectorName:             util.GetTFStringPointer(model.ConnectorName),
		ConnectionType:            util.GetTFStringPointer(model.ConnectionType),
		ConnectorImplementationId: util.GetTFStringPointer(model.ConnectorImplementationId),
	}
}

func (r *sourceResource) getConnectorAttributes(model *sourceModel, diagnostics *diag.Diagnostics) map[string]interface{} {
	if model.ConnectorAttributes.IsNull() || model.ConnectorAttributes.IsUnknown() {
		return nil
	}
	connectorAttributes := util.UnmarshalJsonTypeNormalized(model.ConnectorAttributes, diagnostics)
	if diagnostics.HasError() {
		return nil
	}
	connectorAttributesCred := util.UnmarshalJsonType(model.ConnectorAttributesCredentials, diagnostics)
	if diagnostics.HasError() {
		return nil
	}
	return r.mergeMaps(connectorAttributesCred, connectorAttributes)
}

func (r *sourceResource) mergeMaps(connectorAttributesCred map[string]interface{}, connectorAttributes map[string]interface{}) map[string]interface{} {
	for key, value := range connectorAttributesCred {
		if _, ok := connectorAttributes[key]; ok {
			if reflect.TypeOf(value).Kind() == reflect.Map {
				connectorAttributes[key] = r.mergeMaps(value.(map[string]interface{}), connectorAttributes[key].(map[string]interface{}))
			} else if reflect.TypeOf(value).Kind() == reflect.Slice {
				connectorAttributes[key] = r.mergeSlices(connectorAttributes[key].([]interface{}), value.([]interface{}))
			} else {
				connectorAttributes[key] = value
			}
		} else {
			connectorAttributes[key] = value
		}
	}
	return connectorAttributes
}

func (r *sourceResource) mergeSlices(leftSlice, rightSlice []interface{}) []interface{} {
	for i, value := range rightSlice {
		if len(leftSlice) <= i {
			leftSlice = append(leftSlice, value)
		} else {
			if reflect.TypeOf(value).Kind() == reflect.Map {
				leftSlice[i] = r.mergeMaps(value.(map[string]interface{}), leftSlice[i].(map[string]interface{}))
			} else {
				leftSlice = append(leftSlice, value)
			}
		}
	}
	return leftSlice
}

func (r *sourceResource) mapToTerraformModel(tfModel *sourceModel, source *sailpoint_v3.Source, diagnostics *diag.Diagnostics) {
	tfModel.Id = types.StringPointerValue(source.Id)
	if val, ok := source.ConnectorAttributes["cloudExternalId"]; ok {
		tfModel.CloudExternalId = types.StringValue(val.(string))
	}
	if val, ok := source.ConnectorAttributes["connector_files"]; ok {
		if strVal, ok := val.(string); ok {
			files := strings.Split(strVal, ",")
			connectorFiles := make([]attr.Value, len(files))
			for i, file := range files {
				connectorFiles[i] = types.StringValue(file)
			}
			tfModel.ConnectorFiles = types.SetValueMust(types.StringType, connectorFiles)
		} else {
			tfModel.ConnectorFiles = types.SetNull(types.StringType)
		}
	} else {
		tfModel.ConnectorFiles = types.SetNull(types.StringType)
	}
	tfModel.Name = types.StringValue(source.Name)
	tfModel.Description = types.StringPointerValue(source.Description)
	tfModel.Owner = util.NewPointerReferenceModel(source.Owner.Type, source.Owner.Id, source.Owner.Name)
	if source.Cluster.Get() != nil {
		tfModel.Cluster = basetypes.NewObjectValueMust(map[string]attr.Type{
			"id":   types.StringType,
			"name": types.StringType,
		}, map[string]attr.Value{
			"id":   types.StringValue(source.Cluster.Get().Id),
			"name": types.StringValue(source.Cluster.Get().Name),
		})
	} else {
		tfModel.Cluster = basetypes.NewObjectNull(map[string]attr.Type{
			"id":   types.StringType,
			"name": types.StringType,
		})
	}
	if source.AccountCorrelationConfig.Get() != nil && tfModel.AccountCorrelationConfig != nil {
		config := source.AccountCorrelationConfig.Get()
		tfModel.AccountCorrelationConfig = util.NewPointerReferenceModel(config.Type, config.Id, config.Name)
	}
	if source.AccountCorrelationRule.Get() != nil && tfModel.AccountCorrelationRule != nil {
		rule := source.AccountCorrelationRule.Get()
		tfModel.AccountCorrelationRule = util.NewPointerReferenceModel(rule.Type, rule.Id, rule.Name)
	}
	if source.ManagerCorrelationMapping != nil {
		tfModel.ManagerCorrelationMapping = &managerCorrelationModel{
			AccountAttributeName:  types.StringPointerValue(source.ManagerCorrelationMapping.AccountAttributeName),
			IdentityAttributeName: types.StringPointerValue(source.ManagerCorrelationMapping.IdentityAttributeName),
		}
	}
	if source.ManagerCorrelationRule.Get() != nil && tfModel.ManagerCorrelationRule != nil {
		rule := source.ManagerCorrelationRule.Get()
		tfModel.ManagerCorrelationRule = util.NewPointerReferenceModel(rule.Type, rule.Id, rule.Name)
	}
	if source.BeforeProvisioningRule.Get() != nil && tfModel.BeforeProvisioningRule != nil {
		rule := source.BeforeProvisioningRule.Get()
		tfModel.BeforeProvisioningRule = util.NewPointerReferenceModel(rule.Type, rule.Id, rule.Name)
	}
	if source.PasswordPolicies != nil {
		for _, item := range source.PasswordPolicies {
			tfModel.PasswordPolicies = append(tfModel.PasswordPolicies, *util.NewPointerReferenceModel(item.Type, item.Id, item.Name))
		}
	}
	features := make([]attr.Value, len(source.Features))
	for i, item := range source.Features {
		features[i] = types.StringValue(item)
	}
	tfModel.Features = types.SetValueMust(types.StringType, features)
	tfModel.Type = types.StringPointerValue(source.Type)
	tfModel.Connector = types.StringValue(source.Connector)
	tfModel.ConnectorClass = types.StringPointerValue(source.ConnectorClass)
	if source.ConnectorAttributes != nil {
		tfModel.ConnectorAttributes = util.MarshalToJsonTypeWithDefinedSchemaNormalized(source.ConnectorAttributes, tfModel.ConnectorAttributes, diagnostics)
	}
	tfModel.DeleteThreshold = types.Int64Value(int64(*source.DeleteThreshold))
	tfModel.Authoritative = types.BoolPointerValue(source.Authoritative)
	if source.ManagementWorkgroup.Get() != nil {
		workgroup := source.ManagementWorkgroup.Get()
		tfModel.ManagementWorkgroup = util.NewPointerReferenceModel(workgroup.Type, workgroup.Id, workgroup.Name)
	}
	tfModel.Status = types.StringPointerValue(source.Status)
	tfModel.ConnectorId = types.StringPointerValue(source.ConnectorId)
	tfModel.ConnectorName = types.StringPointerValue(source.ConnectorName)
	tfModel.ConnectionType = types.StringPointerValue(source.ConnectionType)
	tfModel.ConnectorImplementationId = types.StringPointerValue(source.ConnectorImplementationId)

}

func (r *sourceResource) generateJsonPatch(newModel *sailpoint_v3.Source, oldModel *sailpoint_v3.Source, diagnostics *diag.Diagnostics) []sailpoint_v3.JsonPatchOperation {
	jsonPatch, err := patch.NewSourcePatchBuilder(newModel, oldModel).GenerateJsonPatch()
	if err != nil {
		diagnostics.AddError(
			"Error Generating Update Patch",
			"Could not generate update patch for Source '"+oldModel.Name+"': "+err.Error(),
		)
		return nil
	}
	v3JsonPatch, err := patch.ConvertFromBetaToV3(jsonPatch)
	if err != nil {
		diagnostics.AddError(
			"Error Generating Update Patch",
			"Could not convert patch to V3 for Source '"+oldModel.Name+"': "+err.Error(),
		)
		return nil
	}
	return v3JsonPatch
}

func (r *sourceResource) uploadConnectorFiles(ctx context.Context, sourceId string, filePath string, diagnostics *diag.Diagnostics) *sailpoint_v3.Source {
	tflog.Info(ctx, "Uploading connector files "+filePath)
	if filePath == "" {
		return nil
	}
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		filePath = filepath.FromSlash(FILE_FOLDER + "/" + filePath)
	}
	file, err := os.Open(filePath)
	if errors.Is(err, os.ErrNotExist) {
		diagnostics.AddError(
			"Error Uploading Source Connector Files",
			"File '"+filePath+"' does not exist: "+err.Error(),
		)
		return nil
	}
	source, spResp, err := r.apiClient.V3.SourcesAPI.ImportConnectorFile(ctx, sourceId).File(file).Execute()
	if err != nil {
		diagnostics.AddError(
			"Error Uploading Connector Files",
			"Could not upload connector files for Source '"+sourceId+"': "+err.Error()+"\n"+util.GetBody(spResp),
		)
		return nil
	}
	return source
}

func (r *sourceResource) getDifference(current, desired types.Set) []string {
	currentArray := getStringArray(current)
	desiredArray := getStringArray(desired)
	return getDifference(desiredArray, currentArray)
}

func getStringArray(set types.Set) []string {
	var array []string
	for _, item := range set.Elements() {
		array = append(array, item.(basetypes.StringValue).ValueString())
	}
	return array
}

func getDifference(a, b []string) (diff []string) {
	m := make(map[string]bool)

	for _, item := range b {
		m[item] = true
	}

	for _, item := range a {
		if _, ok := m[item]; !ok {
			diff = append(diff, item)
		}
	}
	return
}
