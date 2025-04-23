package provider

import (
	"context"
	"fmt"
	"terraform-provider-provs/internal/client"
	"terraform-provider-provs/internal/model"

	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ ephemeral.EphemeralResource              = &secretEphemeral{}
	_ ephemeral.EphemeralResourceWithConfigure = &secretEphemeral{}
)

func NewEphemeralSecret() ephemeral.EphemeralResource {
	return &secretEphemeral{}
}

// This ephemeral resource is responsible with reading the secret from the secret manager to be able to be used as an ephemeral value
type secretEphemeral struct {
	client client.Client[*model.SecretManager]
}

func (r *secretEphemeral) Metadata(_ context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + typeSecret
}

// Schema defines the schema for the resource.
func (r *secretEphemeral) Schema(_ context.Context, _ ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"secret_manager_id": schema.StringAttribute{
				Required: true,
			},
			"secret_name": schema.StringAttribute{
				Required: true,
			},
			"secret": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *secretEphemeral) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	var cfg secretModel
	diags := req.Config.Get(ctx, &cfg)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	mgr, err := r.client.GetByID(cfg.SecretManagerID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to retrieve the secret manager",
			fmt.Sprintf("Error retrieving secret manager for id %q: %s", cfg.SecretManagerID.ValueString(), err),
		)
		return
	}
	secret, ok := mgr.Secrets[cfg.SecretName.ValueString()]
	if !ok {
		resp.Diagnostics.AddError(
			"Failed to retrieve the secret from secret manager",
			fmt.Sprintf(
				"Error retrieving secret %q from secret manager for id %q: no such secret existing",
				cfg.SecretName.ValueString(),
				cfg.SecretManagerID.ValueString()),
		)
		return
	}
	cfg.Secret = types.StringValue(secret)
	resp.Diagnostics.Append(resp.Result.Set(ctx, &cfg)...)
}

// Configure adds the provider configured client to the resource.
func (r *secretEphemeral) Configure(_ context.Context, req ephemeral.ConfigureRequest, resp *ephemeral.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(client.BackendClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Secret Manager Client Type",
			fmt.Sprintf("Expected client.BackendClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	// typeSecretManager because it reads data from there
	r.client = client.NewClient[*model.SecretManager](c, typeSecretManager)
}
