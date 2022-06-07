package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// NameResource -
type nameResourceData struct {
	ID         types.String `tfsdk:"id"`
	Convention types.String `tfsdk:"convention"`
	Name       types.String `tfsdk:"name"`
	Inputs     types.Map    `tfsdk:"inputs"`
	Result     types.String `tfsdk:"result"`
}

// Provider schema struct
type conventionResourceData struct {
	ID         types.String `tfsdk:"id"`
	Definition types.String `tfsdk:"definition"`
	Variables  []Variable   `tfsdk:"variables"`
	Convention types.String `tfsdk:"convention"`
}

// Convention contains the validated convention
type Convention struct {
	Definition string     `json:"definition"`
	Variables  []Variable `json:"variables"`
}

// Variable -
type Variable struct {
	Name      types.String `tfsdk:"name"`
	Default   types.String `tfsdk:"default"`
	Generated types.Bool   `tfsdk:"generated"`
	MaxLength types.Int64  `tfsdk:"max_length"`
}
