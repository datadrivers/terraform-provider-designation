package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ tfsdk.ResourceType = nameResourceType{}
var _ tfsdk.Resource = nameResource{}

type nameResourceType struct{}

// Convention Resource schema
func (r nameResourceType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "This resource is used to get a name with a configured convention",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Computed:            true,
				MarkdownDescription: "The name identifier",
				Type:                types.StringType,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
			},
			"name": {
				MarkdownDescription: "This is the required convention option for the name",
				Required:            true,
				Type:                types.StringType,
			},
			"inputs": {
				MarkdownDescription: "Map of input values for variables in provider defined convention",
				Required:            true,
				Type: types.MapType{
					ElemType: types.StringType,
				},
			},
			"convention": {
				MarkdownDescription: "The validated convention formated as a json string",
				Required:            true,
				Type:                types.StringType,
			},
			"result": {
				Computed:            true,
				MarkdownDescription: "The result is the generated name.",
				Type:                types.StringType,
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

	var resourceData nameResourceData
	diags := req.Config.Get(ctx, &resourceData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var convention Convention
	if err := json.Unmarshal([]byte(resourceData.Convention.Value), &convention); err != nil {
		resp.Diagnostics.AddError(
			"Convention Reading Error",
			err.Error(),
		)
		return
	}

	inputs := map[string]string{}
	for key, value := range resourceData.Inputs.Elems {
		inputs[key] = value.(types.String).Value
	}

	diags = validateInputs(inputs, convention.Variables)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, diags := generateName(resourceData.Name.Value, inputs, convention)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceData.ID = types.String{Value: fmt.Sprintf("%s/%s", convention.Definition, resourceData.Name.Value)}
	resourceData.Result = result

	diags = resp.State.Set(ctx, resourceData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r nameResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var resourceData nameResourceData
	diags := req.State.Get(ctx, &resourceData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &resourceData)
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

	var resourceData nameResourceData
	diags := req.Config.Get(ctx, &resourceData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var convention Convention
	if err := json.Unmarshal([]byte(resourceData.Convention.Value), &convention); err != nil {
		resp.Diagnostics.AddError(
			"Convention Reading Error",
			err.Error(),
		)
		return
	}

	inputs := map[string]string{}
	for key, value := range resourceData.Inputs.Elems {
		inputs[key] = value.(types.String).Value
	}

	diags = validateInputs(inputs, convention.Variables)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, diags := generateName(resourceData.Name.Value, inputs, convention)
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
func (r nameResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	resp.State.RemoveResource(ctx)
}
