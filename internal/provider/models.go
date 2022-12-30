package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// NameResourceData -
type NameResourceData struct {
	ID         types.String `tfsdk:"id"`
	Convention types.String `tfsdk:"convention"`
	Name       types.String `tfsdk:"name"`
	Inputs     types.Map    `tfsdk:"inputs"`
	Result     types.String `tfsdk:"result"`
}

// ConventionResourceData schema struct
type ConventionResourceData struct {
	ID         types.String `tfsdk:"id"`
	Definition types.String `tfsdk:"definition"`
	Variables  []Variable   `tfsdk:"variables"`
	Convention types.String `tfsdk:"convention"`
}

// Variable -
type Variable struct {
	Name      types.String `tfsdk:"name"`
	Default   types.String `tfsdk:"default"`
	Generated types.Bool   `tfsdk:"generated"`
	MaxLength types.Int64  `tfsdk:"max_length"`
}

// Convention contains the validated convention
type Convention struct {
	Definition string               `json:"definition"`
	Variables  []ConventionVariable `json:"variables"`
}

// ConventionVariable -
type ConventionVariable struct {
	Name      string `json:"name"`
	Default   string `json:"default"`
	Generated string `json:"generated"`
	MaxLength string `json:"max_length"`
}
