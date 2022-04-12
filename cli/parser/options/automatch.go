package options

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/switchupcb/copygen/cli/models"
)

const (
	CategoryAutomatch = "automatch"

	// FormatAutomatch represents an end-user facing format for automatch options.
	// <option> refers to the "automatch" option.
	FormatAutomatch = "<option>:<whitespaces><regex>"
)

// ParseAutomatch parses a automatch option.
func ParseAutomatch(option string) (*Option, error) {
	splitoption := strings.Fields(option)
	if len(splitoption) == 0 {
		return nil, fmt.Errorf("there is an unspecified %s option at an unknown line", CategoryAutomatch)
	} else if len(splitoption) != 1 {
		return nil, fmt.Errorf("there is a misconfigured %s option: %q.\nIs it in format %s?", CategoryAutomatch, option, FormatAutomatch)
	}

	fromRe, err := regexp.Compile("^" + splitoption[0] + "$")
	if err != nil {
		return nil, fmt.Errorf("an error occurred compiling the regex for the from field in the %s option: %q\n%w", CategoryAutomatch, option, err)
	}

	// map options are compared in the matcher
	return &Option{
		Category: CategoryAutomatch,
		Regex:    map[int]*regexp.Regexp{0: fromRe},
		Value:    true,
	}, nil
}

// SetAutomatch sets a field's deepcopy option.
func SetAutomatch(field *models.Field, option Option) {
	// an automatch option can only be set to a field once.
	if field.Options.Automatch {
		return
	}

	// only one matching method may be specified for a field.
	if field.Options.Map != "" {
		return
	}

	if option.Regex[0] != nil && option.Regex[0].MatchString(field.FullNameWithoutContainer("")) {
		field.Options.Automatch = true
	}
}
