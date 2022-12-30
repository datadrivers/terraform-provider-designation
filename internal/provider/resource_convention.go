package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &ConventionResource{}

func NewConventionResource() resource.Resource {
	return &ConventionResource{}
}

type ConventionResource struct{}

func (r *ConventionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_convention"
}

// Schema returns Convention Resource schema
func (r ConventionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This resource is used to define a convention. The convention will be validated.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The name identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"definition": schema.StringAttribute{
				MarkdownDescription: "The definition of the convention. Must include the block `(name)` and all variable blocks.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
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
						"generated": schema.BoolAttribute{
							MarkdownDescription: "Activates the generation of a random string",
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

func (r *ConventionResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}
}

// Create a new resource
func (r ConventionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var resourceData ConventionResourceData

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &resourceData)...)

	if resp.Diagnostics.HasError() {
		return
	}

	diags := req.Config.Get(ctx, &resourceData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = validateConvention(&resourceData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(generateConventionsString(&resourceData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceData.ID = resourceData.Definition

	diags = resp.State.Set(ctx, resourceData)
	resp.Diagnostics.Append(diags...)
}

// Read resource information
func (r ConventionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var resourceData ConventionResourceData

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &resourceData)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &resourceData)...)
}

// Update resource
func (r ConventionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var resourceData ConventionResourceData

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &resourceData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(validateConvention(&resourceData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(generateConventionsString(&resourceData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, resourceData)...)
}

// Delete resource
func (r ConventionResource) Delete(ctx context.Context, _ resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.State.RemoveResource(ctx)
}

// generateConventionsString generates the json string for the convention
func generateConventionsString(resourceData *ConventionResourceData) diag.Diagnostics {
	var diags diag.Diagnostics

	var variables []ConventionVariable
	for _, variable := range resourceData.Variables {
		variables = append(variables, convertVariableToConventionVariable(variable))
	}
	convention, err := json.Marshal(Convention{
		Definition: resourceData.Definition.ValueString(),
		Variables:  variables,
	})
	if err != nil {
		diags.AddError(
			"Convention generation Error",
			err.Error(),
		)
		return diags
	}

	resourceData.Convention = types.StringValue(string(convention))
	return diags
}

// validateConvention checks the configured convention to ensure that it can be used without errors
func validateConvention(resourceData *ConventionResourceData) diag.Diagnostics {
	var diags diag.Diagnostics

	if !strings.Contains(strings.ToLower(resourceData.Definition.ValueString()), "(name)") {
		diags.AddError(
			"Convention Validate Error",
			"The defined convention must include the block '(name)'.",
		)
	}
	var missingVariables []string

	for _, variable := range resourceData.Variables {
		block := fmt.Sprintf("(%s)", variable.Name.ValueString())
		if !strings.Contains(strings.ToLower(resourceData.Definition.ValueString()), block) {
			missingVariables = append(missingVariables, block)
		}

	}
	if len(missingVariables) > 0 {
		diags.AddError(
			"Convention Validate Error",
			fmt.Sprintf("The definied convention must include all variables blocks. Missing blocks: %s", strings.Join(missingVariables, ", ")),
		)
	}

	return diags
}
