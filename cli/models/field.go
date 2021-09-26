package models

// Field represents a field to be copied to/from.
type Field struct {
	Parent     Type         // The type that contains this field.
	Name       string       // The name of the field.
	Definition string       // The type definition of the field.
	Convert    string       // The convert-function used to copy the field.
	Fields     []Field      // The fields of the field.
	From       *Field       // The field that the field will be copied from.
	To         *Field       // The field that the field will be copied to.
	Options    FieldOptions // The custom options of a field.
}

type FieldOptions struct {
	Deepcopy string                 // Whether the field should be deepcopied.
	Custom   map[string]interface{} // The custom options of a field.
}
