package options

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/switchupcb/copygen/cli/models"
)

// IsMatchOptionCategory determines if a given string is a match option category.
func IsMatchOptionCategory(c string) bool {
	return c == CategoryAutomatch || c == CategoryMap || c == CategoryTag
}

// IsMatchOptionSet determines if a match option is already set for a given field.
func IsMatchOptionSet(field models.Field) bool {
	return field.Options.Automatch || field.Options.Map != "" || field.Options.Tag != ""
}

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

	fieldRe, err := regexp.Compile("^" + splitoption[0] + "$")
	if err != nil {
		return nil, fmt.Errorf("an error occurred compiling the regex for the field in the %s option: %q\n%w", CategoryAutomatch, option, err)
	}

	return &Option{
		Category: CategoryAutomatch,
		Regex:    map[int]*regexp.Regexp{0: fieldRe},
		Value:    true,
	}, nil
}

// SetAutomatch sets a field's automatch option.
func SetAutomatch(field *models.Field, option Option) {
	if IsMatchOptionSet(*field) {
		return
	}

	if option.Regex[0] != nil && option.Regex[0].MatchString(field.FullNameWithoutContainer("")) {
		field.Options.Automatch = true
	}
}

const (
	CategoryMap = "map"

	// FormatMap represents an end-user facing format for map options.
	// <option> refers to the "map" option.
	FormatMap = "<option>:<whitespaces><regex><whitespaces><field>"
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

// SetMap sets a field's map option.
func SetMap(field *models.Field, option Option) {
	if IsMatchOptionSet(*field) {
		return
	}

	if option.Regex[0] != nil && option.Regex[0].MatchString(field.FullNameWithoutContainer("")) {
		if value, ok := option.Value.(string); ok {
			field.Options.Map = value
		}
	}
}

const (
	CategoryTag = "tag"

	// FormatTag represents an end-user facing format for tag options.
	// <option> refers to the "tag" option.
	FormatTag = "<option>:<whitespaces><regex><whitespaces><tag>"
)

// ParseTag parses a tag option.
func ParseTag(option string) (*Option, error) {
	splitoption := strings.Fields(option)
	if len(splitoption) == 0 {
		return nil, fmt.Errorf("there is an unspecified %s option at an unknown line", CategoryTag)
	} else if len(splitoption) != 2 {
		return nil, fmt.Errorf("there is a misconfigured %s option: %q.\nIs it in format %s?", CategoryTag, option, FormatTag)
	}

	fieldRe, err := regexp.Compile("^" + splitoption[0] + "$")
	if err != nil {
		return nil, fmt.Errorf("an error occurred compiling the regex for the field in the %s option: %q\n%w", CategoryTag, option, err)
	}

	return &Option{
		Category: CategoryTag,
		Regex:    map[int]*regexp.Regexp{0: fieldRe},
		Value:    splitoption[1],
	}, nil
}

// SetTag sets a field's tag option.
func SetTag(field *models.Field, option Option) {
	if IsMatchOptionSet(*field) {
		return
	}

	if option.Regex[0] != nil && option.Regex[0].MatchString(field.FullNameWithoutContainer("")) {
		if optionvalue, ok := option.Value.(string); ok {
			for tagcat, tagmeta := range field.Tags {

				// match the tag by it's category.
				// i.e api == api in `api:"id"`.
				if optionvalue == tagcat {

					// gets the name of any valid tag.
					// i.e `api:id` in `api:"id", api:"name"` (invalid).
					for tagname := range tagmeta {
						field.Options.Tag = tagcat + ":" + tagname
						return
					}
				}
			}
		}
	}
}
