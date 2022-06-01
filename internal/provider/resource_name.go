package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"k8s.io/apimachinery/pkg/util/rand"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ tfsdk.ResourceType = nameResourceType{}
var _ tfsdk.Resource = nameResource{}

type nameResourceType struct{}

// Convention Resource schema
func (r nameResourceType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "This resource is used to get a name with the in provider configured convention",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Computed:            true,
				MarkdownDescription: "The name identifier",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
				Type: types.StringType,
			},
			"name": {
				MarkdownDescription: "This is the required convention option for the name",
				Type:                types.StringType,
				Required:            true,
			},
			"inputs": {
				MarkdownDescription: "Map of input values for variables in provider defined convention",
				Type: types.MapType{
					ElemType: types.StringType,
				},
				Required: true,
			},
			"result": {
				MarkdownDescription: "The result is the generated name.",
				Type:                types.StringType,
				Computed:            true,
			},
		},
	}, nil
}

// New resource instance
func (r nameResourceType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return nameResource{
		provider: provider,
	}, diags
}

type nameResource struct {
	provider provider
}

// Create a new resource
func (r nameResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	if !r.provider.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, likely because it depends on an unknown value from another resource. This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		)
		return
	}

	var data nameResourceData
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	inputs := map[string]string{}
	for key, value := range data.Inputs.Elems {
		inputs[key] = value.(types.String).Value
	}

	diags = validateInputs(inputs, r.provider.data.Variables)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, diags := generateName(data.Name.Value, inputs, r.provider.data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.ID = data.Name
	data.Result = result

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r nameResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var data nameResourceData
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

// Update resource
func (r nameResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	if !r.provider.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, likely because it depends on an unknown value from another resource. This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		)
		return
	}

	var data nameResourceData
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	inputs := map[string]string{}
	for key, value := range data.Inputs.Elems {
		inputs[key] = value.(types.String).Value
	}

	diags = validateInputs(inputs, r.provider.data.Variables)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, diags := generateName(data.Name.Value, inputs, r.provider.data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Result = result

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete resource
func (r nameResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	resp.State.RemoveResource(ctx)
}

// validateInputs checks if all needed inputs are set
func validateInputs(inputs map[string]string, variables []Variable) diag.Diagnostics {
	var diags diag.Diagnostics

	missingInputs := []string{}
	for _, variable := range variables {
		if variable.Default.Null && !variable.Generated.Value {
			if inputs[variable.Name.Value] != "" {
				missingInputs = append(missingInputs, variable.Name.Value)
			}
		}

	}
	if len(missingInputs) > 0 {
		diags.AddError(
			"Convention Usage Error",
			fmt.Sprintf("All provider variables that are not generated or have a default must be present. Missing inputs: %s", strings.Join(missingInputs, ", ")),
		)
	}

	return diags
}

// generateName generates the name with the inputs, variables and definition
func generateName(name string, inputs map[string]string, providerData *providerData) (types.String, diag.Diagnostics) {
	var diags diag.Diagnostics
	result := strings.Replace(providerData.Definition.Value, "(name)", name, -1)

	for _, variable := range providerData.Variables {
		block := fmt.Sprintf("(%s)", variable.Name.Value)
		replacement := inputs[variable.Name.Value]
		length := 0

		if !variable.MaxLength.Null && variable.MaxLength.Value > int64(0) {
			length = int(variable.MaxLength.Value)
		}
		if variable.Generated.Value {
			replacement = rand.String(length)
		}

		if replacement == "" {
			replacement = variable.Default.Value
		}
		if length > 0 && len(replacement) > length {
			replacement = replacement[0:length]
		}
		result = strings.Replace(result, block, replacement, -1)
	}

	return types.String{Value: result}, diags
}
