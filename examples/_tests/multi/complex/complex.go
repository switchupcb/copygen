package complex

import "github.com/switchupcb/copygen/examples/_tests/multi/external"

// Collection represents a type that holds collection field types.
type Collection struct {
	Arr [16]map[byte]string
	S   []map[string][16]int
	M   map[string]bool
	C   chan int
}

// ComplexCollection represents a type that holds collection field types.
type ComplexCollection struct {
	Arr [16]map[*external.Collection]string
	S   []map[string]func(*external.Collection) string
	M   map[string]bool
	C   chan int
}
