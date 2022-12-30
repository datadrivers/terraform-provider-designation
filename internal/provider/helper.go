package provider

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"k8s.io/apimachinery/pkg/util/rand"
)

func convertVariableToConventionVariable(input Variable) ConventionVariable {
	return ConventionVariable{
		Name:      input.Name.String(),
		Generated: input.Generated.String(),
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

	var generated types.Bool
	switch input.Generated {
	case attr.NullValueString:
		generated = types.BoolNull()
	case attr.UnknownValueString:
		generated = types.BoolUnknown()
	default:
		boolValue, _ := strconv.ParseBool(input.Generated)
		generated = types.BoolValue(boolValue)
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
		Generated: generated,
		MaxLength: maxLength,
	}
}

// validateInputs checks if all needed inputs are set
func validateInputs(inputs map[string]string, variables []ConventionVariable) diag.Diagnostics {
	var diags diag.Diagnostics

	missingInputs := []string{}
	for _, conventionVariable := range variables {
		variable := convertConventionVariableToVariable(conventionVariable)
		if variable.Default.IsNull() && (variable.Generated.IsNull() || !variable.Generated.ValueBool()) {
			if _, ok := inputs[variable.Name.ValueString()]; !ok && variable.Name.ValueString() != "name" {
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

// generateName generates the name with the inputs, variables and definition
func generateName(name string, inputs map[string]string, convention Convention) (types.String, diag.Diagnostics) {
	var diags diag.Diagnostics

	result := convention.Definition
	variableNameConfigured := false

	for _, conventionVariable := range convention.Variables {
		variable := convertConventionVariableToVariable(conventionVariable)
		block := fmt.Sprintf("(%s)", variable.Name.ValueString())
		replacement := inputs[variable.Name.ValueString()]
		if variable.Name.ValueString() == "name" {
			replacement = name
			variableNameConfigured = true
		}
		length := 0

		if !variable.MaxLength.IsNull() && variable.MaxLength.ValueInt64() > int64(0) {
			length = int(variable.MaxLength.ValueInt64())
		}

		if !variable.Generated.IsNull() && variable.Generated.ValueBool() {
			replacement = rand.String(length)
		}

		if replacement == "" {
			replacement = variable.Default.ValueString()
		}

		if length > 0 && len(replacement) > length {
			replacement = replacement[0:length]
		}
		result = strings.Replace(result, block, replacement, -1)
	}

	if !variableNameConfigured {
		result = strings.Replace(result, "(name)", name, -1)
	}

	return types.StringValue(result), diags
}
