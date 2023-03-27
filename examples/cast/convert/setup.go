package convert

// Copygen defines the functions that will be generated.
type Copygen interface {
	// cast bool Placeholder
	ConvertBool(bool) Placeholder
	// cast Placeholder bool
	ConvertPlaceholder(Placeholder) bool
	// map Placeholder bool -cast
	MapConvertPlaceholder(Placeholder) bool
}
