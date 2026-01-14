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
var _ resource.Resource = &ProjectServiceAccountResource{}
var _ resource.ResourceWithImportState = &ProjectServiceAccountResource{}

func NewProjectServiceAccountResource() resource.Resource {
	return &ProjectServiceAccountResource{}
}

// ProjectServiceAccountResource defines the resource implementation.
type ProjectServiceAccountResource struct {
	client *client.Client
}

// ProjectServiceAccountResourceModel describes the resource data model.
type ProjectServiceAccountResourceModel struct {
	ID          types.String `tfsdk:"id"`
	ProjectID   types.String `tfsdk:"project_id"`
	Name        types.String `tfsdk:"name"`
	Role        types.String `tfsdk:"role"`
	CreatedAt   types.Int64  `tfsdk:"created_at"`
	APIKeyID    types.String `tfsdk:"api_key_id"`
	APIKeyValue types.String `tfsdk:"api_key_value"`
}

func (r *ProjectServiceAccountResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project_service_account"
}

func (r *ProjectServiceAccountResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates an OpenAI service account within a project. Service accounts are used for programmatic API access and are created with an associated API key.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the service account.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"project_id": schema.StringAttribute{
				Description: "The ID of the project to create the service account in.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the service account.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"role": schema.StringAttribute{
				Description: "The role of the service account. Currently always 'member'.",
				Computed:    true,
			},
			"created_at": schema.Int64Attribute{
				Description: "The Unix timestamp (in seconds) of when the service account was created.",
				Computed:    true,
			},
			"api_key_id": schema.StringAttribute{
				Description: "The ID of the API key associated with this service account.",
				Computed:    true,
			},
			"api_key_value": schema.StringAttribute{
				Description: "The API key value. Only available immediately after creation.",
				Computed:    true,
				Sensitive:   true,
			},
		},
	}
}

func (r *ProjectServiceAccountResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ProjectServiceAccountResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ProjectServiceAccountResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sa, err := r.client.CreateServiceAccount(ctx, data.ProjectID.ValueString(), &client.CreateServiceAccountRequest{
		Name: data.Name.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create service account: %s", err))
		return
	}

	data.ID = types.StringValue(sa.ID)
	data.Name = types.StringValue(sa.Name)
	data.Role = types.StringValue(sa.Role)
	data.CreatedAt = types.Int64Value(sa.CreatedAt)
	data.APIKeyID = types.StringValue(sa.APIKey.ID)
	data.APIKeyValue = types.StringValue(sa.APIKey.Value)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProjectServiceAccountResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ProjectServiceAccountResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sa, err := r.client.GetServiceAccount(ctx, data.ProjectID.ValueString(), data.ID.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read service account: %s", err))
		return
	}

	data.ID = types.StringValue(sa.ID)
	data.Name = types.StringValue(sa.Name)
	data.Role = types.StringValue(sa.Role)
	data.CreatedAt = types.Int64Value(sa.CreatedAt)
	// API key value is only available at creation time, preserve from state
	// API key ID is not returned on read, preserve from state

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProjectServiceAccountResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Service accounts cannot be updated, only recreated
	resp.Diagnostics.AddError(
		"Update Not Supported",
		"Service accounts cannot be updated. Any changes require recreating the resource.",
	)
}

func (r *ProjectServiceAccountResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ProjectServiceAccountResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteServiceAccount(ctx, data.ProjectID.ValueString(), data.ID.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete service account: %s", err))
		return
	}
}

func (r *ProjectServiceAccountResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import format: project_id/service_account_id
	parts := strings.Split(req.ID, "/")
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Expected import ID format: project_id/service_account_id, got: %s", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[1])...)
}
