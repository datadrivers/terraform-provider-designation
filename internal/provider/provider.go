package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure DesignationProvider satisfies various provider interfaces.
var _ provider.Provider = &DesignationProvider{}

// DesignationProvider satisfies the provider.Provider interface and usually is included
// with all Resource and DataSource implementations.
type DesignationProvider struct {

	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// DesignationProviderModel describes the provider data model.
type DesignationProviderModel struct{}

func (p *DesignationProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "designation"
	resp.Version = p.version
}

// Schema defines the arguments and attributes of this provider
func (p *DesignationProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{}
}

func (p *DesignationProvider) Configure(_ context.Context, _ provider.ConfigureRequest, _ *provider.ConfigureResponse) {

}

// Resources - Defines provider resources
func (p *DesignationProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewConventionResource,
		NewNameResource,
	}
}

// DataSources - Defines provider data sources
func (p *DesignationProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &DesignationProvider{
			version: version,
		}
	}
}
