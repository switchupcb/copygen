package convert

// Copygen defines the functions that are generated.
type Copygen interface {
	// cast bool Placeholder
	ConvertBool(bool) Placeholder
	// cast Placeholder bool
	ConvertPlaceholder(Placeholder) bool
	// map Placeholder bool
	// cast Placeholder bool
	MapConvertPlaceholder(Placeholder) bool
}
