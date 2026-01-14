package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/terraform-mars/terraform-provider-openai/internal/client"
)

// Ensure OpenAIProvider satisfies various provider interfaces.
var _ provider.Provider = &OpenAIProvider{}

// OpenAIProvider defines the provider implementation.
type OpenAIProvider struct {
	version string
}

// OpenAIProviderModel describes the provider data model.
type OpenAIProviderModel struct {
	AdminKey types.String `tfsdk:"admin_key"`
	BaseURL  types.String `tfsdk:"base_url"`
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &OpenAIProvider{
			version: version,
		}
	}
}

func (p *OpenAIProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "openai"
	resp.Version = p.version
}

func (p *OpenAIProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Terraform provider for managing OpenAI organization resources (projects, API keys, service accounts).",
		Attributes: map[string]schema.Attribute{
			"admin_key": schema.StringAttribute{
				Description: "OpenAI Admin API key. Can also be set via OPENAI_ADMIN_KEY environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
			"base_url": schema.StringAttribute{
				Description: "OpenAI API base URL. Defaults to https://api.openai.com/v1",
				Optional:    true,
			},
		},
	}
}

func (p *OpenAIProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config OpenAIProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check for admin key in config or environment
	adminKey := os.Getenv("OPENAI_ADMIN_KEY")
	if !config.AdminKey.IsNull() {
		adminKey = config.AdminKey.ValueString()
	}

	if adminKey == "" {
		resp.Diagnostics.AddError(
			"Missing OpenAI Admin Key",
			"The provider cannot create the OpenAI API client because the admin key is missing. "+
				"Set admin_key in the provider configuration or set the OPENAI_ADMIN_KEY environment variable.",
		)
		return
	}

	// Determine base URL
	baseURL := "https://api.openai.com/v1"
	if envURL := os.Getenv("OPENAI_BASE_URL"); envURL != "" {
		baseURL = envURL
	}
	if !config.BaseURL.IsNull() {
		baseURL = config.BaseURL.ValueString()
	}

	// Create client
	c := client.NewClient(adminKey, baseURL)

	resp.DataSourceData = c
	resp.ResourceData = c
}

func (p *OpenAIProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewProjectResource,
		NewProjectAPIKeyResource,
		NewProjectServiceAccountResource,
	}
}

func (p *OpenAIProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewProjectDataSource,
		NewProjectsDataSource,
	}
}
