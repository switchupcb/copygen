package models

// Field represents a field (or value) to be copied to.
type Field struct {
	Parent     Type         // The type that contains this field.
	Name       string       // The name of the field.
	Definition string       // The type definition of the field.
	Convert    string       // The convert-function used to copy this field.
	Fields     []Field      // The fields of the field.
	From       *Field       // The field that this field will be copied from.
	To         *Field       // The field that this field will be copied to.
	Options    FieldOptions // The custom options of a field.
}

type FieldOptions struct {
	Depth    string                 // Whether the field should be copied recursively (in-depth).
	Deepcopy string                 // Whether the field should be deepcopied.
	Custom   map[string]interface{} // The custom options of a field.
}
