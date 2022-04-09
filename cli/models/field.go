// Package models defines the domain models that model field relations and manipulation.
package models

import (
	"fmt"
)

// Field represents a field to be copied to/from.
// A field's struct properties are set in the parser unless its stated otherwise.
type Field struct {
	// VariableName represents name that is used to assign the field.
	//
	// This value will always be unique in the context of the application.
	// TypeField variable names do not contain '.' (i.e 'tA' in 'tA.UserID').
	// Field variable names are defined by their specifier (i.e '.UserID' in 'domain.Account.UserID').
	VariableName string

	// The package the field is defined in (i.e `log` in `log.Logger`).
	Package string

	// The name of the field (i.e `ID` in `ID int`).
	Name string

	// The type definition of the field (i.e `int` in `ID int`, `Logger` in `log.Logger`).
	Definition string

	// The container type that contains the field's definition in string format.
	//
	// This can be any of the following examples or `nil``.
	// pointer(s): **
	// array: [5]
	// slice: []
	// map: map
	// chan: chan
	Container string

	// The tag defined in a struct field (i.e `json:"tag"`)
	Tag string

	// The file the field was imported from.
	Import string

	// The type or field that contains this field.
	Parent *Field

	// The field that this field will be copied from (or nil).
	// Set in the matcher.
	From *Field

	// The field that this field will be copied to (or nil).
	// Set in the matcher.
	To *Field

	// The fields of this field.
	Fields []*Field

	// The custom options of a field.
	Options FieldOptions
}

// FieldOptions represent options for a Field.
type FieldOptions struct {
	// The function the field is converted with (as a parameter).
	Convert string

	// The field to map this field to, if any.
	Map string

	// The level at which sub-fields are discovered.
	Depth int

	// Whether the field should be deepcopied.
	Deepcopy bool
}

// isStruct returns whether the field is a struct.
func (f *Field) IsStruct() bool {
	return f.Definition == "struct"
}

// isInterface returns whether the field is an interface.
func (f *Field) IsInterface() bool {
	return f.Definition == "interface"
}

// IsNoContainer returns whether the field has no container.
func (f *Field) IsNoContainer() bool {
	return f.Container == ""
}

// IsPointer returns whether the field is a pointer of a type definition.
func (f *Field) IsPointer() bool {
	return f.Container != "" && f.Container[0] == '*'
}

// IsArray returns whether the field is an array.
// assumes the caller is checking a valid container.
func (f *Field) IsArray() bool {
	return len(f.Container) >= 3 && f.Container[0] == '['
}

// IsSlice returns whether the field is a slice.
func (f *Field) IsSlice() bool {
	return len(f.Container) >= 2 && f.Container[0] == '[' && f.Container[1] == ']'
}

// IsMap returns whether the field is a map.
// assumes the caller is checking a valid container.
func (f *Field) IsMap() bool {
	return len(f.Container) >= 3 && f.Container[0] == 'm'
}

// IsMap returns whether the field is a chan.
// assumes the caller is checking a valid container.
func (f *Field) IsChan() bool {
	return len(f.Container) >= 4 && f.Container[0] == 'c'
}

// IsType returns whether the field is a type.
func (f *Field) IsType() bool {
	return f.Parent == nil
}

// FullName gets the full name of a field including its parents (i.e *domain.Account.User.ID).
func (f *Field) FullName(name string) string {
	if !f.IsType() {
		// add names in reverse order
		if name == "" {
			name = f.Name
		} else {
			name = f.Name + "." + name
		}

		return f.Parent.FullName(name)
	}

	if name != "" {
		name = "." + name
	}

	return fmt.Sprintf("%v%v.%v%v", f.Container, f.Package, f.Name, name)
}

// FullNameWithoutContainer gets the full name of a field including its parents
// without the container (i.e domain.Account.User.ID).
func (f *Field) FullNameWithoutContainer(name string) string {
	if !f.IsType() {
		// add names in reverse order
		if name == "" {
			name = f.Name
		} else {
			name = f.Name + "." + name
		}

		return f.Parent.FullName(name)[len(f.Parent.Container):]
	}

	if name != "" {
		name = "." + name
	}

	return fmt.Sprintf("%v.%v%v", f.Package, f.Name, name)
}

// FullVariableName gets the full variable name of a field (i.e tA.User.UserID).
func (f *Field) FullVariableName(name string) string {
	if !f.IsType() {
		return f.Parent.FullVariableName(f.VariableName + name)
	}

	return f.VariableName + name
}

// AllFields gets all the fields in the scope of a field (including itself).
func (f *Field) AllFields(fields []*Field) []*Field {
	fields = append(fields, f)

	if len(f.Fields) != 0 {
		for i := 0; i < len(f.Fields); i++ {
			fields = f.Fields[i].AllFields(fields)
		}
	}

	return fields
}

func (f *Field) String() string {
	var direction, parent string

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

	if f.Parent != nil {
		parent = f.Parent.FullName("")
	}

	return fmt.Sprintf("%v Field %q of Definition %q: Parent %q Fields[%v]", direction, f.FullName(""), f.Definition, parent, len(f.Fields))
}
