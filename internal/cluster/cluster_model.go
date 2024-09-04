package cluster

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type clusterModel struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	Pod  types.String `tfsdk:"pod"`
	Org  types.String `tfsdk:"org"`
}
