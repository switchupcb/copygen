package parser

import (
	"fmt"
	"go/ast"

	"github.com/switchupcb/copygen/cli/models"
	"github.com/switchupcb/copygen/cli/parser/options"
)

// parseFunctions parses the AST for functions in the setup file.
func (p *Parser) parseFunctions(copygen *ast.InterfaceType) ([]models.Function, error) {
	functions := make([]models.Function, 0, len(copygen.Methods.List))

	for _, method := range copygen.Methods.List {
		options, manual := p.getNodeOptions(method)

		parsed, err := p.parseTypes(method, options)
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

// parseMethodForName parses a method inside of a Copygen interface to provide its name.
func parseMethodForName(method *ast.Field) string {
	var funcname string // i.e 'ModelsToDomain' in func ModelsToDomain(models.Account, *models.User) *domain.Account

	// ast Note: "Field.Names contains a single name "type" for elements of interface type lists"
	for _, name := range method.Names {
		funcname += name.String() // i.e ModelsToDomain
	}

	return funcname
}

// getNodeOptions gets an ast.Node options from its comments.
// To reduce overhead, it also returns whether the function uses a manual matcher.
func (p *Parser) getNodeOptions(x ast.Node) ([]*options.Option, bool) {
	nodeOptions := make([]*options.Option, 0, len(p.CommentOptionMap))
	var manual bool

	ast.Inspect(x, func(node ast.Node) bool {
		commentGroup, ok := node.(*ast.CommentGroup)
		if !ok {
			return true
		}

		for _, comment := range commentGroup.List {
			if p.CommentOptionMap[comment.Text] != nil {
				nodeOptions = append(nodeOptions, p.CommentOptionMap[comment.Text])
				if p.CommentOptionMap[comment.Text].Category == options.CategoryMap {
					manual = true
				}
			}
		}

		return true
	})

	return nodeOptions, manual
}

// assignCustomOption parses a functions *ast.CommentGroups for custom options to return a Custom map.
func (p *Parser) assignCustomOption(o []*options.Option) map[string][]string {
	optionmap := make(map[string][]string)

	// functions only have custom options
	for i := 0; i < len(o); i++ {
		switch o[i].Category {
		case options.CategoryConvert, options.CategoryDeepCopy, options.CategoryDepth, options.CategoryMap:
		default:
			if customoptionmap, ok := o[i].Value.(map[string]string); ok {
				for customoption, value := range customoptionmap {
					optionmap[customoption] = append(optionmap[customoption], value)
				}
			} else if customoptionmap != nil {
				fmt.Printf("WARNING: Failed to assign custom option: %v\n", o[i].Value)
			}
		}
	}

	return optionmap
}
