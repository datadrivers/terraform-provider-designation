package provider

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"k8s.io/apimachinery/pkg/util/rand"
)

// validateInputs checks if all needed inputs are set
func validateInputs(inputs map[string]string, variables []Variable) diag.Diagnostics {
	var diags diag.Diagnostics

	missingInputs := []string{}
	for _, variable := range variables {
		if variable.Default.IsNull() && (variable.Generated.IsNull() || !variable.Generated.Value) {
			if _, ok := inputs[variable.Name.Value]; !ok && variable.Name.Value != "name" {
				missingInputs = append(missingInputs, variable.Name.Value)
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

	for _, variable := range convention.Variables {
		block := fmt.Sprintf("(%s)", variable.Name.Value)
		replacement := inputs[variable.Name.Value]
		if variable.Name.Value == "name" {
			replacement = name
			variableNameConfigured = true
		}
		length := 0

		if !variable.MaxLength.IsNull() && variable.MaxLength.Value > int64(0) {
			length = int(variable.MaxLength.Value)
		}

		if !variable.Generated.IsNull() && variable.Generated.Value {
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

	if !variableNameConfigured {
		result = strings.Replace(result, "(name)", name, -1)
	}

	return types.String{Value: result}, diags
}
