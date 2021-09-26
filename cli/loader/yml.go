// package loader loads generator information from an external file.
package loader

import (
	"fmt"
	"os"
	"strconv"

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

	g, err := parseYML(m)
	if err != nil {
		return nil, err
	}
	g.Loadpath = filepath
	return g, nil
}

// parseYML parses a YML into a Generator.
func parseYML(m YML) (*models.Generator, error) {
	var g models.Generator
	importMap := make(map[string]string) // a 'set' of imports.

	// define the generator options.
	if filepath, ok := m.Generated["filepath"].(string); ok {
		g.Filepath = filepath
	} else {
		return nil, fmt.Errorf("There is an issue with the .yml configuration for generated.filepath.")
	}

	if pkg, ok := m.Generated["package"].(string); ok {
		g.Package = pkg
	} else {
		return nil, fmt.Errorf("There is an issue with .yml configuration for generated.package.")
	}

	for _, imprt := range m.Import {
		importMap[imprt] = ""
	}

	g.Template = models.Template{
		Headpath: parseTemplate(m.Generated, "header"),
		Funcpath: parseTemplate(m.Generated, "function"),
	}

	// define the generator functions.
	for name, function := range m.Functions {
		var gf models.Function
		gf.Name = name

		// define the To types of the function.
		gtMap := make(map[string]bool) // A "set" of parameters
		for toName, toType := range function.To {
			var gtType models.Type
			varName := createVariable(gtMap, "t"+string(toName[0]), 0)
			gtMap[varName] = true
			gtType.Name = toName
			gtType.VariableName = varName
			gtType.Package = toType.Package
			gtType.Options = models.TypeOptions{
				Import:   toType.Import,
				Pointer:  toType.Pointer,
				Depth:    toType.Depth,
				Deepcopy: toType.Deepcopy,
				Custom:   toType.Options,
			}
			importMap[gtType.Options.Import] = ""

			// define the From types of the function.
			gfMap := make(map[string]bool) // A "set" of parameters
			for fromName, fromType := range function.From {
				var gfType models.Type
				varName := createVariable(gfMap, "f"+string(fromName[0]), 0)
				gfMap[varName] = true
				gfType.Name = fromName
				gfType.VariableName = varName
				gfType.Package = fromType.Package
				gfType.Options = models.TypeOptions{
					Import:   fromType.Import,
					Pointer:  fromType.Pointer,
					Depth:    fromType.Depth,
					Deepcopy: fromType.Deepcopy,
					Custom:   fromType.Options,
				}
				importMap[gfType.Options.Import] = ""

				// define the fields of each type using the FromType.
				var toFields, fromFields []models.Field
				if len(fromType.Fields) == 0 {
					var err error
					toFields, fromFields, err = Automatch(&gtType, &gfType)
					if err != nil {
						return nil, err
					}
				} else {
					toFields, fromFields = DefineFieldsByFromType(&fromType)
				}
				gtType.Fields = append(gtType.Fields, toFields...)
				gfType.Fields = append(gfType.Fields, fromFields...)
				for _, field := range gtType.Fields {
					if len(field.Fields) != 0 {
						field.Parent = gtType
					}
				}
				gfType.Fields = fromFields
				for _, field := range gfType.Fields {
					if len(field.Fields) != 0 {
						field.Parent = gtType
					}
				}
				gf.From = append(gf.From, gfType)
			}
			gf.To = append(gf.To, gtType)
		}
		gf.Options = models.FunctionOptions{
			Custom: function.Options,
		}
		g.Functions = append(g.Functions, gf)
	}
	for imprt := range importMap {
		if imprt != "" {
			g.Imports = append(g.Imports, imprt)
		}
	}
	return &g, nil

}

// parseTemplate parses a template map for a template key (option).
func parseTemplate(m map[string]interface{}, k string) string {
	if template, exists := m["templates"]; exists {
		if templateMap, ok := template.(map[string]interface{}); ok {
			if option, exists := templateMap[k]; exists {
				if value, ok := option.(string); ok {
					return value
				}
			}
		}
	}
	return ""
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
