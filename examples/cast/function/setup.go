package function

// Copygen defines the functions that will be generated.
type Copygen interface {
	// cast Custom string .String()
	TypeFuncString(Custom) string

	// cast Custom string Convert()
	FuncString(Custom) string
}
