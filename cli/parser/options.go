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

	// The values to assign to a type (function or field) if the option applies (map[position]value).
	Value map[int]interface{}
}

// parseOptions parses the AST for options in the setup file.
func (p *Parser) parseOptions() error {
	if p.Options == nil {
		p.Options = make(OptionMap)
	}

	for _, commentgroup := range p.SetupFile.Comments {
		for _, comment := range commentgroup.List {
			text := comment.Text
			splitcomments := strings.Split(text[2:], ":")

			// map[comment]map[optionname]map[]
			// determine if the comment is an option.
			if len(splitcomments) >= 2 {
				category := strings.TrimSpace(splitcomments[0])
				option := strings.TrimSpace(strings.Join(splitcomments[1:], ":"))
				switch category {
				case "convert":
					re, err := regexp.Compile(option)
					if err != nil {
						return fmt.Errorf("An error occurred compiling the regex for a convert option: %q\n%v", option, err)
					}

					// funcname assigned in parseKeep()
					p.Options[text] = Option{
						Category: "convert",
						Regex:    map[int]*regexp.Regexp{0: re},
						Value:    make(map[int]interface{}, 1),
					}
				case "deepcopy":
					re, err := regexp.Compile(option)
					if err != nil {
						return fmt.Errorf("An error occurred compiling the regex for a deepcopy option: %q\n%v", option, err)
					}
					p.Options[text] = Option{
						Category: "deepcopy",
						Regex:    map[int]*regexp.Regexp{0: re},
						Value:    map[int]interface{}{0: true},
					}
				case "depth":
					opt, err := parseDepth(option)
					if err != nil {
						return err
					}
					p.Options[text] = *opt
				case "map":
					opt, err := parseMap(option)
					if err != nil {
						return err
					}
					p.Options[text] = *opt
				default:
					p.Options[text] = Option{
						Category: "custom",
						Regex:    nil,
						Value:    map[int]interface{}{0: map[string]string{category: option}},
					}
				}
			}
		}
	}
	return nil
}

// parseDepth parses a depth option.
func parseDepth(option string) (*Option, error) {
	splitoption := strings.Split(option, " ")
	if len(splitoption) == 0 {
		return nil, fmt.Errorf("There is an unspecified depth option at an unknown line.")
	} else if len(splitoption) == 1 || len(splitoption) > 2 {
		return nil, fmt.Errorf("There is a misconfigured depth option: %q.\nIs it in format <option>:<whitespaces><regex><whitespaces><int>?", option)
	}

	re, err := regexp.Compile(splitoption[0])
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
		Value:    map[int]interface{}{0: depth},
	}, nil
}

// parseMap parses a map option.
func parseMap(option string) (*Option, error) {
	splitoption := strings.Split(option, " ")
	if len(splitoption) == 0 {
		return nil, fmt.Errorf("There is an unspecified map option at an unknown line.")
	} else if len(splitoption) == 1 || len(splitoption) > 2 {
		return nil, fmt.Errorf("There is a misconfigured map option: %q.\nIs it in format <option>:<whitespaces><regex><whitespaces><regex>?", option)
	}

	fromRe, err := regexp.Compile(splitoption[0])
	if err != nil {
		return nil, fmt.Errorf("An error occurred compiling the regex for the first field in the map option: %q.\n%v", option, err)
	}

	toRe, err := regexp.Compile(splitoption[1])
	if err != nil {
		return nil, fmt.Errorf("An error occurred compiling the regex for the second field in the map option: %q.\n%v", option, err)
	}

	return &Option{
		Category: "map",
		Regex:    map[int]*regexp.Regexp{0: fromRe},
		Value:    map[int]interface{}{0: toRe}, // compared in the matcher
	}, nil
}
