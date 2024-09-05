package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"terraform-provider-identitynow/internal/cluster"
	"terraform-provider-identitynow/internal/connector"
	"terraform-provider-identitynow/internal/connector_rule"
	"terraform-provider-identitynow/internal/entitlement"
	"terraform-provider-identitynow/internal/identity"
	"terraform-provider-identitynow/internal/identity_attribute"
	"terraform-provider-identitynow/internal/identity_profile"
	"terraform-provider-identitynow/internal/lifecycle_state"
	"terraform-provider-identitynow/internal/role"
	"terraform-provider-identitynow/internal/sailpoint/custom"
	"terraform-provider-identitynow/internal/source"
	"terraform-provider-identitynow/internal/source_schema"
	"terraform-provider-identitynow/internal/transform"
	"terraform-provider-identitynow/internal/workflow"

	sailpoint "github.com/sailpoint-oss/golang-sdk/v2"
	"os"
)

// Ensure identityNowProvider satisfies various provider interfaces.
var (
	_ provider.Provider = &identityNowProvider{}
)

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &identityNowProvider{
			version: version,
		}
	}
}

// identityNowProvider defines the provider implementation.
type identityNowProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

type identityNowProviderModel struct {
	Host         types.String `tfsdk:"host"`
	ClientId     types.String `tfsdk:"client_id"`
	ClientSecret types.String `tfsdk:"client_secret"`
}

// ScaffoldingProviderModel describes the provider data model.
type ScaffoldingProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
}

func (p *identityNowProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "identitynow"
	resp.Version = p.version
}

func (p *identityNowProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Description: "URI for IdentityNow API Tenant. May also be provided via IDN_HOST environment variable.",
				Optional:    true,
			},
			"client_id": schema.StringAttribute{
				Description: "Client ID for authentication with IdentityNow API Tenant. May also be provided via IDN_CLIENT_ID environment variable.",
				Optional:    true,
			},
			"client_secret": schema.StringAttribute{
				Description: "Client Secret for authentication with IdentityNow API Tenant. May also be provided via IDN_CLIENT_SECRET environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

func (p *identityNowProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring IdentityNow API Client")

	var config identityNowProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown IdentityNow API Host",
			"The provider cannot create the IdentityNow API client as there is an unknown configuration value for the IdentityNow API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the IDN_HOST environment variable.",
		)
	}

	if config.ClientSecret.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("clientId"),
			"Unknown IdentityNow API client_id",
			"The provider cannot create the IdentityNow API client as there is an unknown configuration value for the IdentityNow API client_id. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the IDN_CLIENT_ID environment variable.",
		)
	}

	if config.ClientSecret.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("client_secret"),
			"Unknown IdentityNow API ClientSecret",
			"The provider cannot create the IdentityNow API client as there is an unknown configuration value for the IdentityNow API client_secret. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the IDN_CLIENT_SECRET environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.
	host := os.Getenv("IDN_HOST")
	clientId := os.Getenv("IDN_CLIENT_ID")
	clientSecret := os.Getenv("IDN_CLIENT_SECRET")

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if !config.ClientId.IsNull() {
		clientId = config.ClientId.ValueString()
	}

	if !config.ClientSecret.IsNull() {
		clientSecret = config.ClientSecret.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing IdentityNow API Host",
			"The provider cannot create the IdentityNow API client as there is a missing or empty value for the IdentityNow API host. "+
				"Set the host value in the configuration or use the IDN_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if clientId == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("client_id"),
			"Missing IdentityNow API ClientId",
			"The provider cannot create the IdentityNow API client as there is a missing or empty value for the IdentityNow API client_id. "+
				"Set the username value in the configuration or use the IDN_CLIENT_ID environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if clientSecret == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("client_secret"),
			"Missing IdentityNow API ClientSecret",
			"The provider cannot create the IdentityNow API client as there is a missing or empty value for the IdentityNow API client_secret. "+
				"Set the password value in the configuration or use the IDN_CLIENT_SECRET environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	configuration := sailpoint.NewConfiguration(sailpoint.ClientConfiguration{
		ClientId:     clientId,
		ClientSecret: clientSecret,
		BaseURL:      host,
		TokenURL:     host + "/oauth/token",
	})

	apiClient := sailpoint.NewAPIClient(configuration)
	configuration.HTTPClient.RetryMax = 5
	client := custom.NewAPIClient(apiClient, configuration)
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *identityNowProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		identity_attribute.NewIdentityAttributeResource,
		transform.NewTransformResource,
		source.NewSourceResource,
		identity_profile.NewIdentityProfileResource,
		source_schema.NewSourceSchemaResource,
		lifecycle_state.NewLifecycleStateResource,
		connector_rule.NewConnectorRuleResource,
		workflow.NewWorkflowResource,
		role.NewRoleResource,
	}
}

func (p *identityNowProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		identity.NewIdentityDataSource,
		cluster.NewClusterDataSource,
		connector.NewConnectorDataSource,
		entitlement.NewEntitlementDataSource,
	}
}
