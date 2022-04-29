package complex

import "github.com/switchupcb/copygen/examples/_tests/multi/external"

// Collection represents a type that holds collection field types.
type Collection struct {
	Arr [16]map[byte]string
	S   []map[string][16]int
	M   map[string]interface{ Error() string }
	C   chan *[]int
	I   interface{ A(rune) *int }
	F   func([]string, uint64) *byte
}

// ComplexCollection represents a type that holds collection field types.
type ComplexCollection struct {
	Arr [16]map[*external.Collection]string
	S   []map[string]func(*external.Collection) string
	M   map[*external.Collection]external.Placeholder
	C   chan *[]external.Collection
	I   interface {
		A(string) map[*external.Collection]bool
		B() (int, byte)
	}
	F func(external.Collection) []string
}
