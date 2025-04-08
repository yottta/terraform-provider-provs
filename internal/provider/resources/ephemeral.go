package resources

import (
	"context"
	"fmt"
	"terraform-provider-provs/internal/client"
	"terraform-provider-provs/internal/model"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	secretResType = "secret"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ ephemeral.EphemeralResource              = &secretResource{}
	_ ephemeral.EphemeralResourceWithRenew     = &secretResource{}
	_ ephemeral.EphemeralResourceWithConfigure = &secretResource{}
)

// orderResourceModel maps the resource schema data.
type secretModel struct {
	ID    types.String `tfsdk:"id"`
	Value types.String `tfsdk:"value"`
	Name  types.String `tfsdk:"name"`
}

func NewSecretResource() ephemeral.EphemeralResource {
	return &secretResource{}
}

type secretResource struct {
	client client.Client[*model.WriteOnly]
}

func (r *secretResource) Metadata(_ context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + secretResType
}

// Schema defines the schema for the resource.
func (r *secretResource) Schema(_ context.Context, _ ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"value": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

func (r *secretResource) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	var cfg secretModel
	diags := req.Config.Get(ctx, &cfg)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	res := fmt.Sprintf("%s-%s", cfg.Name.ValueString(), uuid.NewString())

	cfg.Value = types.StringValue(res)
	resp.Diagnostics.Append(resp.Result.Set(ctx, &cfg)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.RenewAt = time.Now().Add(2 * time.Minute)
}

func (r *secretResource) Renew(ctx context.Context, req ephemeral.RenewRequest, resp *ephemeral.RenewResponse) {
	val, diags := req.Private.GetKey(ctx, "name")
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	res := fmt.Sprintf("%s-%s", val, uuid.NewString())

	resp.RenewAt = time.Now().Add(2 * time.Minute)
	resp.Diagnostics.Append(resp.Private.SetKey(ctx, "value", []byte(res))...)
}

// Configure adds the provider configured client to the resource.
func (r *secretResource) Configure(_ context.Context, req ephemeral.ConfigureRequest, resp *ephemeral.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(client.BackendClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected client.BackendClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client.NewClient[*model.WriteOnly](c, secretResType)
}
