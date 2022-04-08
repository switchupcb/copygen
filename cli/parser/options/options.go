// Package options parses function comments and sets them to fields.
package options

import (
	"fmt"
	"regexp"
	"strings"
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

// splitOption splits an option string and validates it.
func splitOption(option, category, format string) ([]string, error) {
	splitoption := strings.Fields(option)

	if len(splitoption) == 0 {
		return nil, fmt.Errorf("there is an unspecified %s option at an unknown line", category)
	} else if len(splitoption) == 1 || len(splitoption) > 2 {
		return nil, fmt.Errorf("there is a misconfigured %s option: %q.\nIs it in format %s?", category, option, format)
	}

	return splitoption, nil
}
