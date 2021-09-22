package models

// Type represents a type (or struct) that will be copied to/from another type.
type Type struct {
	Name         string      // The name of the type in the provided file.
	VariableName string      // The variable name the type is assigned.
	Package      string      // The package the type is defined in.
	Fields       []Field     // The fields of the type.
	Options      TypeOptions // The type options used for this type.
}

// TypeOptions represent options for a Type.
type TypeOptions struct {
	Pointer bool                   // Whether this type should be used with a pointer.
	Custom  map[string]interface{} // The custom options of a function.
}
