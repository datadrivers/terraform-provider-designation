package provider

import (
	"testing"

	// "github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestVariableConversion(t *testing.T) {
	var conversionVariables []ConventionVariable
	var newVariables []Variable
	variables := []Variable{
		{
			Name:      types.StringValue("name"),
			Generated: types.BoolNull(),
			Default:   types.StringNull(),
			MaxLength: types.Int64Value(int64(8)),
		},
		{
			Name:      types.StringValue("not-generated"),
			Generated: types.BoolValue(false),
			Default:   types.StringNull(),
			MaxLength: types.Int64Null(),
		},
		{
			Name:      types.StringValue("generated"),
			Generated: types.BoolValue(true),
			Default:   types.StringNull(),
			MaxLength: types.Int64Null(),
		},
		{
			Name:      types.StringValue("default"),
			Generated: types.BoolNull(),
			Default:   types.StringValue("default"),
			MaxLength: types.Int64Null(),
		},
	}

	for _, v := range variables {
		conversionVariables = append(conversionVariables, convertVariableToConventionVariable(v))
	}
	for _, v := range conversionVariables {
		newVariables = append(newVariables, convertConventionVariableToVariable(v))
	}
	assert.Equal(t, variables, newVariables)
}
func TestValidateInputs(t *testing.T) {
	variables := []ConventionVariable{
		convertVariableToConventionVariable(Variable{
			Name:      types.StringValue("name"),
			Generated: types.BoolNull(),
			Default:   types.StringNull(),
			MaxLength: types.Int64Value(int64(8)),
		}),
		convertVariableToConventionVariable(Variable{
			Name:      types.StringValue("not-generated"),
			Generated: types.BoolValue(false),
			Default:   types.StringNull(),
			MaxLength: types.Int64Null(),
		}),
		convertVariableToConventionVariable(Variable{
			Name:      types.StringValue("generated"),
			Generated: types.BoolValue(true),
			Default:   types.StringNull(),
			MaxLength: types.Int64Null(),
		}),
		convertVariableToConventionVariable(Variable{
			Name:      types.StringValue("default"),
			Generated: types.BoolNull(),
			Default:   types.StringValue("default"),
			MaxLength: types.Int64Null(),
		}),
	}

	inputs := map[string]string{
		"not-generated": "test",
	}

	diags := validateInputs(inputs, variables)
	assert.False(t, diags.HasError())

	inputs2 := map[string]string{
		"generated": "test",
	}
	diags2 := validateInputs(inputs2, variables)

	assert.True(t, diags2.HasError())
}

func TestGenerateName(t *testing.T) {
	name := "foobar"
	convention := Convention{
		Definition: "(name)-(not-generated)-(default)-(generated)",
		Variables: []ConventionVariable{
			convertVariableToConventionVariable(Variable{
				Name:      types.StringValue("not-generated"),
				Generated: types.BoolValue(false),
				Default:   types.StringNull(),
				MaxLength: types.Int64Null(),
			}),
			convertVariableToConventionVariable(Variable{
				Name:      types.StringValue("generated"),
				Generated: types.BoolValue(true),
				Default:   types.StringNull(),
				MaxLength: types.Int64Value(int64(4)),
			}),
			convertVariableToConventionVariable(Variable{
				Name:      types.StringValue("default"),
				Generated: types.BoolNull(),
				Default:   types.StringValue("default"),
				MaxLength: types.Int64Null(),
			}),
		},
	}

	inputs := map[string]string{
		"not-generated": "bar",
	}

	result, diags := generateName(name, inputs, convention)
	assert.False(t, diags.HasError())
	assert.Contains(t, result.ValueString(), "foobar-bar-default-")
	assert.Len(t, result.ValueString(), 23)
}

func TestGenerateNameWithConfiguredNameVariable(t *testing.T) {
	name := "foobar"
	convention := Convention{
		Definition: "(name)-(not-generated)-(default)-(generated)",
		Variables: []ConventionVariable{
			convertVariableToConventionVariable(Variable{
				Name:      types.StringValue("name"),
				Generated: types.BoolNull(),
				Default:   types.StringNull(),
				MaxLength: types.Int64Value(int64(3)),
			}),
			convertVariableToConventionVariable(Variable{
				Name:      types.StringValue("not-generated"),
				Generated: types.BoolValue(false),
				Default:   types.StringNull(),
				MaxLength: types.Int64Null(),
			}),
			convertVariableToConventionVariable(Variable{
				Name:      types.StringValue("generated"),
				Generated: types.BoolValue(true),
				Default:   types.StringNull(),
				MaxLength: types.Int64Value(int64(4)),
			}),
			convertVariableToConventionVariable(Variable{
				Name:      types.StringValue("default"),
				Generated: types.BoolNull(),
				Default:   types.StringValue("default"),
				MaxLength: types.Int64Null(),
			}),
		},
	}

	inputs := map[string]string{
		"not-generated": "bar",
	}

	result, diags := generateName(name, inputs, convention)
	assert.False(t, diags.HasError())
	assert.Contains(t, result.ValueString(), "foo-bar-default-")
	assert.Len(t, result.ValueString(), 20)
}
