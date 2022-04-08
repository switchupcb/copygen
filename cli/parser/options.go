package parser

import (
	"go/ast"
	"strings"

	"github.com/switchupcb/copygen/cli/parser/options"
)

// NOTE: This file refers to the Parser Options Process.
// See CONTRIBUTING.md#options

// CommentOptionMap represents a map of comments to an option.
type CommentOptionMap map[string]*options.Option

// MapCommentsToOptions parses a list of ast.Comments into a CommentOptionMap.
func MapOptions(comments []*ast.Comment) (CommentOptionMap, error) {
	optionmap := make(CommentOptionMap, len(comments))

	for _, comment := range comments {
		text := comment.Text
		splitcomments := strings.Fields(text[2:])

		if len(splitcomments) >= 1 {
			category := splitcomments[0]
			option := strings.Join(splitcomments[1:], " ")

			// convert options have already been set.
			if category != options.CategoryConvert {
				opt, err := assignFunctionOption(category, option)
				if err != nil {
					return nil, err
				}

				optionmap[text] = opt
			}
		}
	}

	return optionmap, nil
}

// assignFunctionOption assigns a function (field-oriented) option.
func assignFunctionOption(category, option string) (*options.Option, error) {
	switch category {
	case options.CategoryDeepCopy:
		opt, err := options.ParseDeepcopy(option)
		if err != nil {
			return nil, err
		}

		return opt, nil

	case options.CategoryDepth:
		opt, err := options.ParseDepth(option)
		if err != nil {
			return nil, err
		}

		return opt, nil

	case options.CategoryMap:
		opt, err := options.ParseMap(option)
		if err != nil {
			return nil, err
		}

		return opt, nil

	default:
		return &options.Option{
			Category: options.CategoryCustom,
			Regex:    nil,
			Value:    map[string]string{category: option},
		}, nil
	}
}
