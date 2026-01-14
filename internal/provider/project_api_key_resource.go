package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/terraform-mars/terraform-provider-openai/internal/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &ProjectAPIKeyResource{}
var _ resource.ResourceWithImportState = &ProjectAPIKeyResource{}

func NewProjectAPIKeyResource() resource.Resource {
	return &ProjectAPIKeyResource{}
}

// ProjectAPIKeyResource defines the resource implementation.
// Note: API keys in OpenAI are typically created via service accounts or the UI.
// This resource allows managing existing API keys (read/delete).
type ProjectAPIKeyResource struct {
	client *client.Client
}

// ProjectAPIKeyResourceModel describes the resource data model.
type ProjectAPIKeyResourceModel struct {
	ID            types.String `tfsdk:"id"`
	ProjectID     types.String `tfsdk:"project_id"`
	Name          types.String `tfsdk:"name"`
	RedactedValue types.String `tfsdk:"redacted_value"`
	CreatedAt     types.Int64  `tfsdk:"created_at"`
	OwnerType     types.String `tfsdk:"owner_type"`
	OwnerID       types.String `tfsdk:"owner_id"`
	OwnerName     types.String `tfsdk:"owner_name"`
}

func (r *ProjectAPIKeyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project_api_key"
}

func (r *ProjectAPIKeyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an OpenAI project API key. Note: API keys are created via service accounts. Use openai_project_service_account to create new API keys. This resource is primarily for importing and managing existing keys.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the API key.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"project_id": schema.StringAttribute{
				Description: "The ID of the project this API key belongs to.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the API key.",
				Computed:    true,
			},
			"redacted_value": schema.StringAttribute{
				Description: "The redacted value of the API key (e.g., 'sk-...abc').",
				Computed:    true,
			},
			"created_at": schema.Int64Attribute{
				Description: "The Unix timestamp (in seconds) of when the API key was created.",
				Computed:    true,
			},
			"owner_type": schema.StringAttribute{
				Description: "The type of owner of this API key ('user' or 'service_account').",
				Computed:    true,
			},
			"owner_id": schema.StringAttribute{
				Description: "The ID of the owner (user or service account).",
				Computed:    true,
			},
			"owner_name": schema.StringAttribute{
				Description: "The name of the owner.",
				Computed:    true,
			},
		},
	}
}

func (r *ProjectAPIKeyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = c
}

func (r *ProjectAPIKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ProjectAPIKeyResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// API keys cannot be created directly - they're created via service accounts
	// This resource is for importing existing keys
	// On "create", we just read the key to verify it exists
	key, err := r.client.GetProjectAPIKey(ctx, data.ProjectID.ValueString(), data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read API key: %s. Note: API keys must be created via service accounts.", err))
		return
	}

	data.Name = types.StringValue(key.Name)
	data.RedactedValue = types.StringValue(key.RedactedValue)
	data.CreatedAt = types.Int64Value(key.CreatedAt)
	data.OwnerType = types.StringValue(key.Owner.Type)

	if key.Owner.User != nil {
		data.OwnerID = types.StringValue(key.Owner.User.ID)
		data.OwnerName = types.StringValue(key.Owner.User.Name)
	} else if key.Owner.ServiceAccount != nil {
		data.OwnerID = types.StringValue(key.Owner.ServiceAccount.ID)
		data.OwnerName = types.StringValue(key.Owner.ServiceAccount.Name)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProjectAPIKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ProjectAPIKeyResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	key, err := r.client.GetProjectAPIKey(ctx, data.ProjectID.ValueString(), data.ID.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read API key: %s", err))
		return
	}

	data.Name = types.StringValue(key.Name)
	data.RedactedValue = types.StringValue(key.RedactedValue)
	data.CreatedAt = types.Int64Value(key.CreatedAt)
	data.OwnerType = types.StringValue(key.Owner.Type)

	if key.Owner.User != nil {
		data.OwnerID = types.StringValue(key.Owner.User.ID)
		data.OwnerName = types.StringValue(key.Owner.User.Name)
	} else if key.Owner.ServiceAccount != nil {
		data.OwnerID = types.StringValue(key.Owner.ServiceAccount.ID)
		data.OwnerName = types.StringValue(key.Owner.ServiceAccount.Name)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProjectAPIKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// API keys cannot be updated
	resp.Diagnostics.AddError(
		"Update Not Supported",
		"API keys cannot be updated. Any changes require recreating the resource.",
	)
}

func (r *ProjectAPIKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ProjectAPIKeyResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteProjectAPIKey(ctx, data.ProjectID.ValueString(), data.ID.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete API key: %s", err))
		return
	}
}

func (r *ProjectAPIKeyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import format: project_id/api_key_id
	parts := strings.Split(req.ID, "/")
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Expected import ID format: project_id/api_key_id, got: %s", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[1])...)
}
