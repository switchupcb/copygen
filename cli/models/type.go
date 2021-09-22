package models

// Type represents a type (or struct) that will be copied to/from another type.
type Type struct {
	Filepath string      // The location of the type in an existing codebase.
	Name     string      // The name of the type in the provided file.
	Fields   []Field     // The fields of the type.
	Options  TypeOptions // The type options used for this type.
}

// TypeOptions represent valid options for a Type.
type TypeOptions struct {
	Pointer  bool // Whether this type should be used with a pointer.
	Deepcopy bool // Whether this type should be deepcopied (and returned).
}
