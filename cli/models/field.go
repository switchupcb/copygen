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

	// Name represents the name of the field (i.e `ID` in `ID int`).
	Name string

	// Definition represents the type definition of the field (i.e `int` in `ID int`, `Logger` in `log.Logger`).
	Definition string

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

// Deepcopy returns a new field with copied properties (excluding Parent, To, and From fields).
func (f *Field) Deepcopy(cyclic map[*Field]bool) *Field {
	copied := &Field{
		VariableName: f.VariableName,
		Import:       f.Import,
		Package:      f.Package,
		Name:         f.Name,
		Definition:   f.Definition,
		Options: FieldOptions{
			Convert:   f.Options.Convert,
			Map:       f.Options.Map,
			Tag:       f.Options.Tag,
			Depth:     f.Options.Depth,
			Automatch: f.Options.Automatch,
			Deepcopy:  f.Options.Deepcopy,
		},
	}

	copied.Tags = make(map[string]map[string][]string, len(f.Tags))
	for k1, mapval := range f.Tags {
		copied.Tags[k1] = make(map[string][]string, len(mapval))
		for k2, sliceval := range f.Tags[k1] {
			copied.Tags[k1][k2] = make([]string, len(sliceval))
			copy(f.Tags[k1][k2], copied.Tags[k1][k2])
		}
	}

	// setup the cache
	if cyclic == nil {
		cyclic = make(map[*Field]bool)
	}

	// copy the subfields.
	cyclic[f] = true
	copied.Fields = make([]*Field, len(f.Fields))
	for i, sf := range f.Fields {
		if cyclic[sf] {
			copied.Fields[i] = sf
			continue
		}

		copied.Fields[i] = sf.Deepcopy(cyclic)
		copied.Fields[i].Parent = copied
	}
	return copied
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

// FullVariableName returns the full variable name of a field (i.e tA.User.UserID).
func (f *Field) FullVariableName(name string) string {
	if !f.IsType() {
		return f.Parent.FullVariableName(f.VariableName + name)
	}

	return f.VariableName + name
}

// FullDefinition returns the full definition of a field including its package
// without its pointer(s) (i.e domain.Account).
func (f *Field) FullDefinitionWithoutPointer() string {
	i := 0
	for i < len(f.Definition) && f.Definition[i] == Pointer {
		i++
	}

	if f.Package == "" {
		return f.Definition[i:]
	}

	return f.Package + "." + f.Definition[i:]
}

// FullDefinition returns the full definition of a field including its package.
func (f *Field) FullDefinition() string {
	if f.Package == "" {
		return f.Definition
	}

	i := 0
	for i < len(f.Definition) && f.Definition[i] == Pointer {
		i++
	}

	return f.Definition[:i] + f.Package + "." + f.Definition[i:]
}

// FullNameWithoutPointer returns the full name of a field including its parents
// without the pointer (i.e domain.Account.User.ID).
func (f *Field) FullNameWithoutPointer(name string) string {
	if !f.IsType() {
		// names are added in reverse.
		if name == "" {
			// reference the field (i.e `ID`).
			name = f.Name
		} else {
			// prepend the field (i.e `User` + `.` + `ID`).
			name = f.Name + "." + name
		}

		return f.Parent.FullNameWithoutPointer(name)
	}

	if name != "" {
		name = "." + name
	}

	return f.FullDefinitionWithoutPointer() + name
}

// FullName returns the full name of a field including its parents (i.e *domain.Account.User.ID).
func (f *Field) FullName() string {
	i := 0
	for i < len(f.Definition) && f.Definition[i] == Pointer {
		i++
	}

	return f.Definition[:i] + f.FullNameWithoutPointer("")
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

	var name string
	if f.Name != "" {
		name = f.Name + " "
	}

	var parent string
	if f.Parent != nil {
		parent = f.Parent.FullName()
	}

	return fmt.Sprintf("%v Field %v%q of Definition %q Fields[%v]: Parent %q", direction, name, f.FullName(), f.FullDefinition(), len(f.Fields), parent)
}
