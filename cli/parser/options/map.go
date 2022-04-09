package options

import (
	"fmt"
	"regexp"

	"github.com/switchupcb/copygen/cli/models"
)

const CategoryMap = "map"

// ParseMap parses a map option.
func ParseMap(option string) (*Option, error) {
	splitoption, err := splitOption(option, CategoryMap, "<option>:<whitespaces><regex><whitespaces><regex>")
	if err != nil {
		return nil, fmt.Errorf("%w", err)
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
func SetMap(field *models.Field, options []*Option) {
	// A map option can only be set to a field once, so use the last one
	for i := len(options) - 1; i > -1; i-- {
		if options[i].Category == CategoryMap &&
			options[i].Regex[0].MatchString(field.FullNameWithoutContainer("")) {
			if value, ok := options[i].Value.(string); ok {
				field.Options.Map = value
				break
			}
		}
	}
}
