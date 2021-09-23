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
		return nil, fmt.Errorf("The specified .yml filepath doesn't exist.")
	}

	var m YML
	err = yaml.Unmarshal(file, &m)
	if err != nil {
		return nil, fmt.Errorf("There is an issue with the provided .yml file.\n%v", err)
	}

	g := parseYML(m)
	return &g, nil
}

// parseYML parses a YML into a Generator.
func parseYML(m YML) models.Generator {
	var g models.Generator

	// define the generator options.
	g.GenFile = m.Generated["filepath"]
	g.GenPackage = m.Generated["package"]
	g.Imports = m.Import

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
				Pointer: toType.Pointer,
				Custom:  toType.Options,
			}

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
					Pointer: fromType.Pointer,
					Custom:  fromType.Options,
				}

				// define the fields of each type using the FromType.
				for fieldName, field := range fromType.Fields {
					// fromField
					var gfField models.Field
					gfField.Parent = gfType
					gfField.Name = fieldName
					gfField.Convert = field.Convert
					gfField.Options = models.FieldOptions{
						Custom: field.Options,
					}

					// toField
					var gtField models.Field
					gtField.Parent = gtType
					gtField.Name = field.To
					gtField.Convert = field.Convert
					gtField.Options = models.FieldOptions{
						Custom: field.Options,
					}

					// point the fields
					gfField.To = &gtField
					gfType.Fields = append(gfType.Fields, gfField)

					gtField.From = &gfField
					gtType.Fields = append(gtType.Fields, gtField)
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
	return g
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
