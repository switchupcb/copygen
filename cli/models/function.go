package models

// Function represents the properties of a generated function.
type Function struct {
	Name    string          // The name of the function.
	From    []Type          // The types to copy fields from.
	To      []Type          // The types to copy fields to.
	Options FunctionOptions // The custom options of a function.
}

// FunctionOptions represent options for a Function.
type FunctionOptions struct {
	Custom map[string][]string // The custom options of a function (map[option]values).
}
