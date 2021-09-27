package models

import "fmt"

// Field represents a field to be copied to/from.
type Field struct {
	Parent     Type         // The type that contains this field.
	Name       string       // The name of the field.
	Definition string       // The type definition of the field.
	Convert    string       // The convert-function used to copy the field.
	Fields     []*Field     // The fields of the field.
	From       *Field       // The field that the field will be copied from.
	To         *Field       // The field that the field will be copied to.
	Options    FieldOptions // The custom options of a field.
}

type FieldOptions struct {
	Deepcopy string                 // Whether the field should be deepcopied.
	Custom   map[string]interface{} // The custom options of a field.
}

func (f Field) String() string {
	var direction string
	if f.From != nil {
		direction = "To"
	}
	if f.To != nil {
		if direction != "" {
			direction += " and "
		}
		direction += "From"
	}
	if direction == "" {
		direction = "Unpointed"
	}

	name := f.Name
	if name == "" {
		name = "\"\""
	}

	convert := f.Convert
	if convert != "" {
		convert = " (Convert " + f.Convert + ")"
	}

	definition := f.Definition
	if definition == "" {
		definition = "\"\""
	}
	return fmt.Sprintf("%v Field %v of Definition %v%v: Parent %p Fields %v", direction, name, definition, convert, &f.Parent, f.Fields)
}
