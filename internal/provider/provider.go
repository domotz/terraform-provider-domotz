package provider

import (
	"context"
	"os"

	"github.com/domotz/terraform-provider-domotz/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	defaultBaseURL = "https://api-eu-west-1-cell-1.domotz.com/public-api/v1"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ provider.Provider = &DomotzProvider{}
)

// DomotzProvider defines the provider implementation
type DomotzProvider struct {
	version string
}

// DomotzProviderModel describes the provider data model
type DomotzProviderModel struct {
	APIKey  types.String `tfsdk:"api_key"`
	BaseURL types.String `tfsdk:"base_url"`
}

// Metadata returns the provider type name
func (p *DomotzProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "domotz"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data
func (p *DomotzProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Terraform provider for Domotz network monitoring platform.",
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				Description: "Domotz API key for authentication. Can also be set via DOMOTZ_API_KEY environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
			"base_url": schema.StringAttribute{
				Description: "Base URL for the Domotz API. Defaults to production endpoint.",
				Optional:    true,
			},
		},
	}
}

// Configure prepares a Domotz API client for data sources and resources
func (p *DomotzProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config DomotzProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read configuration from environment variables if not set
	apiKey := config.APIKey.ValueString()
	if apiKey == "" {
		apiKey = os.Getenv("DOMOTZ_API_KEY")
	}

	baseURL := config.BaseURL.ValueString()
	if baseURL == "" {
		baseURL = os.Getenv("DOMOTZ_BASE_URL")
	}
	if baseURL == "" {
		baseURL = defaultBaseURL
	}

	// Validate required configuration
	if apiKey == "" {
		resp.Diagnostics.AddError(
			"Missing API Key Configuration",
			"The provider cannot create the Domotz API client as there is a missing or empty value for the API key. "+
				"Set the api_key value in the provider configuration or use the DOMOTZ_API_KEY environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
		return
	}

	// Create API client
	c := client.NewClient(baseURL, apiKey)

	// Make the client available to resources and data sources
	resp.DataSourceData = c
	resp.ResourceData = c
}

// Resources defines the resources implemented in the provider
func (p *DomotzProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewDeviceResource,
		NewCustomTagResource,
		NewDeviceTagBindingResource,
		NewSNMPSensorResource,
		NewTCPSensorResource,
	}
}

// DataSources defines the data sources implemented in the provider
func (p *DomotzProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewAgentDataSource,
		NewDeviceDataSource,
		NewDevicesDataSource,
		NewDeviceVariablesDataSource,
	}
}

// New returns a new provider instance
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &DomotzProvider{
			version: version,
		}
	}
}
