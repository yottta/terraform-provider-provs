package provider

import (
	"context"
	"os"
	"terraform-provider-provs/internal/client/filesystem"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider                       = &provsProvider{}
	_ provider.ProviderWithEphemeralResources = &provsProvider{}
	_ provider.ProviderWithFunctions          = &provsProvider{}
)

// provsProviderModel maps provider schema data to a Go type.
type provsProviderModel struct {
	Path types.String `tfsdk:"path"`
}

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &provsProvider{
			version: version,
		}
	}
}

// provsProvider is the provider implementation.
type provsProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// Metadata returns the provider type name.
func (p *provsProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "provs"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *provsProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"path": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

func (p *provsProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring provider client")
	// Retrieve provider data from configuration
	var config provsProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.
	if config.Path.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("path"),
			"Unknown path for storing provider data",
			"The provider cannot create the fs client as there is an unknown configuration value for the the storage path. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the PROVS_PATH environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.
	storagePath := os.Getenv("PROVS_PATH")
	if !config.Path.IsNull() {
		storagePath = config.Path.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if storagePath == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("path"),
			"Missing storage path",
			"The provider cannot create the Provs API client as there is a missing or empty value for the Provs API host. "+
				"Set the host value in the configuration or use the PROVS_PATH environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}
	ctx = tflog.SetField(ctx, "provs_path", storagePath)

	tflog.Debug(ctx, "Creating Provs client")

	// Create a new Provs client using the configuration values
	c, err := filesystem.NewFsClient(storagePath)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create storage client",
			"An unexpected error occurred when creating the Provs API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Client Error: "+err.Error(),
		)
		return
	}

	// Make the client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = c
	resp.ResourceData = c
	tflog.Info(ctx, "Configured Provs client", map[string]any{"success": true})
}

// DataSources defines the data sources implemented in the provider.
func (p *provsProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewCoffeesDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *provsProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewResourceOrder,
		NewResourceIssue2372,
		NewResourceSecret,
		NewResourceSecretManager,
	}
}

func (p *provsProvider) Functions(context.Context) []func() function.Function {
	return []func() function.Function{
		NewFunctionComputeTax,
	}
}

func (p *provsProvider) EphemeralResources(_ context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{
		NewEphemeralRandom,
		NewEphemeralSecret,
	}
}
