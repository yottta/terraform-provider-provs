package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ ephemeral.EphemeralResource          = &randomEphemeral{}
	_ ephemeral.EphemeralResourceWithRenew = &randomEphemeral{}
)

// orderResourceModel maps the resource schema data.
type randomModel struct {
	ID     types.String `tfsdk:"id"`
	Value  types.String `tfsdk:"value"`
	Prefix types.String `tfsdk:"prefix"`
}

func NewEphemeralRandom() ephemeral.EphemeralResource {
	return &randomEphemeral{}
}

type randomEphemeral struct{}

func (r *randomEphemeral) Metadata(_ context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + typeRandom
}

// Schema defines the schema for the resource.
func (r *randomEphemeral) Schema(_ context.Context, _ ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"value": schema.StringAttribute{
				Computed: true,
			},
			"prefix": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

func (r *randomEphemeral) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	var cfg randomModel
	diags := req.Config.Get(ctx, &cfg)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res := generateSecret([]byte(cfg.Prefix.ValueString()))
	cfg.Value = types.StringValue(res)
	resp.Diagnostics.Append(resp.Result.Set(ctx, &cfg)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.RenewAt = time.Now().Add(2 * time.Minute)
}

func (r *randomEphemeral) Renew(ctx context.Context, req ephemeral.RenewRequest, resp *ephemeral.RenewResponse) {
	val, diags := req.Private.GetKey(ctx, "prefix")
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	res := generateSecret(val)

	resp.RenewAt = time.Now().Add(2 * time.Minute)
	resp.Diagnostics.Append(resp.Private.SetKey(ctx, "value", []byte(res))...)
}

func generateSecret(prefix []byte) string {
	res := uuid.NewString()
	if len(prefix) > 0 {
		res = fmt.Sprintf("%s-%s", prefix, res)
	}
	return res
}
