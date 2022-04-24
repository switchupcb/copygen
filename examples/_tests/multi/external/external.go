package external

type Placeholder bool

// Collection represents a type that holds collection field types.
type Collection struct {
	Arr [16]byte
	S   []string
	M   map[string]bool
	C   chan int
}
