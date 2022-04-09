// Package options parses function comments and sets them to fields.
package options

import (
	"regexp"
)

// Option represents an option applied to functions and fields.
type Option struct {
	// The compiled regex the option uses for its arguments (map[position]regex).
	Regex map[int]*regexp.Regexp

	// The values to assign to a type (function or field) if the option applies.
	Value interface{}

	// The category the option falls under.
	// There are currently five: convert, depth, deepcopy, map, custom
	Category string
}
