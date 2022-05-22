package options

import (
	"fmt"
)

const (
	CategoryCustom = "custom"
	FormatCustom   = "<option><whitespaces><value>"
)

// MapCustomOption maps a custom option in an optionmap[category][]values.
func MapCustomOption(optionmap map[string][]string, option *Option) (map[string][]string, error) {
	if optionmap == nil {
		optionmap = make(map[string][]string)
	}

	if option.Category == CategoryCustom {
		if customoptionvalue, ok := option.Value.(map[string]string); ok {
			for k, v := range customoptionvalue {
				optionmap[k] = append(optionmap[k], v)
			}
		} else {
			return optionmap, fmt.Errorf("failed to map custom option: %v", option.Value)
		}
	}

	return optionmap, nil
}

// MapCustomOptions maps options with custom categories in a list of options to a customoptionmap[category][]value.
func MapCustomOptions(options []*Option) (map[string][]string, error) {
	var (
		optionmap map[string][]string
		err       error
	)

	for _, option := range options {
		optionmap, err = MapCustomOption(optionmap, option)
		if err != nil {
			return nil, err
		}
	}

	return optionmap, nil
}
