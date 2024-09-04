package source_aggregation_schedule

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type sourceAggregationScheduleModel struct {
	SourceCloudId   types.String `tfsdk:"source_cloud_id"`
	CronExpression  types.String `tfsdk:"cron_expression"`
	AggregationType types.String `tfsdk:"aggregation_type"`
}
