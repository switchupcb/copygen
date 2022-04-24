package complex

// Collection represents a type that holds collection field types.
type Collection struct {
	Arr [16]map[byte]string
	S   []string
	M   map[string]bool
	C   chan int
}
