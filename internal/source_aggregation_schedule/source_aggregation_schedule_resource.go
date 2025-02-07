package source_aggregation_schedule

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"net/http"
	"terraform-provider-identitynow/internal/sailpoint/custom"
	"terraform-provider-identitynow/internal/util"
)

var (
	_ resource.Resource              = &sourceAggregationScheduleResource{}
	_ resource.ResourceWithConfigure = &sourceAggregationScheduleResource{}
)

const (
	aggregationTypeAccount     = "account"
	aggregationTypeEntitlement = "entitlement"
)

func NewSourceAggregationScheduleResource() resource.Resource {
	return &sourceAggregationScheduleResource{}
}

type sourceAggregationScheduleResource struct {
	apiClient *custom.APIClient
}

func (r *sourceAggregationScheduleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
	r.apiClient = client
}

func (r *sourceAggregationScheduleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_source_aggregation_schedule"
}

func (r *sourceAggregationScheduleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"source_cloud_id": schema.StringAttribute{
				Description: "Legacy Source ID",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"cron_expression": schema.StringAttribute{
				Description: "Cron Expression for the Schedule",
				Required:    true,
			},
			"aggregation_type": schema.StringAttribute{
				Description: "Aggregation type one of 'account' or 'entitlement'",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf(aggregationTypeAccount, aggregationTypeEntitlement),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *sourceAggregationScheduleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan sourceAggregationScheduleModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	sourceCloudId := plan.SourceCloudId.ValueString()
	cronExpression := plan.CronExpression.ValueString()
	var spResp *http.Response
	var err error
	if plan.AggregationType.ValueString() == aggregationTypeAccount {
		_, spResp, err = r.apiClient.ModifySourceAccountAggregationSchedule(ctx, sourceCloudId, cronExpression)
	} else if plan.AggregationType.ValueString() == aggregationTypeEntitlement {
		_, spResp, err = r.apiClient.ModifySourceEntitlementAggregationSchedule(ctx, sourceCloudId, cronExpression)
	} else {
		resp.Diagnostics.AddError(
			"Invalid Aggregation Type",
			"Aggregation Type must be either 'account' or 'entitlement'",
		)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Source Aggregation Schedule",
			"Could not create Source Aggregation Schedule : "+err.Error()+"\n"+util.GetBody(spResp),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *sourceAggregationScheduleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state sourceAggregationScheduleModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	sourceCloudId := state.SourceCloudId.ValueString()
	aggregationType := state.AggregationType.ValueString()
	var spResp *http.Response
	var err error
	var schedule *custom.SourceAggregationSchedule
	if aggregationType == aggregationTypeAccount {
		schedule, spResp, err = r.apiClient.ReadSourceAccountAggregationSchedule(ctx, sourceCloudId)
	} else if aggregationType == aggregationTypeEntitlement {
		schedule, spResp, err = r.apiClient.ReadSourceEntitlementAggregationSchedule(ctx, sourceCloudId)
	} else {
		resp.Diagnostics.AddError(
			"Invalid Aggregation Type",
			"Aggregation Type must be either 'account' or 'entitlement'",
		)
		return
	}
	if spResp.StatusCode == 404 {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Source Aggregation Schedule",
			"Could not read Source Aggregation Schedule': "+err.Error()+"\n"+util.GetBody(spResp),
		)
		return
	}
	if schedule == nil || len(schedule.CronExpressions) == 0 {
		return
	}

	state.CronExpression = types.StringValue(schedule.CronExpressions[0])
	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *sourceAggregationScheduleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan sourceAggregationScheduleModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	sourceCloudId := plan.SourceCloudId.ValueString()
	cronExpression := plan.CronExpression.ValueString()
	var spResp *http.Response
	var err error
	if plan.AggregationType.ValueString() == aggregationTypeAccount {
		_, spResp, err = r.apiClient.ModifySourceAccountAggregationSchedule(ctx, sourceCloudId, cronExpression)
	} else if plan.AggregationType.ValueString() == aggregationTypeEntitlement {
		_, spResp, err = r.apiClient.ModifySourceEntitlementAggregationSchedule(ctx, sourceCloudId, cronExpression)
	} else {
		resp.Diagnostics.AddError(
			"Invalid Aggregation Type",
			"Aggregation Type must be either 'account' or 'entitlement'",
		)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Source Aggregation Schedule",
			"Could not update Source Aggregation Schedule : "+err.Error()+"\n"+util.GetBody(spResp),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *sourceAggregationScheduleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state sourceAggregationScheduleModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	sourceCloudId := state.SourceCloudId.ValueString()
	var spResp *http.Response
	var err error
	if state.AggregationType.ValueString() == aggregationTypeAccount {
		spResp, err = r.apiClient.DeleteSourceAccountAggregationSchedule(ctx, sourceCloudId)
	} else if state.AggregationType.ValueString() == aggregationTypeEntitlement {
		spResp, err = r.apiClient.DeleteSourceEntitlementAggregationSchedule(ctx, sourceCloudId)
	} else {
		resp.Diagnostics.AddError(
			"Invalid Aggregation Type",
			"Aggregation Type must be either 'account' or 'entitlement'",
		)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Source Aggregation Schedule",
			"Could not delete Source Aggregation Schedule : "+err.Error()+"\n"+util.GetBody(spResp),
		)
		return
	}
}
