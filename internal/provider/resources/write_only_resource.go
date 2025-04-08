package resources

import (
	"context"
	"fmt"
	"terraform-provider-provs/internal/client"
	"terraform-provider-provs/internal/model"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	writeOnlyResType = "write_only_test"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &writeOnlyResource{}
	_ resource.ResourceWithConfigure   = &writeOnlyResource{}
	_ resource.ResourceWithImportState = &writeOnlyResource{}
)

// orderResourceModel maps the resource schema data.
type writeOnlyModel struct {
	ID            types.String `tfsdk:"id"`
	Attribute     types.String `tfsdk:"attr"`
	WOAttribute   types.String `tfsdk:"attr_wo"`
	HasAttrWO     types.Bool   `tfsdk:"has_attr_wo"`
	AttrVersionWO types.Int32  `tfsdk:"attr_wo_version"`
}

func NewWriteOnlyResource() resource.Resource {
	return &writeOnlyResource{}
}

type writeOnlyResource struct {
	client client.Client[*model.WriteOnly]
}

// Metadata returns the resource type name.
func (r *writeOnlyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + writeOnlyResType
}

// Schema defines the schema for the resource.
func (r *writeOnlyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"attr": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.PreferWriteOnlyAttribute(path.MatchRoot("attr_wo")),
					stringvalidator.ExactlyOneOf(path.MatchRoot("attr_wo")),
				},
			},
			"attr_wo": schema.StringAttribute{
				Optional:  true,
				WriteOnly: true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot("attr")),
					//stringvalidator.AlsoRequires(path.MatchRoot("has_attr_wo")),
					stringvalidator.AlsoRequires(path.MatchRoot("attr_wo_version")),
				},
			},
			"has_attr_wo": schema.BoolAttribute{
				Computed: true,
			},
			"attr_wo_version": schema.Int32Attribute{
				Optional: true,
			},
		},
	}
}

// Create a new resource.
func (r *writeOnlyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan, config writeOnlyModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	item := model.WriteOnly{
		ID:   uuid.NewString(),
		Attr: plan.Attribute.ValueString(),
	}
	plan.HasAttrWO = types.BoolValue(false)
	// WriteOnly attributes should be read only from the req.config, not from req.plan
	if !config.WOAttribute.IsNull() {
		item.Attr = config.WOAttribute.ValueString()
		plan.HasAttrWO = types.BoolValue(true)
	}

	if _, err := r.client.Create(&item); err != nil {
		resp.Diagnostics.AddError(
			"Error creating writeOnlyRes resource",
			"Could not create writeOnlyRes, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.StringValue(item.ID)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information.
func (r *writeOnlyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state writeOnlyModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed writeOnlyRes value
	writeOnlyRes, err := r.client.GetByID(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading WriteOnlyRes",
			"Could not read WriteOnlyRes ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}
	state.Attribute = types.StringValue(writeOnlyRes.Attr)
	if state.HasAttrWO.ValueBool() {
		state.Attribute = types.StringNull()
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *writeOnlyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan, config writeOnlyModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	o := model.WriteOnly{
		ID:   plan.ID.ValueString(),
		Attr: plan.Attribute.ValueString(),
	}
	plan.HasAttrWO = types.BoolValue(false)
	if !config.WOAttribute.IsNull() {
		o.Attr = config.WOAttribute.ValueString()
		plan.HasAttrWO = types.BoolValue(true)
	}

	// Update existing res
	if err := r.client.Update(&o); err != nil {
		resp.Diagnostics.AddError(
			"Error Updating WriteOnlyRes",
			"Could not update writeOnlyRes, unexpected error: "+err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *writeOnlyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from plan
	var state writeOnlyModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.Delete(state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting WriteOnlyRes",
			"Could not delete writeOnlyRes, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *writeOnlyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Configure adds the provider configured client to the resource.
func (r *writeOnlyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
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

	r.client = client.NewClient[*model.WriteOnly](c, writeOnlyResType)
}
