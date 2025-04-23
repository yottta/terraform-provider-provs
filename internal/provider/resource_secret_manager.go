package provider

import (
	"context"
	"fmt"
	"terraform-provider-provs/internal/client"
	"terraform-provider-provs/internal/model"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &secretManagerResource{}
	_ resource.ResourceWithConfigure = &secretManagerResource{}
)

// orderResourceModel maps the resource schema data.
type secretManagerModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func NewResourceSecretManager() resource.Resource {
	return &secretManagerResource{}
}

// This resource is for managing the existence of a secret manager. For handling the secret stored, check ephemeral/secret_manager.go
type secretManagerResource struct {
	client client.Client[*model.SecretManager]
}

// Metadata returns the resource type name.
func (r *secretManagerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + typeSecretManager
}

// Schema defines the schema for the resource.
func (r *secretManagerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

// Create a new resource.
func (r *secretManagerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan secretManagerModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	item := model.SecretManager{
		ID:   uuid.NewString(),
		Name: plan.Name.ValueString(),
	}

	if _, err := r.client.Create(&item); err != nil {
		resp.Diagnostics.AddError(
			"Error creating secret manager resource",
			fmt.Sprintf("Could not create secret manager: %s ", err),
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
func (r *secretManagerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state secretManagerModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	mgr, err := r.client.GetByID(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading SecretManager",
			fmt.Sprintf("Could not read SecretManager ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}

	state.Name = types.StringValue(mgr.Name)
	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *secretManagerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan secretManagerModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	mgr, err := r.client.GetByID(plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading SecretManager",
			fmt.Sprintf("Could not read SecretManager ID %s: %s", plan.ID.ValueString(), err),
		)
		return
	}
	mgr.Name = plan.Name.ValueString()
	if err := r.client.Update(mgr); err != nil {
		resp.Diagnostics.AddError(
			"Error Updating SecretManager",
			fmt.Sprintf("Could not update SecretManager ID %s: %s", plan.ID.ValueString(), err),
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
func (r *secretManagerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from plan
	var state secretManagerModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.Delete(state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting SecretManager",
			fmt.Sprintf("Failed to delete the secret manager with id %q: %s", state.ID.ValueString(), err),
		)
		return
	}
}

func (r *secretManagerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Configure adds the provider configured client to the resource.
func (r *secretManagerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = client.NewClient[*model.SecretManager](c, typeSecretManager)
}
