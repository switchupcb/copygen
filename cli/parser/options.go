package parser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// OptionMap represents a map of comment-option relations (map[comment]Option).
type OptionMap map[string]Option

// Option represents an option applied to functions and fields.
type Option struct {
	// The category the option falls under.
	// There are currently five: convert, depth, deepcopy, map, custom
	Category string

	// The compiled regex the option uses for its arguments (map[position]regex).
	Regex map[int]*regexp.Regexp

	// The values to assign to a type (function or field) if the option applies.
	Value interface{}
}

// parseConvert parses a convert option.
func parseConvert(option, value string) (*Option, error) {
	splitoption := strings.Fields(option)
	if len(splitoption) == 0 {
		return nil, fmt.Errorf("There is an unspecified convert option at an unknown line.")
	} else if len(splitoption) == 1 || len(splitoption) > 2 {
		return nil, fmt.Errorf("There is a misconfigured convert option: %q.\nIs it in format <option>:<whitespaces><regex><whitespaces><regex>?", option)
	}

	funcRe, err := regexp.Compile("^" + splitoption[0] + "$")
	if err != nil {
		return nil, fmt.Errorf("An error occurred compiling the regex for the first field in the convert option: %q.\n%v", option, err)
	}

	fieldRe, err := regexp.Compile("^" + splitoption[1] + "$")
	if err != nil {
		return nil, fmt.Errorf("An error occurred compiling the regex for the second field in the convert option: %q.\n%v", option, err)
	}

	return &Option{
		Category: "convert",
		Regex:    map[int]*regexp.Regexp{0: funcRe, 1: fieldRe},
		Value:    value,
	}, nil
}

// parseDeepcopy parses a deepcopy option.
func parseDeepcopy(option string) (*Option, error) {
	re, err := regexp.Compile("^" + option + "$")
	if err != nil {
		return nil, fmt.Errorf("An error occurred compiling the regex for a deepcopy option: %q\n%v", option, err)
	}
	return &Option{
		Category: "deepcopy",
		Regex:    map[int]*regexp.Regexp{0: re},
		Value:    true,
	}, nil
}

// parseDepth parses a depth option.
func parseDepth(option string) (*Option, error) {
	splitoption := strings.Fields(option)
	if len(splitoption) == 0 {
		return nil, fmt.Errorf("There is an unspecified depth option at an unknown line.")
	} else if len(splitoption) == 1 || len(splitoption) > 2 {
		return nil, fmt.Errorf("There is a misconfigured depth option: %q.\nIs it in format <option>:<whitespaces><regex><whitespaces><int>?", option)
	}

	re, err := regexp.Compile("^" + splitoption[0] + "$")
	if err != nil {
		return nil, fmt.Errorf("An error occurred compiling the regex for a depth option: %q.\n%v", option, err)
	}

	depth, err := strconv.Atoi(splitoption[1])
	if err != nil {
		return nil, fmt.Errorf("An error occurred parsing the integer depth value of a depth option: %q\n%v", option, err)
	}

	return &Option{
		Category: "depth",
		Regex:    map[int]*regexp.Regexp{0: re},
		Value:    depth,
	}, nil
}

// parseMap parses a map option.
func parseMap(option string) (*Option, error) {
	splitoption := strings.Fields(option)
	if len(splitoption) == 0 {
		return nil, fmt.Errorf("There is an unspecified map option at an unknown line.")
	} else if len(splitoption) == 1 || len(splitoption) > 2 {
		return nil, fmt.Errorf("There is a misconfigured map option: %q.\nIs it in format <option>:<whitespaces><regex><whitespaces><regex>?", option)
	}

	fromRe, err := regexp.Compile("^" + splitoption[0] + "$")
	if err != nil {
		return nil, fmt.Errorf("An error occurred compiling the regex for the from field in the map option: %q.\n%v", option, err)
	}

	// map options are compared in the matcher
	return &Option{
		Category: "map",
		Regex:    map[int]*regexp.Regexp{0: fromRe},
		Value:    splitoption[1],
	}, nil
}
