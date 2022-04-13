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

// NewFieldOption creates a new field-oriented option from the given category and text.
func NewFieldOption(category, text string) (*Option, error) {
	var option *Option
	var err error

	switch category {
	case CategoryAutomatch:
		option, err = ParseAutomatch(text)

	case CategoryMap:
		option, err = ParseMap(text)

	case CategoryTag:
		option, err = ParseTag(text)

	case CategoryDeepcopy:
		option, err = ParseDeepcopy(text)

	case CategoryDepth:
		option, err = ParseDepth(text)

	default:
		option = &Option{
			Category: CategoryCustom,
			Regex:    nil,
			Value:    map[string]string{category: text},
		}
	}

	if err != nil {
		return nil, err
	}
	return option, nil
}
