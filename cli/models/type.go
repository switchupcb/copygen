package models

import "fmt"

// Type represents a field that isn't contained.
type Type struct {
	Field   Field       // The field information for the type.
	Options TypeOptions // The type options used for the type.
}

// TypeOptions represent options for a Type.
type TypeOptions struct {
	Alloc bool // Whether a new type will be allocated.
}

// isStruct returns whether the type is a struct.
func (t Type) isStruct() bool {
	return t.Field.Definition == "struct"
}

// isInterface returns whether the type is an interface.
func (t Type) isInterface() bool {
	return t.Field.Definition == "interface"
}

func (t Type) String() string {
	return fmt.Sprintf("type %v.%v %v", t.Field.Package, t.Field.Name, t.Field.Definition)
}
