package options

import (
	"fmt"
	"regexp"

	"github.com/switchupcb/copygen/cli/models"
)

const CategoryDeepcopy = "deepcopy"

// ParseDeepcopy parses a deepcopy option.
func ParseDeepcopy(option string) (*Option, error) {
	re, err := regexp.Compile("^" + option + "$")
	if err != nil {
		return nil, fmt.Errorf("an error occurred compiling the regex for a %s option: %q\n%w", CategoryDeepcopy, option, err)
	}

	return &Option{
		Category: CategoryDeepcopy,
		Regex:    map[int]*regexp.Regexp{0: re},
		Value:    true,
	}, nil
}

// SetDeepcopy sets a field's deepcopy option.
func SetDeepcopy(field *models.Field, option Option) {
	// A deepcopy option can only be set to a field once.
	if field.Options.Deepcopy {
		return
	}

	if option.Regex[0] != nil && option.Regex[0].MatchString(field.FullNameWithoutPointer("")) {
		field.Options.Deepcopy = true
	}
}
