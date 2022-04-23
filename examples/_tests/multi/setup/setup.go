// Package copygen contains the setup information for copygen generated code.
package copygen

type Placeholder bool

// Copygen defines the functions that will be generated.
type Copygen interface {
	Interface(IFC error) ifcHolder
	Func(Placeholder) Placeholder

	/* Container Types */
	Array(Placeholder) Placeholder
	Slice(Placeholder) Placeholder
	Map(Placeholder) Placeholder
	Chan(Placeholder) Placeholder
}

// ifcHolder represents a type that holds an interface.
type ifcHolder struct {
	IFC ifc
}

// ifc represents an interface type (equivalent to `error`).
type ifc interface {
	Error() string
}
