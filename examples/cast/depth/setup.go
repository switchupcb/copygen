package depth

type Copygen interface {
	ZeroDepth([]string) []string
	DefaultDepth([]string) One

	// cast .* .* 2
	CustomDepth([]string) Two

	// cast .* .* 2
	ReverseDepth(Two) []string

	// cast .* .* 0
	DisableCast([]string) []string
}
