// package loader loads generator information from an external file.
package loader

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/switchupcb/copygen/cli/models"
	"gopkg.in/yaml.v3"
)

// YML loads a .yml file into a Generator.
func LoadYML(filepath string) (*models.Generator, error) {
	file, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("The specified .yml filepath doesn't exist: %v.\n%v", filepath, err)
	}

	var m YML
	err = yaml.Unmarshal(file, &m)
	if err != nil {
		return nil, fmt.Errorf("There is an issue with the provided .yml file: %v\n%v", filepath, err)
	}

	gen, err := ParseYML(m)
	if err != nil {
		return nil, err
	}
	gen.Loadpath = filepath
	return gen, nil
}

// ParseYML parses a YML into a Generator.
func ParseYML(m YML) (*models.Generator, error) {
	var gen models.Generator

	// define the generator options.
	if filepath, ok := m.Generated["filepath"].(string); ok {
		gen.Filepath = filepath
	} else {
		return nil, fmt.Errorf("There is an issue with the .yml configuration for generated.filepath.")
	}

	if pkg, ok := m.Generated["package"].(string); ok {
		gen.Package = pkg
	} else {
		return nil, fmt.Errorf("There is an issue with .yml configuration for generated.package.")
	}

	gen.Template = models.Template{
		Headpath: parseTemplate(m.Generated, "header"),
		Funcpath: parseTemplate(m.Generated, "function"),
	}

	// define the generator functions
	a := AST{}
	for name, function := range m.Functions {
		modelFunction, err := parseFunction(function, name, &a)
		if err != nil {
			return nil, err
		}
		gen.Functions = append(gen.Functions, *modelFunction)
	}

	imprtMap := make(map[string]bool) // a 'set' of imports.
	for _, imprt := range m.Import {
		if imprt != "" {
			imprtMap[strings.TrimSpace(imprt)] = true
		}
	}
	for _, function := range gen.Functions {
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
	for imprt := range imprtMap {
		gen.Imports = append(gen.Imports, imprt)
	}
	return &gen, nil
}

// parseTemplate parses a template map for a template key (option).
func parseTemplate(m map[string]interface{}, k string) string {
	if template, exists := m["templates"]; exists {
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

// parseFunction parses a YML function.
func parseFunction(f Function, name string, a *AST) (*models.Function, error) {
	var function models.Function
	function.Name = name

	// define the To types of the function.
	toParams := make(map[string]bool) // A "set" of parameters
	for toname, to := range f.To {
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
		for fromname, from := range f.From {
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

			// define the fields of a from type
			toFields, fromFields, err := parseFields(from, &toType, &fromType, a)
			if err != nil {
				return nil, err
			}
			fromType.Fields = fromFields
			toType.Fields = append(toType.Fields, toFields...)
			function.From = append(function.From, fromType)
		}
		function.To = append(function.To, toType)
	}
	function.Options = models.FunctionOptions{
		Custom: f.Options,
	}
	return &function, nil
}

// parseFields parses the fields of two types.
func parseFields(from From, toType *models.Type, fromType *models.Type, a *AST) ([]models.Field, []models.Field, error) {
	if len(from.Fields) == 0 {
		var err error
		toFields, fromFields, err := a.Automatch(toType, fromType)
		if err != nil {
			return nil, nil, err
		}
		return toFields, fromFields, nil
	} else {
		toFields, fromFields := DefineFieldsByFrom(&from, toType, fromType)
		return toFields, fromFields, nil
	}
}

// createVariable generates a valid variable name for a list of parameters.
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
