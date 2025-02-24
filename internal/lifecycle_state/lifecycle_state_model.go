package lifecycle_state

import "github.com/hashicorp/terraform-plugin-framework/types"

type lifecycleStateModel struct {
	Id                      types.String             `tfsdk:"id"`
	IdentityProfileId       types.String             `tfsdk:"identity_profile_id"`
	Name                    types.String             `tfsdk:"name"`
	Enabled                 types.Bool               `tfsdk:"enabled"`
	TechnicalName           types.String             `tfsdk:"technical_name"`
	Description             types.String             `tfsdk:"description"`
	EmailNotificationOption *emailNotificationOption `tfsdk:"email_notification_option"`
	AccountActions          []accountAction          `tfsdk:"account_actions"`
	AccessProfileIds        []types.String           `tfsdk:"access_profile_ids"`
	IdentityState           types.String             `tfsdk:"identity_state"`
}

type emailNotificationOption struct {
	NotifyManagers      types.Bool     `tfsdk:"notify_managers"`
	NotifyAllAdmins     types.Bool     `tfsdk:"notify_all_admins"`
	NotifySpecificUsers types.Bool     `tfsdk:"notify_specific_users"`
	EmailAddressList    []types.String `tfsdk:"email_address_list"`
}

type accountAction struct {
	Action    types.String   `tfsdk:"action"`
	SourceIds []types.String `tfsdk:"source_ids"`
}
