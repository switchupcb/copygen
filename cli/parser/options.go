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
	// The compiled regex the option uses for its arguments (map[position]regex).
	Regex map[int]*regexp.Regexp

	// The values to assign to a type (function or field) if the option applies.
	Value interface{}

	// The category the option falls under.
	// There are currently five: convert, depth, deepcopy, map, custom
	Category string
}

const (
	categoryConvert   = "convert"
	categoryCustom    = "custom"
	categoryDeepCopy  = "deepcopy"
	categoryDepth     = "depth"
	categoryMap       = "map"
	categoryCommonTag = "tag"
)

// parseConvert parses a convert option.
func parseConvert(option, value string) (*Option, error) {
	splitoption, err := splitOption(option, categoryConvert, "<option>:<whitespaces><tag name><whitespaces><regex>")
	if err != nil {
		return nil, err
	}

	funcRe, err := regexp.Compile("^" + splitoption[0] + "$")
	if err != nil {
		return nil, fmt.Errorf("an error occurred compiling the regex for the first field in the %s option: %q\n%v", categoryConvert, option, err)
	}

	fieldRe, err := regexp.Compile("^" + splitoption[1] + "$")
	if err != nil {
		return nil, fmt.Errorf("an error occurred compiling the regex for the second field in the %s option: %q\n%v", categoryConvert, option, err)
	}

	return &Option{
		Category: categoryConvert,
		Regex:    map[int]*regexp.Regexp{0: funcRe, 1: fieldRe},
		Value:    value,
	}, nil
}

// parseDeepcopy parses a deepcopy option.
func parseDeepcopy(option string) (*Option, error) {
	re, err := regexp.Compile("^" + option + "$")
	if err != nil {
		return nil, fmt.Errorf("an error occurred compiling the regex for a %s option: %q\n%v", categoryDeepCopy, option, err)
	}

	return &Option{
		Category: categoryDeepCopy,
		Regex:    map[int]*regexp.Regexp{0: re},
		Value:    true,
	}, nil
}

// parseDepth parses a depth option.
func parseDepth(option string) (*Option, error) {
	splitoption, err := splitOption(option, categoryDepth, "<option>:<whitespaces><regex><whitespaces><int>")
	if err != nil {
		return nil, err
	}

	re, err := regexp.Compile("^" + splitoption[0] + "$")
	if err != nil {
		return nil, fmt.Errorf("an error occurred compiling the regex for a %s option: %q\n%v", categoryDepth, option, err)
	}

	depth, err := strconv.Atoi(splitoption[1])
	if err != nil {
		return nil, fmt.Errorf("an error occurred parsing the integer depth value of a %s option: %q\n%v", categoryDepth, option, err)
	}

	return &Option{
		Category: categoryDepth,
		Regex:    map[int]*regexp.Regexp{0: re},
		Value:    depth,
	}, nil
}

// parseMap parses a map option.
func parseMap(option string) (*Option, error) {
	splitoption, err := splitOption(option, categoryMap, "<option>:<whitespaces><regex><whitespaces><regex>")
	if err != nil {
		return nil, err
	}

	fromRe, err := regexp.Compile("^" + splitoption[0] + "$")
	if err != nil {
		return nil, fmt.Errorf("an error occurred compiling the regex for the from field in the %s option: %q\n%v", categoryMap, option, err)
	}

	// map options are compared in the matcher
	return &Option{
		Category: categoryMap,
		Regex:    map[int]*regexp.Regexp{0: fromRe},
		Value:    splitoption[1],
	}, nil
}

// splitOption splits option string and validates it.
func splitOption(option, category, format string) ([]string, error) {
	splitoption := strings.Fields(option)

	if len(splitoption) == 0 {
		return nil, fmt.Errorf("there is an unspecified %s option at an unknown line", category)
	} else if len(splitoption) == 1 || len(splitoption) > 2 {
		return nil, fmt.Errorf("there is a misconfigured %s option: %q.\nIs it in format %s?", category, option, format)
	}

	return splitoption, nil
}

// parseMatchByTag parses a map option.
func parseMatchByTag(option string) (*Option, error) {
	splitoption, err := splitOption(option, categoryCommonTag, "<option>:<whitespaces><model><whitespaces><tag>")
	if err != nil {
		return nil, err
	}

	fromRe, err := regexp.Compile("^" + splitoption[0] + "$")
	if err != nil {
		return nil, fmt.Errorf("an error occurred compiling the regex for the from field in the %s option: %q\n%v", categoryMap, option, err)
	}

	// map options are compared in the matcher
	return &Option{
		Category: categoryCommonTag,
		Regex:    map[int]*regexp.Regexp{0: fromRe},
		Value:    splitoption[1],
	}, nil
}
