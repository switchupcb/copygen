package models

// Type represents a type that contains fields to be copied to/from.
type Type struct {
	Package      string      // The package the type is defined in.
	Name         string      // The name of the type in the provided file.
	VariableName string      // The variable name the type is assigned.
	Fields       []*Field    // The fields of the type.
	Options      TypeOptions // The type options used for the type.
}

// TypeOptions represent options for a Type.
type TypeOptions struct {
	Import   string                 // The import path for the type.
	Pointer  bool                   // Whether the type should be used with a pointer.
	Depth    int                    // The level fields should be copied to/from.
	Deepcopy string                 // Whether the type should be deepcopied.
	Custom   map[string]interface{} // The custom options of a function.
}
