package role

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-identitynow/internal/util"
)

type roleModel struct {
	Id                      types.String        `tfsdk:"id"`
	Name                    types.String        `tfsdk:"name"`
	Description             types.String        `tfsdk:"description"`
	Owner                   util.ReferenceModel `tfsdk:"owner"`
	AccessProfiles          types.Set           `tfsdk:"access_profiles"`
	Entitlements            types.Set           `tfsdk:"entitlements"`
	Membership              *roleMembership     `tfsdk:"membership"`
	Enabled                 types.Bool          `tfsdk:"enabled"`
	Requestable             types.Bool          `tfsdk:"requestable"`
	AccessRequestConfig     *requestConfig      `tfsdk:"access_request_config"`
	RevocationRequestConfig *requestConfig      `tfsdk:"revocation_request_config"`
	Segments                types.Set           `tfsdk:"segments"`
}

type roleMembership struct {
	Type       types.String            `tfsdk:"type"`
	Criteria   *roleMembershipCriteria `tfsdk:"criteria"`
	Identities []util.ReferenceModel   `tfsdk:"identities"`
}

type roleMembershipCriteria struct {
	Operation   types.String       `tfsdk:"operation"`
	Key         *roleCriteriaKey   `tfsdk:"key"`
	StringValue types.String       `tfsdk:"string_value"`
	Children    []criteriaChildren `tfsdk:"children"`
}

type roleCriteriaKey struct {
	Type     types.String `tfsdk:"type"`
	Property types.String `tfsdk:"property"`
	SourceId types.String `tfsdk:"source_id"`
}

type criteriaChildren struct {
	Operation   types.String     `tfsdk:"operation"`
	Key         *roleCriteriaKey `tfsdk:"key"`
	StringValue types.String     `tfsdk:"string_value"`
}

type requestConfig struct {
	CommentsRequired       types.Bool        `tfsdk:"comments_required"`
	DenialCommentsRequired types.Bool        `tfsdk:"denial_comments_required"`
	ApprovalSchemas        []approvalSchemas `tfsdk:"approval_schemas"`
}

type approvalSchemas struct {
	ApproverType types.String `tfsdk:"approver_type"`
	ApproverId   types.String `tfsdk:"approver_id"`
}
