package provider

import (
	"context"
	"fmt"
	"terraform-provider-provs/internal/client"
	"terraform-provider-provs/internal/model"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &secretResource{}
	_ resource.ResourceWithConfigure   = &secretResource{}
	_ resource.ResourceWithImportState = &secretResource{}
)

// orderResourceModel maps the resource schema data.
type secretModel struct {
	SecretManagerID types.String `tfsdk:"secret_manager_id"`
	SecretName      types.String `tfsdk:"secret_name"`
	Secret          types.String `tfsdk:"secret"`
	SecretWO        types.String `tfsdk:"secret_wo"`
	HasSecretWO     types.Bool   `tfsdk:"has_secret_wo"`
	SecretWOVersion types.Int32  `tfsdk:"secret_wo_version"`
}

func NewResourceSecret() resource.Resource {
	return &secretResource{}
}

type secretResource struct {
	client client.Client[*model.SecretManager]
}

// Metadata returns the resource type name.
func (r *secretResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + typeSecret
}

// Schema defines the schema for the resource.
func (r *secretResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"secret_manager_id": schema.StringAttribute{
				Required: true,
			},
			"secret_name": schema.StringAttribute{
				Required: true,
			},
			"secret": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.PreferWriteOnlyAttribute(path.MatchRoot("secret_wo")),
					stringvalidator.ExactlyOneOf(path.MatchRoot("secret_wo")),
				},
			},
			"secret_wo": schema.StringAttribute{
				Optional:  true,
				WriteOnly: true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot("secret")),
					//stringvalidator.AlsoRequires(path.MatchRoot("has_attr_wo")),
					stringvalidator.AlsoRequires(path.MatchRoot("secret_wo_version")),
				},
			},
			"has_secret_wo": schema.BoolAttribute{
				Computed: true,
			},
			"secret_wo_version": schema.Int32Attribute{
				Optional: true,
			},
		},
	}
}

// Create a new resource.
func (r *secretResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan, config secretModel
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

	mgr, err := r.client.GetByID(plan.SecretManagerID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating secret resource",
			fmt.Sprintf("Could not find secret manager for id %q: %s ", plan.SecretManagerID.ValueString(), err),
		)
		return
	}
	if mgr.Secrets == nil {
		mgr.Secrets = map[string]string{}
	}
	mgr.Secrets[plan.SecretName.ValueString()] = plan.Secret.ValueString()
	// Generate API request body from plan
	plan.HasSecretWO = types.BoolValue(false)
	// Secret attributes should be read only from the req.config, not from req.plan
	if !config.SecretWOVersion.IsNull() {
		mgr.Secrets[plan.SecretName.ValueString()] = config.SecretWO.ValueString()
		plan.HasSecretWO = types.BoolValue(true)
	}

	if err := r.client.Update(mgr); err != nil {
		resp.Diagnostics.AddError(
			"Error updating secret manager resource",
			"Could not update secret manager, unexpected error: "+err.Error(),
		)
		return
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information.
func (r *secretResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state secretModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed mgr value
	mgr, err := r.client.GetByID(state.SecretManagerID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Secret manager",
			fmt.Sprintf("Could not read secret manager with id %q: %s", state.SecretManagerID.ValueString(), err),
		)
		return
	}
	addMissingSecretDiag := func() {
		resp.Diagnostics.AddError(
			"Error Reading secret from secret manager",
			fmt.Sprintf(
				"Could not read secret %q from secret manager with id %q: no such secret",
				state.SecretName.ValueString(),
				state.SecretManagerID.ValueString(),
			),
		)
	}
	if mgr.Secrets == nil {
		addMissingSecretDiag()
		return
	}
	secret, ok := mgr.Secrets[state.SecretName.ValueString()]
	if !ok {
		addMissingSecretDiag()
		return
	}
	state.Secret = types.StringValue(secret)
	if state.HasSecretWO.ValueBool() { // DO NOT RETURN WRITE ONLY VALUES
		state.Secret = types.StringNull()
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *secretResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan, config secretModel
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

	mgr, err := r.client.GetByID(plan.SecretManagerID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating secret resource",
			fmt.Sprintf("Could not find secret manager for id %q: %s ", plan.SecretManagerID.ValueString(), err),
		)
		return
	}

	addMissingSecretDiag := func() {
		resp.Diagnostics.AddError(
			"Error updating secret inside secret manager",
			fmt.Sprintf(
				"The secret manager %q contains no secret named %q",
				plan.SecretManagerID.ValueString(),
				plan.SecretName.ValueString(),
			),
		)
	}
	if mgr.Secrets == nil {
		addMissingSecretDiag()
		return
	}
	if _, ok := mgr.Secrets[plan.SecretName.ValueString()]; !ok {
		addMissingSecretDiag()
		return
	}
	mgr.Secrets[plan.SecretName.ValueString()] = plan.Secret.ValueString()
	plan.HasSecretWO = types.BoolValue(false)
	if !config.SecretWO.IsNull() {
		mgr.Secrets[plan.SecretName.ValueString()] = plan.SecretWO.ValueString()
		plan.HasSecretWO = types.BoolValue(true)
	}

	// Update existing res
	if err := r.client.Update(mgr); err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Secret Manager",
			fmt.Sprintf("Could not update SecretManager with id %q: %s", plan.SecretManagerID.ValueString(), err),
		)
		return
	}

	// Here is not needed to handle unsetting the secret field since the provider framework is already doing it based on the schema
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *secretResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from plan
	var state secretModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	mgr, err := r.client.GetByID(state.SecretManagerID.ValueString())
	if err != nil { // TODO improve errors to check for not found. If secret manager was deleted should return no error
		resp.Diagnostics.AddError(
			"Error deleting secret resource",
			fmt.Sprintf("Could not find secret manager for id %q: %s ", state.SecretManagerID.ValueString(), err),
		)
		return
	}

	delete(mgr.Secrets, state.SecretName.ValueString())
	// Update existing res
	if err := r.client.Update(mgr); err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Secret Manager",
			fmt.Sprintf("Could not update SecretManager with id %q: %s", state.SecretManagerID.ValueString(), err),
		)
		return
	}
}

func (r *secretResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp) // TODO this might be wrong
}

// Configure adds the provider configured client to the resource.
func (r *secretResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	// typeSecretManager because it updates that data
	r.client = client.NewClient[*model.SecretManager](c, typeSecretManager)
}
