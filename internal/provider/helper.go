package provider

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// generateConventionsString generates the json string for the convention
func generateConventionsString(data *ConventionDataSourceData) diag.Diagnostics {
	var diags diag.Diagnostics

	var variables []ConventionVariable
	for _, variable := range data.Variables {
		variables = append(variables, convertVariableToConventionVariable(variable))
	}
	convention, err := json.Marshal(Convention{
		Definition: data.Definition.ValueString(),
		Variables:  variables,
	})
	if err != nil {
		diags.AddError(
			"Convention generation Error",
			err.Error(),
		)
		return diags
	}

	data.Convention = types.StringValue(string(convention))
	return diags
}

// validateConvention checks the configured convention to ensure that it can be used without errors
func validateConvention(data *ConventionDataSourceData) diag.Diagnostics {
	var diags diag.Diagnostics

	if !strings.Contains(strings.ToLower(data.Definition.ValueString()), "(name)") {
		diags.AddError(
			"Convention Validate Error",
			"The defined convention must include the block '(name)'.",
		)
	}
	var missingVariables []string

	for _, variable := range data.Variables {
		block := fmt.Sprintf("(%s)", variable.Name.ValueString())
		if !strings.Contains(strings.ToLower(data.Definition.ValueString()), block) {
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

func convertVariableToConventionVariable(input Variable) ConventionVariable {
	return ConventionVariable{
		Name:      input.Name.String(),
		Default:   input.Default.String(),
		MaxLength: input.MaxLength.String(),
	}
}

func convertConventionVariableToVariable(input ConventionVariable) Variable {
	var name types.String
	switch input.Name {
	case attr.NullValueString:
		name = types.StringNull()
	case attr.UnknownValueString:
		name = types.StringUnknown()
	default:
		name = types.StringValue(strings.Trim(input.Name, "\\\""))
	}

	var def types.String
	switch input.Default {
	case attr.NullValueString:
		def = types.StringNull()
	case attr.UnknownValueString:
		def = types.StringUnknown()
	default:
		def = types.StringValue(strings.Trim(input.Default, "\\\""))
	}

	var maxLength types.Int64
	switch input.MaxLength {
	case attr.NullValueString:
		maxLength = types.Int64Null()
	case attr.UnknownValueString:
		maxLength = types.Int64Unknown()
	default:
		n, _ := strconv.ParseInt(input.MaxLength, 10, 64)
		maxLength = types.Int64Value(n)
	}

	return Variable{
		Name:      name,
		Default:   def,
		MaxLength: maxLength,
	}
}

// validateInputs checks if all needed inputs are set
func validateInputs(inputs map[string]string, variables []ConventionVariable) diag.Diagnostics {
	var diags diag.Diagnostics

	missingInputs := []string{}
	for _, conventionVariable := range variables {
		variable := convertConventionVariableToVariable(conventionVariable)
		if variable.Default.IsNull() {
			if _, ok := inputs[variable.Name.ValueString()]; !ok {
				missingInputs = append(missingInputs, variable.Name.ValueString())
			}
		}

	}
	if len(missingInputs) > 0 {
		diags.AddError(
			"Convention Usage Error",
			fmt.Sprintf("All convention variables that are not generated or have a default must be present. Missing inputs: %s", strings.Join(missingInputs, ", ")),
		)
	}

	return diags
}

// generateNameWithSeperatedName generates the name with the name, inputs map and definition
func generateName(inputs map[string]*string, convention Convention) (string, error) {
	result := convention.Definition

	for _, conventionVariable := range convention.Variables {
		variable := convertConventionVariableToVariable(conventionVariable)
		block := fmt.Sprintf("(%s)", variable.Name.ValueString())

		replacement := ""
		length := 0

		if inputs[variable.Name.ValueString()] != nil {
			replacement = *inputs[variable.Name.ValueString()]
		}

		if !variable.MaxLength.IsNull() && variable.MaxLength.ValueInt64() > int64(0) {
			length = int(variable.MaxLength.ValueInt64())
		}

		if replacement == "" {
			replacement = variable.Default.ValueString()
		}

		if length > 0 && len(replacement) > length {
			replacement = replacement[0:length]
		}

		if replacement == "" {
			return "nil", errors.New(fmt.Sprintf("Missing value %q", variable.Name.ValueString()))
		}

		result = strings.Replace(result, block, replacement, -1)
	}

	return result, nil
}
