package expression

// Copygen defines the functions that are generated.
type Copygen interface {
	// cast int int * 2
	ExprDouble(int) int

	// map bool bool -cast ^ true
	MapExprXOR(bool) bool
}
