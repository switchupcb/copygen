package models

import "fmt"

// Field represents a field (or value) to be copied to.
type Field struct {
	Parent  Type         // The type that contains this field.
	Name    string       // The name of the field.
	Convert string       // The convert-function used to copy this field.
	From    *Field       // The field that this field will be copied from.
	To      *Field       // The field that this field will be copied to.
	Options FieldOptions // The custom options of a field.
}

type FieldOptions struct {
	Custom map[string]interface{} // The custom options of a field.
}

// Validate validates a field by ensuring it doesn't point to itself.
func (f Field) Validate() error {
	if &f == f.From {
		return fmt.Errorf("A field cannot point to itself.")
	}
	return nil
}
