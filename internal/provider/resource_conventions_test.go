package provider

import (
	"testing"

	// "github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestGenerateConventionsString(t *testing.T) {
	resourceData := ConventionResourceData{
		ID:         types.StringValue("test"),
		Definition: types.StringValue("(name)-(not-generated)-(default)-(generated)"),
		Variables: []Variable{
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
				MaxLength: types.Int64Value(int64(4)),
			},
			{
				Name:      types.StringValue("default"),
				Generated: types.BoolNull(),
				Default:   types.StringValue("default"),
				MaxLength: types.Int64Null(),
			},
		},
		Convention: types.StringNull(),
	}

	diags := generateConventionsString(&resourceData)
	assert.False(t, diags.HasError())
	expected := "{\"definition\":\"(name)-(not-generated)-(default)-(generated)\",\"variables\":[{\"name\":\"\\\"not-generated\\\"\",\"default\":\"\\u003cnull\\u003e\",\"generated\":\"false\",\"max_length\":\"\\u003cnull\\u003e\"},{\"name\":\"\\\"generated\\\"\",\"default\":\"\\u003cnull\\u003e\",\"generated\":\"true\",\"max_length\":\"4\"},{\"name\":\"\\\"default\\\"\",\"default\":\"\\\"default\\\"\",\"generated\":\"\\u003cnull\\u003e\",\"max_length\":\"\\u003cnull\\u003e\"}]}"
	assert.Equal(t, expected, resourceData.Convention.ValueString())

}
