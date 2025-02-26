package convert

// Copygen defines the functions that are generated.
type Copygen interface {
	// cast convert.Placeholder bool
	NoConvertPlaceholder(Placeholder) bool
	// cast bool convert.Placeholder
	NoConvertBool(bool) Placeholder
	// map convert.Placeholder bool
	// cast convert.Placeholder bool
	MapConvertPlaceholder(Placeholder) bool

	// map convert.Placeholder int
	// cast convert.Placeholder int + 5
	MapConvertPlaceholderWithExpression(Placeholder) int
}
