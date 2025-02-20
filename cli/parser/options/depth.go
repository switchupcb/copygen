package options

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/switchupcb/copygen/cli/models"
)

const (
	CategoryDepth = "depth"

	// FormatDepth represents an end-user facing format for a depth option.
	// <option> refers to the "depth" option.
	FormatDepth = "<option>:<whitespaces><regex><whitespaces><int>"
)

// ParseDepth parses a depth option.
func ParseDepth(option string) (*Option, error) {
	splitoption := strings.Fields(option)
	if len(splitoption) == 0 {
		return nil, fmt.Errorf("there is an unspecified %s option at an unknown line", CategoryDepth)
	} else if len(splitoption) != 2 {
		return nil, fmt.Errorf("there is a misconfigured %s option: %q.\nIs it in format %s?", CategoryDepth, option, FormatDepth)
	}

	re, err := regexp.Compile("^" + splitoption[0] + "$")
	if err != nil {
		return nil, fmt.Errorf("an error occurred compiling the regex for a %s option: %q\n%w", CategoryDepth, option, err)
	}

	depth, err := strconv.Atoi(splitoption[1])
	if err != nil {
		return nil, fmt.Errorf("an error occurred parsing the integer depth value of a %s option: %q\n%w", CategoryDepth, option, err)
	}

	return &Option{
		Category: CategoryDepth,
		Regex:    map[int]*regexp.Regexp{0: re},
		Value:    depth, // int
	}, nil
}

// SetDepth sets a field's depth option.
func SetDepth(field *models.Field, option Option) {
	// A depth option can only be set to a field once.
	if field.Options.Depth != 0 {
		return
	}

	if option.Regex[0] != nil && option.Regex[0].MatchString(field.FullNameWithoutPointer("")) {
		if value, ok := option.Value.(int); ok {
			// Automatch all is on by default; if a user specifies 0 depth-level, guarantee it.
			if value == 0 {
				value = -1
			}

			field.Options.Depth = value
		}
	}
}

// FilterDepth filters a field according to it's depth level and returns the removed fields.
func FilterDepth(field *models.Field, maxdepth, curdepth int) []*models.Field {
	if maxdepth == 0 || len(field.Fields) == 0 {
		return nil
	}

	var removed []*models.Field
	if maxdepth < 0 || maxdepth <= curdepth {
		field.Fields, removed = removed, field.Fields
		return removed
	}

	for _, f := range field.Fields {
		removed = append(removed, FilterDepth(f, maxdepth+f.Options.Depth, curdepth+1)...)
	}

	return removed
}
