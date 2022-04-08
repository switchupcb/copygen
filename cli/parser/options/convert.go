package options

import (
	"fmt"
	"regexp"

	"github.com/switchupcb/copygen/cli/models"
)

const CategoryConvert = "convert"

// ParseConvert parses a convert option.
func ParseConvert(option, value string) (*Option, error) {
	splitoption, err := splitOption(option, CategoryConvert, "<option>:<whitespaces><regex><whitespaces><regex>")
	if err != nil {
		return nil, err
	}

	funcRe, err := regexp.Compile("^" + splitoption[0] + "$")
	if err != nil {
		return nil, fmt.Errorf("an error occurred compiling the regex for the first field in the %s option: %q\n%v", CategoryConvert, option, err)
	}

	fieldRe, err := regexp.Compile("^" + splitoption[1] + "$")
	if err != nil {
		return nil, fmt.Errorf("an error occurred compiling the regex for the second field in the %s option: %q\n%v", CategoryConvert, option, err)
	}

	return &Option{
		Category: CategoryConvert,
		Regex:    map[int]*regexp.Regexp{0: funcRe, 1: fieldRe},
		Value:    value,
	}, nil
}

// SetConvert sets a field's convert option.
func SetConvert(field *models.Field, options []*Option) {
	// A convert option can only be set to a field once, so use the last one
	for i := len(options) - 1; i > -1; i-- {
		if options[i].Category == CategoryConvert && options[i].Regex[1].MatchString(field.FullName("")) {
			if value, ok := options[i].Value.(string); ok {
				field.Options.Convert = value
				break
			}
		}
	}
}
