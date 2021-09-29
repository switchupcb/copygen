// Package models defines the domain models that model field relations and manipulation.
package models

import "fmt"

// Field represents a field to be copied to/from.
type Field struct {
	// The variable name the field is assigned for assignment.
	// This value will always be unique in the context of the application.
	// Type variable names do not contain '.' (i.e 'tA' in 'tA.UserID')
	// Field variable names are defined by their specifier (i.e '.UserID' in 'domain.Account.UserID').
	VariableName string

	// The package the field is defined in.
	Package string

	// The name of the field (i.e ID in `ID int`).
	Name string

	// The type definition of the field (i.e int in `ID int`, struct, or interface).
	Definition string

	// The type or field that contains this field.
	Parent *Field

	// The field that this field will be copied from (or nil).
	From *Field

	// The field that this field will be copied to (or nil).
	To *Field

	// The fields of this field.
	Fields []*Field

	// The custom options of a field.
	Options FieldOptions
}

// FieldOptions represent options for a Field.
type FieldOptions struct {
	Depth    int    // The level at which sub-fields are discovered.
	Deepcopy bool   // Whether the field should be deepcopied.
	Convert  string // The function the field is converted with (as a parameter).
}

// IsType returns whether the field is a type.
func (f Field) IsType() bool {
	return f.Parent == nil
}

// FullName gets the full name of a field including its parents (i.e domain.Account.User.ID).
func (f Field) FullName(name string) string {
	if !f.IsType() {
		// add names in reverse order
		if name == "" {
			name = "." + f.Name
		} else {
			name = f.Name + "." + name
		}
		f.Parent.FullName(name)
	}
	return fmt.Sprintf("%v.%v.%v", f.Package, f.Name, name)
}

// FullVariableName gets the full variable name of a field (i.e tA.User.UserID)
func (f Field) FullVariableName(name string) string {
	if !f.IsType() {
		// add names in reverse order
		if name == "" {
			name = f.VariableName
		} else {
			name = f.VariableName + name
		}
		f.Parent.FullName(name)
	}
	return f.VariableName + name
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
	var parent string
	if f.Parent != nil {
		parent = f.Parent.FullName("")
	}
	return fmt.Sprintf("%v Field %q of Definition %q: Parent %q Fields[%v]", direction, f.FullName(""), f.Definition, parent, len(f.Fields))
}
