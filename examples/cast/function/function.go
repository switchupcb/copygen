package function

// Custom represents a custom type.
type Custom struct{}

func (Custom) String() string {
	return ""
}

// Convert converts a Custom struct into a string.
func Convert(Custom) string {
	return ""
}
