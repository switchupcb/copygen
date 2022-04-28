// Package copygen contains the setup information for copygen generated code.
package copygen

import (
	"github.com/switchupcb/copygen/examples/_tests/multi/complex"
	"github.com/switchupcb/copygen/examples/_tests/multi/external"
)

// Placeholder represents a basic type.
type Placeholder bool

// Copygen defines the functions that will be generated.
type Copygen interface {
	NoMatchBasic(A Placeholder) (B Placeholder)
	NoMatchBasicExternal(A *Placeholder) (A external.Placeholder, B *external.Placeholder, C bool)
	Basic(bool) bool
	BasicPointer(bool) *bool
	BasicDoublePointer(*bool) **bool // Fail Pointer Semantics
	BasicSimple(Placeholder) Placeholder
	BasicSimplePointer(Placeholder) *Placeholder
	BasicPointerMulti(A *Placeholder) (A *Placeholder, B *Placeholder, C string)
	BasicExternal(*external.Placeholder) external.Placeholder
	BasicExternalMulti(*external.Placeholder) (external.Placeholder, *external.Placeholder)

	NoMatchArraySimple([16]byte) Collection
	Array([16]byte) [16]byte
	ArraySimple(Arr [16]byte) *Collection
	ArrayExternal([16]external.Placeholder) [16]external.Placeholder
	ArrayComplex(Arr [16]map[byte]string) *complex.Collection
	ArrayExternalComplex(Arr [16]map[*external.Collection]string) *complex.ComplexCollection

	NoMatchSliceSimple([]string) Collection
	Slice([]string) []string
	SliceSimple(S []string) *Collection
	SliceExternal([]external.Placeholder) []external.Placeholder
	SliceComplex(S []map[string][16]int) *complex.Collection
	SliceExternalComplex(S []map[string]func(*external.Collection) string) *complex.ComplexCollection

	NoMatchMap(map[string]bool) Collection
	Map(map[string]bool) map[string]bool
	MapSimple(M map[string]bool) *Collection
	MapExternal(map[string]external.Placeholder) map[string]external.Placeholder
	MapComplex(M map[string]interface{ Error() string }) *complex.Collection
	MapExternalComplex(M map[*external.Collection]external.Placeholder) *complex.ComplexCollection

	NoMatchChan(chan int) Collection
	Chan(chan int) chan int
	ChanSimple(C chan int) *Collection
	ChanExternal(chan external.Placeholder) chan external.Placeholder
	ChanComplex(C chan *[]int) *complex.Collection
	ChanExternalComplex(C chan *[]external.Collection) complex.ComplexCollection

	NoMatchInterface(error) Collection
	Interface(interface{}) interface{}
	InterfaceSimple(I error) *Collection
	InterfaceExternal(I error) *external.Collection
	InterfaceComplex(I interface{ A(rune) *int }) *complex.Collection
	InterfaceExternalComplex(I interface {
		A(string) map[*external.Collection]bool
		B() (int, byte)
	}) complex.ComplexCollection

	NoMatchFunc(func() int) Collection
	Func(func() int) func() int
	FuncSimple(F func() int) *Collection
	FuncExternal(func(external.Placeholder) int) func(external.Placeholder) int
	FuncComplex(F func([]string, uint64) *byte) *complex.Collection
	FuncExternalComplex(F func(external.Collection) []string) *complex.ComplexCollection

	NoMatchComplex([]external.Collection) (Struct []external.Collection)
	EmptyStruct(struct{}) empty
	Struct(Collection) Collection
	StructExternal(external.Collection) *Collection
}

// Collection represents a type that holds collection field types.
type Collection struct {
	Arr [16]byte
	S   []string
	M   map[string]bool
	C   chan int
	I   error
	F   func() int
}

// empty represents a struct that contains an empty struct.
type empty struct {
	e struct{}
}

// freefloat serves the purpose of checking for free-floating comments.
type freefloat struct {
	A string
}
