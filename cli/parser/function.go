package parser

import (
	"fmt"
	"go/ast"
	"go/types"

	"github.com/switchupcb/copygen/cli/models"
	"github.com/switchupcb/copygen/cli/parser/options"
)

// parseFunctions parses the AST for functions in the setup file.
// astcopygen is used to assign options from *ast.Comments.
func (p *Parser) parseFunctions(copygen *ast.InterfaceType) ([]models.Function, error) {
	numMethods := len(copygen.Methods.List)
	if numMethods == 0 {
		fmt.Println("WARNING: no functions are defined in the \"type Copygen interface\"")
	}

	// create models.Function objects.
	functions := make([]models.Function, numMethods)
	for i := 0; i < numMethods; i++ {
		method := p.Config.SetupPkg.TypesInfo.Defs[copygen.Methods.List[i].Names[0]]

		// create models.Type objects.
		fieldoptions, manual := getNodeOptions(copygen.Methods.List[i], p.Options.CommentOptionMap)
		fieldoptions = append(fieldoptions, p.Options.ConvertOptions...)
		parsed, err := parseTypes(method.(*types.Func))
		if err != nil {
			return nil, fmt.Errorf("an error occurred while parsing the types of function %q.\n%w", method.Name(), err)
		}

		// set the options for each field.
		setTypeOptions(parsed.fromTypes, fieldoptions, method.Name())
		setTypeOptions(parsed.toTypes, fieldoptions, method.Name())

		// map the function custom options.
		customoptionmap := make(map[string][]string)
		for _, option := range fieldoptions {
			customoptionmap, err = options.MapCustomOption(customoptionmap, option)
			if err != nil {
				fmt.Printf("WARNING: %v\n", err)
			}
		}

		// create the models.Function object.
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
func getNodeOptions(x ast.Node, commentoptionmap map[string]*options.Option) ([]*options.Option, bool) {
	nodeOptions := make([]*options.Option, 0, len(commentoptionmap))
	var manual bool

	ast.Inspect(x, func(node ast.Node) bool {
		commentGroup, ok := node.(*ast.CommentGroup)
		if !ok {
			return true
		}

		for _, comment := range commentGroup.List {
			if commentoptionmap[comment.Text] != nil {
				nodeOptions = append(nodeOptions, commentoptionmap[comment.Text])

				// specifying a match option disables automatching by default.
				if options.IsMatchOptionCategory(commentoptionmap[comment.Text].Category) {
					manual = true
				}
			}
		}

		return true
	})

	return nodeOptions, manual
}

// setTypeOptions sets the options for all fields in the given types.
func setTypeOptions(types []models.Type, fieldoptions []*options.Option, functionName string) {
	for _, t := range types {
		for _, field := range t.Field.AllFields(nil, nil) {
			options.SetFieldOptions(field, fieldoptions, functionName)
			options.FilterDepth(field, field.Options.Depth, 0)
		}
	}
}
