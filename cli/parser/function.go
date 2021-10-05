package parser

import (
	"fmt"
	"go/ast"

	"github.com/switchupcb/copygen/cli/models"
)

// parseFunctions parses the AST for functions in the setup file.
func (p *Parser) parseFunctions(copygen *ast.InterfaceType) ([]models.Function, error) {
	functions := make([]models.Function, 0, len(copygen.Methods.List))

	for _, method := range copygen.Methods.List {
		options, manual := p.setOptionMap(method)
		fieldsearcher := FieldSearcher{Options: options}
		parsed, err := p.parseTypes(method, &fieldsearcher)
		if err != nil {
			return nil, fmt.Errorf("an error occurred while parsing the types of function %q.\n%v", parseMethodForName(method), err)
		}

		function := models.Function{
			Name: parseMethodForName(method),
			To:   parsed.toTypes,
			From: parsed.fromTypes,
			Options: models.FunctionOptions{
				Custom: p.assignCustomOption(options),
				Manual: manual,
			},
		}

		functions = append(functions, function)
	}
	return functions, nil
}

// setOptionMap filters an Option map for options that only pertain to the fields of a function.
// To reduce overhead, it also returns whether the function uses a manual matcher.
func (p *Parser) setOptionMap(x ast.Node) ([]Option, bool) {
	var options []Option
	var manual bool
	ast.Inspect(x, func(node ast.Node) bool {
		if xcg, ok := node.(*ast.CommentGroup); ok {
			for _, comment := range xcg.List {
				if _, exists := p.Options[comment.Text]; exists {
					options = append(options, p.Options[comment.Text])
					if p.Options[comment.Text].Category == categoryMap {
						manual = true
					}
				}
			}
		}
		return true
	})

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
		case "convert":
		case "deepcopy":
		case "depth":
		case "map":
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

// parseMethodForName parses a method inside of a Copygen interface to provide its name.
func parseMethodForName(method *ast.Field) string {
	var funcname string // i.e 'ModelsToDomain' in func ModelsToDomain(models.Account, *models.User) *domain.Account

	// ast Note: "Field.Names contains a single name "type" for elements of interface type lists"
	for _, name := range method.Names {
		funcname += name.String() // i.e ModelsToDomain
	}
	return funcname
}
