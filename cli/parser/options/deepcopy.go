package options

import (
	"fmt"
	"regexp"

	"github.com/switchupcb/copygen/cli/models"
)

const CategoryDeepCopy = "deepcopy"

// ParseDeepcopy parses a deepcopy option.
func ParseDeepcopy(option string) (*Option, error) {
	re, err := regexp.Compile("^" + option + "$")
	if err != nil {
		return nil, fmt.Errorf("an error occurred compiling the regex for a %s option: %q\n%w", CategoryDeepCopy, option, err)
	}

	return &Option{
		Category: CategoryDeepCopy,
		Regex:    map[int]*regexp.Regexp{0: re},
		Value:    true,
	}, nil
}

// SetDeepcopy sets a field's deepcopy option.
func SetDeepcopy(field *models.Field, options []*Option) {
	// A deepcopy option can only be set to a field once, so use the last one
	for i := len(options) - 1; i > -1; i-- {
		if options[i].Category == CategoryDeepCopy &&
			options[i].Regex[0].MatchString(field.FullNameWithoutContainer("")) {
			field.Options.Deepcopy = true
			break
		}
	}
}
