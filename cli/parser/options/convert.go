package options

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/switchupcb/copygen/cli/models"
)

const (
	CategoryConvert = "convert"

	// FormatConvert represents an end-user facing format for convert options.
	// <option> refers to the "convert" option.
	FormatConvert = "<option>:<whitespaces><regex><whitespaces><regex>"
)

// ParseConvert parses a convert option.
func ParseConvert(option, value string) (*Option, error) {
	splitoption := strings.Fields(option)
	if len(splitoption) == 0 {
		return nil, fmt.Errorf("there is an unspecified %s option at an unknown line", CategoryConvert)
	} else if len(splitoption) != 2 {
		return nil, fmt.Errorf("there is a misconfigured %s option: %q.\nIs it in format %s?", CategoryConvert, option, FormatConvert)
	}

	funcRe, err := regexp.Compile("^" + splitoption[0] + "$")
	if err != nil {
		return nil, fmt.Errorf("an error occurred compiling the regex for the first field in the %s option: %q\n%w", CategoryConvert, option, err)
	}

	fieldRe, err := regexp.Compile("^" + splitoption[1] + "$")
	if err != nil {
		return nil, fmt.Errorf("an error occurred compiling the regex for the second field in the %s option: %q\n%w", CategoryConvert, option, err)
	}

	return &Option{
		Category: CategoryConvert,
		Regex:    map[int]*regexp.Regexp{0: funcRe, 1: fieldRe},
		Value:    value,
	}, nil
}

// SetConvert sets a field's convert option.
func SetConvert(field *models.Field, option Option, functionName string) {
	// A convert option can only be set to a field once.
	if field.Options.Convert != "" {
		return
	}

	if option.Regex[0] != nil && option.Regex[0].MatchString(functionName) {
		if option.Regex[1] != nil && option.Regex[1].MatchString(field.FullNameWithoutPointer("")) {
			if value, ok := option.Value.(string); ok {
				field.Options.Convert = value
			}
		}
	}

}
