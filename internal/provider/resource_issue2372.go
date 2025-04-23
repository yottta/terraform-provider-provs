package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource = &issue2372Resource{}
)

// NewResourceIssue2372 is a helper function to simplify the provider implementation.
func NewResourceIssue2372() resource.Resource {
	return &issue2372Resource{}
}

type issue2372Resource struct {
}

type issue2372Model struct {
	Id          types.String `tfsdk:"id"`
	NestedValue types.Object `tfsdk:"nested_value"`
}

type issue2372SingleNestedValResourceModel struct {
	ChildComputed types.String `tfsdk:"child_computed"`
}

// GetSchema defines the schema for the resource.
func (r *issue2372Resource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Issue 2372 resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Issue 2372 identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"nested_value": schema.SingleNestedAttribute{
				Description: "This contains some nested values",
				Optional:    true,

				Attributes: map[string]schema.Attribute{
					"child_computed": schema.StringAttribute{
						Description: "This is a computed child attribute",
						Computed:    true,
					},
				},
			},
		},
	}
}

// Metadata returns the resource type name.
func (r *issue2372Resource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_issue2372"
}

func (r *issue2372Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state issue2372Model

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Stub API
	data := plan
	data.SetId()

	resp.Diagnostics.Append(state.toState(ctx, data)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *issue2372Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data, state issue2372Model

	// Stub API
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data = state

	resp.Diagnostics.Append(state.toState(ctx, data)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *issue2372Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state issue2372Model

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Stub API
	data := plan
	data.SetId()

	resp.Diagnostics.Append(state.toState(ctx, data)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

// No backend so no logic needed
func (r *issue2372Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *issue2372Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (m *issue2372Model) SetId() {
	m.Id = types.StringValue("id")
}

func (m *issue2372Model) toState(ctx context.Context, data issue2372Model) (diags diag.Diagnostics) {
	m.Id = data.Id

	singleNestedObjectTyp := map[string]attr.Type{
		"child_computed": types.StringType,
	}

	var d diag.Diagnostics

	singleNestedPlan := issue2372SingleNestedValResourceModel{
		ChildComputed: types.StringValue("required string"),
	}

	if !data.NestedValue.IsNull() {
		d = data.NestedValue.As(context.Background(), &singleNestedPlan, basetypes.ObjectAsOptions{})
		diags.Append(d...)
		if diags.HasError() {
			return
		}
	}

	m.NestedValue, d = types.ObjectValueFrom(ctx, singleNestedObjectTyp, singleNestedPlan)
	diags = append(diags, d...)

	return
}
