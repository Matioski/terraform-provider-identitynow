package source

import (
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-identitynow/internal/util"
)

type sourceModel struct {
	Id                             types.String             `tfsdk:"id"`
	CloudExternalId                types.String             `tfsdk:"cloud_external_id"`
	Name                           types.String             `tfsdk:"name"`
	Description                    types.String             `tfsdk:"description"`
	Owner                          *util.ReferenceModel     `tfsdk:"owner"`
	Cluster                        types.Object             `tfsdk:"cluster"`
	AccountCorrelationConfig       *util.ReferenceModel     `tfsdk:"account_correlation_config"`
	AccountCorrelationRule         *util.ReferenceModel     `tfsdk:"account_correlation_rule"`
	ManagerCorrelationMapping      *managerCorrelationModel `tfsdk:"manager_correlation_mapping"`
	ManagerCorrelationRule         *util.ReferenceModel     `tfsdk:"manager_correlation_rule"`
	BeforeProvisioningRule         *util.ReferenceModel     `tfsdk:"before_provisioning_rule"`
	PasswordPolicies               []util.ReferenceModel    `tfsdk:"password_policies"`
	Features                       types.Set                `tfsdk:"features"`
	Type                           types.String             `tfsdk:"type"`
	Connector                      types.String             `tfsdk:"connector"`
	ConnectorClass                 types.String             `tfsdk:"connector_class"`
	ConnectorAttributes            jsontypes.Normalized     `tfsdk:"connector_attributes"`
	ConnectorAttributesCredentials jsontypes.Exact          `tfsdk:"connector_attributes_credentials"`
	DeleteThreshold                types.Int64              `tfsdk:"delete_threshold"`
	Authoritative                  types.Bool               `tfsdk:"authoritative"`
	ManagementWorkgroup            *util.ReferenceModel     `tfsdk:"management_workgroup"`
	Status                         types.String             `tfsdk:"status"`
	ConnectorId                    types.String             `tfsdk:"connector_id"`
	ConnectorName                  types.String             `tfsdk:"connector_name"`
	ConnectionType                 types.String             `tfsdk:"connection_type"`
	ConnectorImplementationId      types.String             `tfsdk:"connector_implementation_id"`
	ConnectorFiles                 types.Set                `tfsdk:"connector_files"`
}

type managerCorrelationModel struct {
	AccountAttributeName  types.String `tfsdk:"account_attribute_name"`
	IdentityAttributeName types.String `tfsdk:"identity_attribute_name"`
}
