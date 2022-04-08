package options

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/switchupcb/copygen/cli/models"
)

const CategoryDepth = "depth"

// ParseDepth parses a depth option.
func ParseDepth(option string) (*Option, error) {
	splitoption, err := splitOption(option, CategoryDepth, "<option>:<whitespaces><regex><whitespaces><int>")
	if err != nil {
		return nil, err
	}

	re, err := regexp.Compile("^" + splitoption[0] + "$")
	if err != nil {
		return nil, fmt.Errorf("an error occurred compiling the regex for a %s option: %q\n%v", CategoryDepth, option, err)
	}

	depth, err := strconv.Atoi(splitoption[1])
	if err != nil {
		return nil, fmt.Errorf("an error occurred parsing the integer depth value of a %s option: %q\n%v", CategoryDepth, option, err)
	}

	return &Option{
		Category: CategoryDepth,
		Regex:    map[int]*regexp.Regexp{0: re},
		Value:    depth,
	}, nil
}

// SetDepth sets a field's depth option.
func SetDepth(field *models.Field, options []*Option) {
	// A depth option can only be set to a field once, so use the last one
	for i := len(options) - 1; i > -1; i-- {
		if options[i].Category == CategoryDepth && options[i].Regex[0].MatchString(field.FullName("")) {
			if value, ok := options[i].Value.(int); ok {
				// Automatch all is on by default; if a user specifies 0 depth-level, guarantee it.
				if value == 0 {
					value = -1
				}

				field.Options.Depth = value

				break
			}
		}
	}
}
