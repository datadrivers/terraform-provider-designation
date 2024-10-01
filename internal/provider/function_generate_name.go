package provider

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the desired interfaces.
var _ function.Function = &GenerateNameFunction{}

// With the function.Function implementation
func NewGenerateNameFunction() function.Function {
	return &GenerateNameFunction{}
}

type GenerateNameFunction struct{}

func (f *GenerateNameFunction) Metadata(ctx context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "generate_name"
}
func (f *GenerateNameFunction) Definition(ctx context.Context, req function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:             "Generate a name for a convention and with specified parameters",
		MarkdownDescription: "Generate a name for a convention and with specified parameters",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:                "convention",
				MarkdownDescription: "The validated convention formated as a json string",
			},
			function.MapParameter{
				Name:                "inputs",
				MarkdownDescription: "Map of input values for variables in provider defined convention",
				ElementType:         types.StringType,
			},
		},
		Return: function.StringReturn{},
	}
}
func (f *GenerateNameFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var conventionString string
	var convention Convention
	var inputs map[string]*string
	// Read Terraform argument data into the variable
	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &conventionString, &inputs))

	if err := json.Unmarshal([]byte(conventionString), &convention); err != nil {
		resp.Error = function.ConcatFuncErrors(resp.Error, function.NewFuncError("Error reading convention: "+err.Error()))
		return
	}

	result, err := generateName(inputs, convention)
	if err != nil {
		resp.Error = function.ConcatFuncErrors(resp.Error, function.NewFuncError("Error generating name: "+err.Error()))
		return
	}

	// Set the result to the same data
	resp.Error = function.ConcatFuncErrors(resp.Error, resp.Result.Set(ctx, result))
}
