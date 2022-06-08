package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ tfsdk.ResourceType = conventionResourceType{}
var _ tfsdk.Resource = conventionResource{}

type conventionResourceType struct{}

// Convention Resource schema
func (r conventionResourceType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "This resource is used to define a convention. The convention will be validated.",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Computed:            true,
				MarkdownDescription: "The name identifier",
				Type:                types.StringType,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
			},
			"definition": {
				MarkdownDescription: "The definition of the convention. Must include the block `(name)` and all variable blocks.",
				Required:            true,
				Type:                types.StringType,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.RequiresReplace(),
				},
			},
			"variables": {
				MarkdownDescription: "A list of variable definition used in the convention definition.",
				Required:            true,
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
					"name": {
						MarkdownDescription: "Name of the variable",
						Required:            true,
						Type:                types.StringType,
					},
					"default": {
						MarkdownDescription: "Define a default value",
						Optional:            true,
						Type:                types.StringType,
					},
					"generated": {
						MarkdownDescription: "Activates the generation of a random string",
						Optional:            true,
						Type:                types.BoolType,
					},
					"max_length": {
						MarkdownDescription: "Set the size limit of the value. Required if value is generated",
						Optional:            true,
						Type:                types.Int64Type,
					},
				}, tfsdk.ListNestedAttributesOptions{}),
			},
			"convention": {
				MarkdownDescription: "The validated convention formated as a json string",
				Computed:            true,
				Type:                types.StringType,
			},
		},
	}, nil
}

// New resource instancexf
func (r conventionResourceType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return conventionResource{
		provider: provider,
	}, diags
}

type conventionResource struct {
	provider provider
}

// Create a new resource
func (r conventionResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	if !r.provider.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, likely because it depends on an unknown value from another resource. This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		)
		return
	}

	var resourceData conventionResourceData
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

	convention, err := json.Marshal(Convention{
		Definition: resourceData.Definition.Value,
		Variables:  resourceData.Variables,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Convention generation Error",
			err.Error(),
		)
		return
	}

	resourceData.Convention = types.String{Value: string(convention)}
	resourceData.ID = resourceData.Definition

	diags = resp.State.Set(ctx, resourceData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r conventionResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var resourceData conventionResourceData
	diags := req.State.Get(ctx, &resourceData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &resourceData)
	resp.Diagnostics.Append(diags...)
}

// Update resource
func (r conventionResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	if !r.provider.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, likely because it depends on an unknown value from another resource. This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		)
		return
	}

	var resourceData conventionResourceData
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

	convention, err := json.Marshal(Convention{
		Definition: resourceData.Definition.Value,
		Variables:  resourceData.Variables,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Convention generation Error",
			err.Error(),
		)
		return
	}

	resourceData.Convention = types.String{Value: string(convention)}

	diags = resp.State.Set(ctx, resourceData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete resource
func (r conventionResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	resp.State.RemoveResource(ctx)
}

// validateConvention checks the configured convention to ensure that it can be used without errors
func validateConvention(resourceData *conventionResourceData) diag.Diagnostics {
	var diags diag.Diagnostics

	if !strings.Contains(strings.ToLower(resourceData.Definition.Value), "(name)") {
		diags.AddError(
			"Convention Validate Error",
			"The definied convention must include the block '(name)'.",
		)
	}
	missingVariables := []string{}

	for _, variable := range resourceData.Variables {
		block := fmt.Sprintf("(%s)", variable.Name.Value)
		if !strings.Contains(strings.ToLower(resourceData.Definition.Value), block) {
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
