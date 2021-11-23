package parser

import (
	"fmt"
	"go/ast"
	"go/types"

	"github.com/switchupcb/copygen/cli/models"
)

// parseFunctions parses the AST for functions in the setup file.
func (p *Parser) parseFunctions() ([]models.Function, error) {
	functions := make([]models.Function, 0)
	for _, def := range p.pkg.TypesInfo.Defs {
		if def != nil && def.Name() == "Copygen" {
			obj := def.(*types.TypeName)
			if t, ok := obj.Type().Underlying().(*types.Interface); ok {
				for i := 0; i < t.NumMethods(); i++ {
					method := t.Method(i)
					decl := p.pkg.Syntax[0].Scope.Lookup("Copygen").Decl.(*ast.TypeSpec)
					doc := decl.Type.(*ast.InterfaceType).Methods.List[i].Doc
					options, manual := p.filterOptionMap(doc)
					parsed, err := p.parseTypes(method.Type().(*types.Signature), options)
					if err != nil {
						return nil, fmt.Errorf("an error occurred while parsing the types of function %q.\n%v", method.Name(), err)
					}

					function := models.Function{
						Name: method.Name(),
						To:   parsed.toTypes,
						From: parsed.fromTypes,
						Options: models.FunctionOptions{
							Custom: p.assignCustomOption(options),
							Manual: manual,
						},
					}

					functions = append(functions, function)
				}
			}
		}
	}
	return functions, nil
}

// filterOptionMap filters an Option map for options that only pertain to the fields of a function.
// To reduce overhead, it also returns whether the function uses a manual matcher.
func (p *Parser) filterOptionMap(x *ast.CommentGroup) ([]Option, bool) {
	var (
		options []Option
		manual  bool
	)
	if x != nil {
		for _, comment := range x.List {
			if _, exists := p.Options[comment.Text]; exists {
				options = append(options, p.Options[comment.Text])
				if p.Options[comment.Text].Category == categoryMap {
					manual = true
				}
			}
		}
	}

	// add all convert options; which aren't in the scope of any functions but may apply
	for _, option := range p.Options {
		if option.Category == categoryConvert {
			options = append(options, option)
		}
	}

	return options, manual
}

// assignCustomOption parses a functions *ast.CommentGroups for custom options to return a Custom map.
func (p *Parser) assignCustomOption(options []Option) map[string][]string {
	optionmap := make(map[string][]string)

	// functions only have custom options
	for i := 0; i < len(options); i++ {
		switch options[i].Category {
		case categoryConvert, categoryDeepCopy, categoryDepth, categoryMap:
		default:
			if customoptionmap, ok := options[i].Value.(map[string]string); ok {
				for customoption, value := range customoptionmap {
					optionmap[customoption] = append(optionmap[customoption], value)
				}
			} else if customoptionmap != nil {
				fmt.Printf("WARNING: Failed to assign custom option: %v\n", options[i].Value)
			}
		}
	}

	return optionmap
}
