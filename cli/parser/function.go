package parser

import (
	"fmt"
	"go/ast"
	"regexp"
	"strings"

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
		optionmap := make(map[string][]string)
		if method.Comment != nil {
			for _, comment := range method.Comment.List {
				option, value, err := parseComment(comment)
				if err != nil {
					return nil, fmt.Errorf("An error occurred while parsing the options of a function.\n%q", err)
				}
				optionmap[option] = append(optionmap[option], value)
			}
		}

		fromTypes, toTypes, err := p.parseTypes(method, optionmap)
		if err != nil {
			return nil, fmt.Errorf("An error occured while parsing the types of function %q.\n%v", parseMethodForName(method), err)
		}

		function := models.Function{
			Name: parseMethodForName(method),
			To:   toTypes,
			From: fromTypes,
		}
		for option, values := range optionmap {
			switch option {
			case "map":
			case "depth":
			case "deepcopy":
			default:
				// functions only have custom options
				function.Options.Custom[option] = values
			}
		}
		functions = append(functions, function)
	}
	return functions, nil
}

// parseComment parses a comment above a Copygen interface function to provide an option.
func parseComment(comment *ast.Comment) (string, string, error) {
	regex := regexp.MustCompile("\\s+")
	splitcomment := regex.Split(strings.TrimSpace(comment.Text), -1)
	if len(splitcomment) == 0 {
		return "", "", fmt.Errorf("There is an unspecified option at %v.\nPlease ensure there aren't empty lines between comments that pertain to a function.", comment.Slash)
	} else if len(splitcomment) != 2 {
		return "", "", fmt.Errorf("There is a misconfigured option at %v.\nIs it in format <option>:<whitespaces><value>?", comment.Slash)
	}
	return splitcomment[0], splitcomment[1], nil
}

// parseMethod parses a method inside of a Copygen interface to provide its name.
func parseMethodForName(method *ast.Field) string {
	var funcname string // i.e 'ModelsToDomain' in func ModelsToDomain(models.Account, *models.User) *domain.Account

	// ast Note: "Field.Names contains a single name "type" for elements of interface type lists"
	for _, name := range method.Names {
		funcname += name.String() // ModelsToDomain
	}
	return funcname
}
