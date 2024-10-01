package provider

import (
	"testing"

	// "github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestGenerateConventionsString(t *testing.T) {
	data := ConventionDataSourceData{
		ID:         types.StringValue("test"),
		Definition: types.StringValue("(name)-(type)-(default)"),
		Variables: []Variable{
			{
				Name:      types.StringValue("name"),
				Default:   types.StringNull(),
				MaxLength: types.Int64Null(),
			},
			{
				Name:      types.StringValue("type"),
				Default:   types.StringNull(),
				MaxLength: types.Int64Value(int64(3)),
			},
			{
				Name:      types.StringValue("default"),
				Default:   types.StringValue("default"),
				MaxLength: types.Int64Null(),
			},
		},
		Convention: types.StringNull(),
	}

	diags := generateConventionsString(&data)
	assert.False(t, diags.HasError())
	expected := "{\"definition\":\"(name)-(type)-(default)\",\"variables\":[{\"name\":\"\\\"name\\\"\",\"default\":\"\\u003cnull\\u003e\",\"max_length\":\"\\u003cnull\\u003e\"},{\"name\":\"\\\"type\\\"\",\"default\":\"\\u003cnull\\u003e\",\"max_length\":\"3\"},{\"name\":\"\\\"default\\\"\",\"default\":\"\\\"default\\\"\",\"max_length\":\"\\u003cnull\\u003e\"}]}"
	assert.Equal(t, expected, data.Convention.ValueString())

}

func TestVariableConversion(t *testing.T) {
	var conversionVariables []ConventionVariable
	var newVariables []Variable
	variables := []Variable{
		{
			Name:      types.StringValue("name"),
			Default:   types.StringNull(),
			MaxLength: types.Int64Value(int64(8)),
		},
		{
			Name:      types.StringValue("type"),
			Default:   types.StringNull(),
			MaxLength: types.Int64Value(int64(3)),
		},
		{
			Name:      types.StringValue("default"),
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
			Default:   types.StringNull(),
			MaxLength: types.Int64Value(int64(8)),
		}),
		convertVariableToConventionVariable(Variable{
			Name:      types.StringValue("type"),
			Default:   types.StringNull(),
			MaxLength: types.Int64Value(int64(3)),
		}),
		convertVariableToConventionVariable(Variable{
			Name:      types.StringValue("default"),
			Default:   types.StringValue("default"),
			MaxLength: types.Int64Null(),
		}),
	}

	inputs := map[string]string{
		"name": "test",
		"type": "service",
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
	convention := Convention{
		Definition: "(name)-(type)-(default)",
		Variables: []ConventionVariable{
			convertVariableToConventionVariable(Variable{
				Name:      types.StringValue("name"),
				Default:   types.StringNull(),
				MaxLength: types.Int64Null(),
			}),
			convertVariableToConventionVariable(Variable{
				Name:      types.StringValue("type"),
				Default:   types.StringNull(),
				MaxLength: types.Int64Value(int64(4)),
			}),
			convertVariableToConventionVariable(Variable{
				Name:      types.StringValue("default"),
				Default:   types.StringValue("default"),
				MaxLength: types.Int64Null(),
			}),
		},
	}

	inputName := "foobar"
	inputType := "service"
	inputs := map[string]*string{
		"name": &inputName,
		"type": &inputType,
	}

	result, err := generateName(inputs, convention)
	assert.Equal(t, result, "foobar-serv-default")
	assert.Nil(t, err)
}
