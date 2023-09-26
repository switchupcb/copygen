// Package options parses function comments and sets them to fields.
package options

import (
	"regexp"

	"github.com/switchupcb/copygen/cli/models"
)

// Option represents an option applied to functions and fields.
type Option struct {
	// The compiled regex the option uses for its arguments (map[position]regex).
	Regex map[int]*regexp.Regexp

	// The values to assign to a type (function or field) if the option applies.
	Value interface{}

	// The category the option falls under.
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

// SetFieldOptions sets a field's (and its subfields) options.
func SetFieldOptions(field *models.Field, fieldoptions []*Option, functionName string) {
	for _, option := range fieldoptions {

		switch option.Category {

		case CategoryAutomatch:
			SetAutomatch(field, *option)

		case CategoryMap:
			SetMap(field, *option)

		case CategoryTag:
			SetTag(field, *option)

		case CategoryConvert:
			SetConvert(field, *option, functionName)

		case CategoryDepth:
			SetDepth(field, *option)

		case CategoryDeepcopy:
			SetDeepcopy(field, *option)

		case CategoryCustom:
			SetConvert(field, *option, functionName)
		}
	}
}
