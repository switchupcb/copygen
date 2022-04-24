package models

// Container refers to a category of types which indicate that
// a field contains multiple fields.
const (
	ContainerStruct    = "struct"
	ContainerInterface = "interface"
)

// IsContainer returns whether the field has a container.
func (f *Field) IsContainer() bool {
	return f.Container != ""
}

// isStruct returns whether the field is a struct.
func (f *Field) IsStruct() bool {
	return f.Container == ContainerStruct
}

// isInterface returns whether the field is an interface.
func (f *Field) IsInterface() bool {
	return f.Container == ContainerInterface
}

// Collection refers to a category of types which indicate that
// a field's definition collects multiple fields (i.e `map[string]bool`).
const (
	CollectionPointer = "*"
	CollectionSlice   = "[]"
	CollectionMap     = "map"
	CollectionChan    = "chan"
	CollectionFunc    = "func"
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

// IsCollection returns whether the field is a collection.
func (f *Field) IsCollection() bool {
	return f.IsPointer() || f.IsComposite() || f.IsFunc()
}
