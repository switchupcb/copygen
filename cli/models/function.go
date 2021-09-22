package models

// Function represents the properties of a generated function.
type Function struct {
	Name       string          // The name of the function.
	Parameters map[string]Type // The parameters (variable name and type) of the function.
	To         []Type          // The types to copy fields to.
	From       []Type          // the types to copy fields from.
	Options    FunctionOptions // The function options for this type.
}

// FunctionOptions represent valid options for a Function.
type FunctionOptions struct {
	Error bool // Whether the function returns an error.
}
