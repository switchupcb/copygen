package options

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/switchupcb/copygen/cli/models"
)

const (
	CategoryCast = "cast"

	// FormatCast represents an end-user facing format for a cast option.
	// <option> refers to the "cast" option.
	FormatCast = "<option><whitespaces><regex><whitespaces><field><whitespaces><modifier>"

	// FormatModifierCast represents an end-user facing format for a cast option modifier.
	// <option> refers to the "cast" option.
	FormatModifierCast = "-<option><whitespaces><regex><whitespaces><modifier>"
)

// ParseCast parses a cast option.
func ParseCast(option string) (*Option, error) {
	splitoption := strings.Fields(option)
	if len(splitoption) == 0 {
		return nil, fmt.Errorf("there is an unspecified %s option at an unknown line", CategoryCast)
	} else if len(splitoption) < 2 {
		return nil, fmt.Errorf("there is a misconfigured %s option: %q.\nIs it in format %s?", CategoryCast, option, FormatCast)
	}

	fromRe, err := regexp.Compile("^" + splitoption[0] + "$")
	if err != nil {
		return nil, fmt.Errorf("an error occurred compiling the regex for the from-field in the %s option: %q\n%w", CategoryCast, option, err)
	}

	return &Option{
		Category: CategoryCast,
		Regex:    map[int]*regexp.Regexp{0: fromRe},
		Value:    []string{splitoption[1], strings.Join(splitoption[2:], " ")}, // []string{from-field, modifier}
	}, nil
}

// ParseModifierCast parses a cast option modifier.
func ParseModifierCast(option string) (*Option, error) {
	// - cast
	// - cast modifier...

	return nil, nil
}

// SetCast sets a field's cast option.
func SetCast(field *models.Field, option Option) {
	// if IsMatchOptionSet(*field) {
	// 	return
	// }

	// if option.Regex[0] != nil && option.Regex[0].MatchString(field.FullNameWithoutPointer("")) {
	// 	if value, ok := option.Value.(string); ok {
	// 		field.Options.Map = value
	// 	}
	// }
}
