package parser

import (
	"fmt"
	"go/ast"
	"go/types"

	"github.com/switchupcb/copygen/cli/models"
	"github.com/switchupcb/copygen/cli/parser/options"
)

const copygenInterfaceName = "Copygen"

// parseFunctions parses the AST for functions in the setup file.
// astcopygen is used to assign options from *ast.Comments.
func (p *Parser) parseFunctions(astcopygen *ast.InterfaceType) ([]models.Function, error) {

	// find the `type Copygen interface` definition in the setup file.
	var copygen *types.Interface

	setpkg := p.Pkgs[0]
	defs := setpkg.TypesInfo.Defs
	for k, v := range defs {
		if k.Name == copygenInterfaceName {
			if it, ok := v.Type().Underlying().(*types.Interface); ok {
				copygen = it
			}
		}
	}

	if copygen == nil {
		return nil, fmt.Errorf("the \"type Copygen interface\" could not be found in the setup file's package")
	}

	if copygen.NumMethods() == 0 {
		return nil, fmt.Errorf("no functions are defined in the \"type Copygen interface\"")
	}

	// create the models.Function objects
	functions := make([]models.Function, copygen.NumMethods())
	for i := 0; i < copygen.NumMethods(); i++ {
		method := copygen.Method(i)

		// create the models.Type objects
		fieldoptions, manual := p.getNodeOptions(astcopygen.Methods.List[i])
		fieldoptions = append(fieldoptions, convertOptions...)
		parsed, err := parseTypes(method, fieldoptions)
		if err != nil {
			return nil, fmt.Errorf("an error occurred while parsing the types of function %q.\n%w", method.Name(), err)
		}

		// map the function custom options.
		var customoptionmap map[string][]string
		for _, option := range fieldoptions {
			err := options.MapCustomOption(customoptionmap, option)
			if err != nil {
				fmt.Printf("WARNING: %v", err)
			}
		}

		// create the models.Function
		function := models.Function{
			Name: method.Name(),
			To:   parsed.toTypes,
			From: parsed.fromTypes,
			Options: models.FunctionOptions{
				Custom: customoptionmap,
				Manual: manual,
			},
		}

		functions[i] = function
	}

	return functions, nil
}

// getNodeOptions gets an ast.Node options from its comments.
// To reduce overhead, it also returns whether a manual matcher is used.
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

				// specifying a map disables automatching.
				if p.CommentOptionMap[comment.Text].Category == options.CategoryMap {
					manual = true
				}
			}
		}

		return true
	})

	return nodeOptions, manual
}
