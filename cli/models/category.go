package models

// IsType returns whether the field is a type.
func (f *Field) IsType() bool {
	return f.Parent == nil
}

// basicMap contains a list of basic types.
var (
	basicMap = map[string]bool{
		"invalid":    true,
		"bool":       true,
		"int":        true,
		"int8":       true,
		"int16":      true,
		"int32":      true,
		"int64":      true,
		"uint":       true,
		"uint8":      true,
		"uint16":     true,
		"uint32":     true,
		"uint64":     true,
		"uintptr":    true,
		"float32":    true,
		"float64":    true,
		"complex64":  true,
		"complex128": true,
		"string":     true,
		"byte":       true,
		"rune":       true,
	}
)

// IsBasic determines whether the field is a basic type.
func (f *Field) IsBasic() bool {
	return basicMap[f.Definition]
}

// Pointer represents the char representation of a pointer.
const Pointer = '*'

// IsPointer returns whether the field is a pointer.
func (f *Field) IsPointer() bool {
	return len(f.Definition) >= 1 && f.Definition[1] == Pointer
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
	return f.Definition != "" && !(f.IsBasic() || f.IsPointer() || f.IsComposite() || f.IsFunc())
}
