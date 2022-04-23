package models

import "fmt"

// Type represents a field that isn't contained.
type Type struct {
	// Field represents field information for the type.
	Field *Field
}

// Name gets the name of the type field.
func (t Type) Name() string {
	return t.Field.FullName("")
}

func (t Type) String() string {
	return fmt.Sprintf("type %v", t.Field.FullName(""))
}
