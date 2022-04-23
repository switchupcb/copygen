// Package copygen contains the setup information for copygen generated code.
package copygen

type Placeholder bool

// Copygen defines the functions that will be generated.
type Copygen interface {
	Interface(IFC error) ifcHolder
	Func(F func() int) funcHolder

	Array(Arr [16]byte) Container
	Slice(S []string) Container
	Map(M map[string]bool) Container
	Chan(C chan int) Container

	Complex(Placeholder) Placeholder
}

// ifcHolder represents a type that holds an interface.
type ifcHolder struct {
	IFC ifc
}

// ifc represents an interface type (equivalent to `error`).
type ifc interface {
	Error() string
}

// funcHolder represents a type that holds a func.
type funcHolder struct {
	F func() int
}

// Container represents a type that holds container types.
type Container struct {
	Arr [16]byte
	S   []string
	M   map[string]bool
	C   chan int
}
