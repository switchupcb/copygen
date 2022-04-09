package models

import "fmt"

// Type represents a field that isn't contained.
type Type struct {
	Field *Field // The field information for the type.
}

// ParameterName gets the parameter name of the type.
func (t Type) ParameterName() string {
	return t.Field.Container + t.Field.Package + "." + t.Field.Name
}

func (t Type) String() string {
	return fmt.Sprintf("type %v", t.Field.FullName(""))
}
