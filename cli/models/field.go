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

	// Import represents the file that field was imported from.
	Import string

	// Package represents the package the field is defined in (i.e `log` in `log.Logger`).
	Package string

	// Container represents the container that contains the field's definition (in string format).
	//
	// This can be any of the following examples or `nil``.
	// pointer(s): **
	// array: [5]
	// slice: []
	// map: map
	// chan: chan
	Container string

	// Name represents the name of the field (i.e `ID` in `ID int`).
	Name string

	// Definition represents the type definition of the field (i.e `int` in `ID int`, `Logger` in `log.Logger`).
	Definition string

	// Collection represents the type of collection that this field represents.
	//
	// a "struct" collection contain subfields with any other type.
	// an "interface" collection contains `func` subfields.
	Collection string

	// The tags defined in a struct field (i.e `json:"tag,omitempty"`)
	// map[tag]map[name][]options (i.e map[json]map[tag]["omitempty"])
	Tags map[string]map[string][]string

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

	// The tag to map this field with, if any.
	Tag string

	// The level at which sub-fields are discovered.
	Depth int

	// Whether the field should be explicitly automatched.
	Automatch bool

	// Whether the field should be deepcopied.
	Deepcopy bool
}

// isStruct returns whether the field is a struct.
func (f *Field) IsStruct() bool {
	return f.Collection == "struct"
}

// isInterface returns whether the field is an interface.
func (f *Field) IsInterface() bool {
	return f.Collection == "interface"
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

// FullDefinition returns the full definition of a field including its package.
func (f *Field) FullDefinition() string {
	if f.Package != "" {
		return f.Package + "." + f.Definition
	}

	return f.Definition
}

// FullNameWithoutContainer returns the full name of a field including its parents
// without the container (i.e domain.Account.User.ID).
func (f *Field) FullNameWithoutContainer(name string) string {
	if !f.IsType() {
		// names are added in reverse.
		if name == "" {
			// reference the field (i.e `ID`).
			name = f.Name
		} else {
			// prepend the field (i.e `User` + `.` + `ID`).
			name = f.Name + "." + name
		}

		return f.Parent.FullNameWithoutContainer(name)
	}

	if name != "" {
		name = "." + name
	}

	if f.Package != "" {
		return f.Package + "." + f.Definition + name
	}

	return f.Definition + name
}

// FullName returns the full name of a field including its parents (i.e *domain.Account.User.ID).
func (f *Field) FullName(name string) string {
	return f.Container + f.FullNameWithoutContainer("")
}

// FullVariableName gets the full variable name of a field (i.e tA.User.UserID).
func (f *Field) FullVariableName(name string) string {
	if !f.IsType() {
		return f.Parent.FullVariableName(f.VariableName + name)
	}

	return f.VariableName + name
}

// AllFields gets all the fields in the scope of a field (including itself).
func (f *Field) AllFields(fields []*Field, cyclic map[*Field]bool) []*Field {
	if cyclic == nil {
		cyclic = make(map[*Field]bool)
	}

	fields = append(fields, f)
	cyclic[f] = true
	for _, subfield := range f.Fields {
		if !cyclic[subfield] {
			fields = subfield.AllFields(fields, cyclic)
		}
	}

	return fields
}

func (f *Field) String() string {
	direction := "Unpointed"
	if f.From != nil {
		direction = "To"
	}

	if f.To != nil {
		switch direction {
		case "To":
			direction = "To and From"
		case "Unpointed":
			direction = "From"
		}
	}

	var parent string
	if f.Parent != nil {
		parent = f.Parent.FullName("")
	}

	return fmt.Sprintf("%v Field %q of Definition %q Fields[%v]: Parent %q", direction, f.FullName(""), f.Definition, len(f.Fields), parent)
}
