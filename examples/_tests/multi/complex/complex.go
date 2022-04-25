package complex

import "github.com/switchupcb/copygen/examples/_tests/multi/external"

// Collection represents a type that holds collection field types.
type Collection struct {
	Arr [16]map[byte]string
	S   []map[string][16]int
	M   map[string]error
	C   chan *[]int
}

// ComplexCollection represents a type that holds collection field types.
type ComplexCollection struct {
	Arr [16]map[*external.Collection]string
	S   []map[string]func(*external.Collection) string
	M   map[*external.Collection]external.Placeholder
	C   chan *[]external.Collection
}

// Container represents a type that holds container field types.
type Container struct {
	I error
	F func([]string, uint64) *byte
}

// ComplexContainer represents a type that holds container field types.
type ComplexContainer struct {
	I interface{ Error() []external.Collection }
	F func(external.Collection) []string
}
