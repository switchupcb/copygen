// Package copygen contains the setup information for copygen generated code.
package copygen

import (
	"github.com/switchupcb/copygen/examples/_tests/multi/complex"
	"github.com/switchupcb/copygen/examples/_tests/multi/external"
)

// Placeholder represents a basic type..
type Placeholder bool

// Copygen defines the functions that will be generated.
type Copygen interface {
	NoMatchBasic(A Placeholder) (B Placeholder)
	NoMatchBasicExternal(A *Placeholder) (A external.Placeholder, B *external.Placeholder, C bool)
	Basic(bool) bool
	BasicSimple(Placeholder) Placeholder
	BasicPointer(Placeholder) *Placeholder
	BasicPointerMulti(A *Placeholder) (A *Placeholder, B *Placeholder, C string)
	BasicExternal(*external.Placeholder) external.Placeholder
	BasicExternalMulti(*external.Placeholder) (external.Placeholder, *external.Placeholder)

	NoMatchArraySimple([16]byte) Collection
	Array([16]byte) [16]byte
	ArraySimple(Arr [16]byte) *Collection
	ArrayExternal(external.Collection) *external.Collection
	ArrayComplex(Arr [16]map[byte]string) *complex.Collection
}

// Collection represents a type that holds collection field types.
type Collection struct {
	Arr [16]byte
	S   []string
	M   map[string]bool
	C   chan int
}
