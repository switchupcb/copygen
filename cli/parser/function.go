package parser

import (
	"fmt"
	"go/ast"

	"github.com/switchupcb/copygen/cli/models"
)

// parseFunctions parses the AST for functions in the setup file.
func (p *Parser) parseFunctions() ([]models.Function, error) {
	typecopygen, err := astTypeSearch(p.SetupFile, "Copygen")
	if err != nil {
		return nil, fmt.Errorf("The \"type Copygen interface\" could not be found in the setup file.")
	}

	copyinterface, ok := typecopygen.Type.(*ast.InterfaceType)
	if !ok {
		return nil, fmt.Errorf("The \"type Copygen\" was found but its not an interface. Please redefine it.")
	}

	var functions []models.Function
	for _, method := range copyinterface.Methods.List {
		fromTypes, toTypes, err := p.parseTypes(method)
		if err != nil {
			return nil, fmt.Errorf("An error occured while parsing the types of function %q.\n%v", parseMethodForName(method), err)
		}

		function := models.Function{
			Name: parseMethodForName(method),
			To:   toTypes,
			From: fromTypes,
		}

		function.Options.Custom = p.assignCustomOption(method)
		functions = append(functions, function)
	}
	return functions, nil
}

// assignCustomOption parses a functions *ast.CommentGroups for custom options to return a Custom map.
func (p *Parser) assignCustomOption(x ast.Node) map[string][]string {
	optionmap := make(map[string][]string)
	ast.Inspect(x, func(node ast.Node) bool {
		switch xcg := node.(type) {
		case *ast.CommentGroup:
			for _, comment := range xcg.List {
				// functions only have custom options
				if option, exists := p.Options[comment.Text]; exists && option.Category == "custom" {
					optionvalue := option.Value[0]
					if customoptionmap, ok := optionvalue.(map[string]string); ok {
						for customoption, value := range customoptionmap {
							optionmap[customoption] = append(optionmap[customoption], value)
						}
					} else {
						fmt.Println("WARNING: Failed to assign custom option.", optionvalue)
					}
				}
			}
		}
		return true
	})
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
