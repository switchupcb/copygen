package models

// Field represents a field (or value) to be copied to.
type Field struct {
	Name    string // The name of the field.
	To      *Field // The field that this field will be copied to.
	Convert string // The convert-function used to copy this field.
}
