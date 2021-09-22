package models

// Function represents the properties of a generated function.
type Function struct {
	Name    string          // The name of the function.
	To      []Type          // The types to copy fields to.
	From    []Type          // The types to copy fields from.
	Options FunctionOptions // The custom options of a function.
}

// FuncOptions represent options for a Function.
type FunctionOptions struct {
	Custom map[string]interface{} // The custom options of a function.
}
