package org_config

import "github.com/hashicorp/terraform-plugin-framework/types"

type orgConfigModel struct {
	TimeZone types.String `tfsdk:"time_zone"`
}
