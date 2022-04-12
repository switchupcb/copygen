package options

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/switchupcb/copygen/cli/models"
)

const (
	CategoryMap = "map"

	// FormatMap represents an end-user facing format for map options.
	// <option> refers to the "map" option.
	FormatMap = "<option>:<whitespaces><regex><whitespaces><regex>"
)

// ParseMap parses a map option.
func ParseMap(option string) (*Option, error) {
	splitoption := strings.Fields(option)
	if len(splitoption) == 0 {
		return nil, fmt.Errorf("there is an unspecified %s option at an unknown line", CategoryMap)
	} else if len(splitoption) != 2 {
		return nil, fmt.Errorf("there is a misconfigured %s option: %q.\nIs it in format %s?", CategoryMap, option, FormatMap)
	}

	fromRe, err := regexp.Compile("^" + splitoption[0] + "$")
	if err != nil {
		return nil, fmt.Errorf("an error occurred compiling the regex for the from field in the %s option: %q\n%w", CategoryMap, option, err)
	}

	// map options are compared in the matcher
	return &Option{
		Category: CategoryMap,
		Regex:    map[int]*regexp.Regexp{0: fromRe},
		Value:    splitoption[1],
	}, nil
}

// SetMap sets a field's deepcopy option.
func SetMap(field *models.Field, option Option) {
	// a map option can only be set to a field once.
	if field.Options.Map != "" {
		return
	}

	// only one matching method may be specified for a field.
	if field.Options.Automatch {
		return
	}

	if option.Regex[0] != nil && option.Regex[0].MatchString(field.FullNameWithoutContainer("")) {
		if value, ok := option.Value.(string); ok {
			field.Options.Map = value
		}
	}
}
