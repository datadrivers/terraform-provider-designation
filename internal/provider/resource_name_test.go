package provider

import (
	"testing"

	// "github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestValidateInputs(t *testing.T) {
	variables := []Variable{
		{
			Name: types.String{Value: "not-generated"},
			Generated: types.Bool{
				Null:  true,
				Value: false},
			Default: types.String{
				Null: true,
			},
			MaxLength: types.Int64{
				Null: true,
			},
		},
		{
			Name:      types.String{Value: "generated"},
			Generated: types.Bool{Value: true},
			Default: types.String{
				Null: true,
			},
			MaxLength: types.Int64{
				Null: true,
			},
		},
		{
			Name: types.String{Value: "default"},
			Generated: types.Bool{
				Null: true,
			},
			Default: types.String{
				Null:  false,
				Value: "default",
			},
			MaxLength: types.Int64{
				Null: true,
			},
		},
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
