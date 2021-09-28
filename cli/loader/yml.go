// Package loader loads generator information from an external file.
package loader

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/switchupcb/copygen/cli/models"
	"gopkg.in/yaml.v3"
)

// Parser represents a YML parser that loads properties into the program models.
type Parser struct {
	YML       YML              // The YML that is parsed.
	AST       AST              // The Abstract Syntax Tree object used during matching.
	Generator models.Generator // The generator that information is parsed to.

}

// LoadYML loads a .yml file into a Generator.
func LoadYML(filepath string) (*models.Generator, error) {
	file, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("The specified .yml filepath doesn't exist: %v.\n%v", filepath, err)
	}

	var p Parser
	err = yaml.Unmarshal(file, &p.YML)
	if err != nil {
		return nil, fmt.Errorf("There is an issue with the provided .yml file: %v\n%v", filepath, err)
	}

	gen, err := p.ParseYML()
	if err != nil {
		return nil, err
	}
	gen.Loadpath = filepath
	return gen, nil
}

// ParseYML parses a YML into a Generator.
func (p *Parser) ParseYML() (*models.Generator, error) {
	// define the generator options.
	if filepath, ok := p.YML.Generated["filepath"].(string); ok {
		p.Generator.Filepath = filepath
	} else {
		return nil, fmt.Errorf("There is an issue with the .yml configuration for generated.filepath.")
	}

	if pkg, ok := p.YML.Generated["package"].(string); ok {
		p.Generator.Package = pkg
	} else {
		return nil, fmt.Errorf("There is an issue with .yml configuration for generated.package.")
	}

	p.Generator.Template = models.Template{
		Headpath: p.parseTemplate("header"),
		Funcpath: p.parseTemplate("function"),
	}

	// define the generator functions
	for name := range p.YML.Functions {
		modelFunction, err := p.parseFunction(name)
		if err != nil {
			return nil, err
		}
		p.Generator.Functions = append(p.Generator.Functions, *modelFunction)
	}

	for imprt := range p.parseImports() {
		p.Generator.Imports = append(p.Generator.Imports, imprt)
	}
	return &p.Generator, nil
}

// parseTemplate parses a template map for a template key (option).
func (p *Parser) parseTemplate(k string) string {
	if template, exists := p.YML.Generated["templates"]; exists {
		if templateMap, ok := template.(map[string]interface{}); ok {
			if option, exists := templateMap[k]; exists {
				if str, ok := option.(string); ok {
					return str
				}
			}
		}
	}
	return ""
}

// parseImports parses a generator's objects to determine its imports.
func (p *Parser) parseImports() map[string]bool {
	imprtMap := make(map[string]bool) // a 'set' of imports.
	for _, imprt := range p.YML.Import {
		if imprt != "" {
			imprtMap[strings.TrimSpace(imprt)] = true
		}
	}
	for _, function := range p.Generator.Functions {
		for _, toType := range function.To {
			if toType.Options.Import != "" {
				imprtMap[strings.TrimSpace(toType.Options.Import)] = true
			}
		}
		for _, fromType := range function.From {
			if fromType.Options.Import != "" {
				imprtMap[strings.TrimSpace(fromType.Options.Import)] = true
			}
		}
	}
	return imprtMap
}

// parseFunction parses a YML function.
func (p *Parser) parseFunction(name string) (*models.Function, error) {
	var function models.Function
	function.Name = name

	// define the To types of the function.
	toParams := make(map[string]bool) // A "set" of parameters
	for toname, to := range p.YML.Functions[name].To {
		tovarname := createVariable(toParams, "t"+string(toname[0]), 0)
		toParams[tovarname] = true
		toType := models.Type{
			Name:         toname,
			VariableName: tovarname,
			Package:      to.Package,
			Options: models.TypeOptions{
				Import:   to.Import,
				Pointer:  to.Pointer,
				Depth:    to.Depth,
				Deepcopy: to.Deepcopy,
				Custom:   to.Options,
			},
		}

		// define the From types of the function.
		fromParams := make(map[string]bool) // A "set" of parameters
		for fromname, from := range p.YML.Functions[name].From {
			fromvarname := createVariable(fromParams, "f"+string(fromname[0]), 0)
			fromParams[fromvarname] = true
			fromType := models.Type{
				Name:         fromname,
				VariableName: fromvarname,
				Package:      from.Package,
				Options: models.TypeOptions{
					Import:   from.Import,
					Pointer:  from.Pointer,
					Depth:    from.Depth,
					Deepcopy: from.Deepcopy,
					Custom:   from.Options,
				},
			}

			// determine the fields of a from type
			toFields, fromFields, err := p.parseFields(from, &toType, &fromType)
			if err != nil {
				return nil, err
			}

			// assign the fields
			fromType.Fields = fromFields
			for i := 0; i < len(toFields); i++ {
				toType.Fields = append(toType.Fields, toFields[i])
			}
			function.From = append(function.From, fromType)
		}
		function.To = append(function.To, toType)
	}
	function.Options = models.FunctionOptions{
		Custom: p.YML.Functions[name].Options,
	}
	return &function, nil
}

// parseFields parses the fields of two types.
func (p *Parser) parseFields(from From, toType *models.Type, fromType *models.Type) ([]*models.Field, []*models.Field, error) {
	if len(from.Fields) == 0 {
		var err error
		toFields, fromFields, err := p.AST.Automatch(toType, fromType)
		if err != nil {
			return nil, nil, err
		}
		return toFields, fromFields, nil
	}
	// otherwise use the match-by-hand method
	toFields, fromFields := DefineFieldsByFrom(&from, toType, fromType)
	return toFields, fromFields, nil
}

// createVariable p.Generatorerates a valid variable name for a list of parameters.
func createVariable(parameters map[string]bool, typename string, occurrence int) string {
	if occurrence < 0 {
		createVariable(parameters, typename, 0)
	}

	varName := typename
	if occurrence > 0 {
		varName += strconv.Itoa(occurrence + 1)
	}

	if _, exists := parameters[varName]; exists {
		createVariable(parameters, typename, occurrence+1)
	}
	return varName
}
