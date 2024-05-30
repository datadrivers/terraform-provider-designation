package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &ConventionDataSource{}

func NewConventionDataSource() datasource.DataSource {
	return &ConventionDataSource{}
}

type ConventionDataSource struct{}

func (r *ConventionDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_convention"
}

// Schema returns Convention Resource schema
func (r ConventionDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This resource is used to define a convention. The convention will be validated.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The name identifier",
			},
			"definition": schema.StringAttribute{
				MarkdownDescription: "The definition of the convention. Must include the block `(name)` and all variable blocks.",
				Required:            true,
			},
			"variables": schema.ListNestedAttribute{
				MarkdownDescription: "A list of variable definition used in the convention definition.",
				Required:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "Name of the variable",
							Required:            true,
						},
						"default": schema.StringAttribute{
							MarkdownDescription: "Define a default value",
							Optional:            true,
						},
						"max_length": schema.Int64Attribute{
							MarkdownDescription: "Set the size limit of the value. Required if value is generated",
							Optional:            true,
						},
					},
				},
			},
			"convention": schema.StringAttribute{
				MarkdownDescription: "The validated convention formated as a json string",
				Computed:            true,
			},
		},
	}
}

// Read resource information
func (r ConventionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ConventionDataSourceData

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(validateConvention(&data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(generateConventionsString(&data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
