package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &NameResource{}

func NewNameResource() resource.Resource {
	return &NameResource{}
}

type NameResource struct{}

func (r *NameResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_name"
}

// Schema returns Name Resource schema
func (r NameResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This resource is used to get a name with a configured convention",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The name identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "This is the required convention option for the name",
				Required:            true,
			},
			"inputs": schema.MapAttribute{
				MarkdownDescription: "Map of input values for variables in provider defined convention",
				Required:            true,
				ElementType:         types.StringType,
			},
			"convention": schema.StringAttribute{
				MarkdownDescription: "The validated convention formated as a json string",
				Required:            true,
			},
			"result": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The result is the generated name.",
			},
		},
	}
}

func (r *NameResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}
}

// Create a new resource
func (r NameResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var resourceData NameResourceData
	diags := req.Config.Get(ctx, &resourceData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var convention Convention
	if err := json.Unmarshal([]byte(resourceData.Convention.ValueString()), &convention); err != nil {
		resp.Diagnostics.AddError(
			"Convention Reading Error",
			err.Error(),
		)
		return
	}

	inputs := map[string]string{}
	for key, value := range resourceData.Inputs.Elements() {
		inputs[key] = value.(types.String).ValueString()
	}

	diags = validateInputs(inputs, convention.Variables)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, diags := generateName(resourceData.Name.ValueString(), inputs, convention)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceData.ID = types.StringValue(fmt.Sprintf("%s/%s", convention.Definition, resourceData.Name.ValueString()))
	resourceData.Result = result

	diags = resp.State.Set(ctx, resourceData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r NameResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var resourceData NameResourceData
	diags := req.State.Get(ctx, &resourceData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &resourceData)
	resp.Diagnostics.Append(diags...)
}

// Update resource
func (r NameResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var resourceData NameResourceData
	diags := req.Config.Get(ctx, &resourceData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var convention Convention
	if err := json.Unmarshal([]byte(resourceData.Convention.ValueString()), &convention); err != nil {
		resp.Diagnostics.AddError(
			"Convention Reading Error",
			err.Error(),
		)
		return
	}

	inputs := map[string]string{}
	for key, value := range resourceData.Inputs.Elements() {
		inputs[key] = value.(types.String).ValueString()
	}

	diags = validateInputs(inputs, convention.Variables)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, diags := generateName(resourceData.Name.ValueString(), inputs, convention)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceData.Result = result

	diags = resp.State.Set(ctx, resourceData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete resource
func (r NameResource) Delete(ctx context.Context, _ resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.State.RemoveResource(ctx)
}
