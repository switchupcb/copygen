package models

// IsType returns whether the field is a type.
func (f *Field) IsType() bool {
	return f.Parent == nil
}

// Pointer represents the string representation of a pointer.
const Pointer = "*"

// UsesPointer returns whether the field uses a pointer.
func (f *Field) UsesPointer() bool {
	return f.Pointer == Pointer
}

// Collection refers to a category of types which indicate that
// a field's definition collects multiple fields (i.e `map[string]bool`).
const (
	CollectionPointer   = "*"
	CollectionSlice     = "[]"
	CollectionMap       = "map"
	CollectionChan      = "chan"
	CollectionFunc      = "func"
	CollectionInterface = "interface"
)

// IsPointer returns whether the field is a pointer.
func (f *Field) IsPointer() bool {
	return len(f.Definition) >= 1 && f.Definition[0:1] == CollectionPointer
}

// IsArray returns whether the field is an array.
func (f *Field) IsArray() bool {
	return len(f.Definition) >= 3 && f.Definition[0] == '[' && ('0' <= f.Definition[1] && f.Definition[1] <= '9')
}

// IsSlice returns whether the field is a slice.
func (f *Field) IsSlice() bool {
	return len(f.Definition) >= 2 && f.Definition[:2] == CollectionSlice
}

// IsMap returns whether the field is a map.
func (f *Field) IsMap() bool {
	return len(f.Definition) >= 3 && f.Definition[:3] == CollectionMap
}

// IsMap returns whether the field is a chan.
func (f *Field) IsChan() bool {
	return len(f.Definition) >= 4 && f.Definition[:4] == CollectionChan
}

// IsComposite returns whether the field is a composite type: array, slice, map, chan.
func (f *Field) IsComposite() bool {
	return f.IsArray() || f.IsSlice() || f.IsMap() || f.IsChan()
}

// IsFunc returns whether the field is a function.
func (f *Field) IsFunc() bool {
	return len(f.Definition) >= 4 && f.Definition[:4] == CollectionFunc
}

// IsInterface returns whether the field is an interface.
func (f *Field) IsInterface() bool {
	return len(f.Definition) >= 9 && f.Definition[:9] == CollectionInterface
}

// IsCollection returns whether the field is a collection.
func (f *Field) IsCollection() bool {
	return f.IsPointer() || f.IsComposite() || f.IsFunc() || f.IsInterface()
}

// IsAlias determines whether the field is a type alias.
func (f *Field) IsAlias() bool {
	return f.Definition != "" && !f.IsCollection()
}
