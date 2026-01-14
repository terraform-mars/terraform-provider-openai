package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/terraform-mars/terraform-provider-openai/internal/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &ProjectsDataSource{}

func NewProjectsDataSource() datasource.DataSource {
	return &ProjectsDataSource{}
}

// ProjectsDataSource defines the data source implementation.
type ProjectsDataSource struct {
	client *client.Client
}

// ProjectsDataSourceModel describes the data source data model.
type ProjectsDataSourceModel struct {
	IncludeArchived types.Bool            `tfsdk:"include_archived"`
	Projects        []ProjectDataModel    `tfsdk:"projects"`
}

// ProjectDataModel describes a project in the list.
type ProjectDataModel struct {
	ID         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	Status     types.String `tfsdk:"status"`
	CreatedAt  types.Int64  `tfsdk:"created_at"`
	ArchivedAt types.Int64  `tfsdk:"archived_at"`
}

func (d *ProjectsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_projects"
}

func (d *ProjectsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a list of all OpenAI projects in the organization.",
		Attributes: map[string]schema.Attribute{
			"include_archived": schema.BoolAttribute{
				Description: "Whether to include archived projects. Defaults to false.",
				Optional:    true,
			},
			"projects": schema.ListNestedAttribute{
				Description: "List of projects in the organization.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "The unique identifier of the project.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "The name of the project.",
							Computed:    true,
						},
						"status": schema.StringAttribute{
							Description: "The status of the project.",
							Computed:    true,
						},
						"created_at": schema.Int64Attribute{
							Description: "The Unix timestamp of when the project was created.",
							Computed:    true,
						},
						"archived_at": schema.Int64Attribute{
							Description: "The Unix timestamp of when the project was archived.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *ProjectsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = c
}

func (d *ProjectsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ProjectsDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	includeArchived := false
	if !data.IncludeArchived.IsNull() {
		includeArchived = data.IncludeArchived.ValueBool()
	}

	result, err := d.client.ListProjects(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list projects: %s", err))
		return
	}

	data.Projects = make([]ProjectDataModel, 0, len(result.Data))
	for _, project := range result.Data {
		// Filter out archived projects if not requested
		if !includeArchived && project.Status == "archived" {
			continue
		}

		p := ProjectDataModel{
			ID:        types.StringValue(project.ID),
			Name:      types.StringValue(project.Name),
			Status:    types.StringValue(project.Status),
			CreatedAt: types.Int64Value(project.CreatedAt),
		}
		if project.ArchivedAt != nil {
			p.ArchivedAt = types.Int64Value(*project.ArchivedAt)
		} else {
			p.ArchivedAt = types.Int64Null()
		}
		data.Projects = append(data.Projects, p)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
