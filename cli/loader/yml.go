// package loader loads generator information from an external file.
package loader

import (
	"fmt"
	"os"

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

	g := parseYMLMap(m)
	return &g, nil
}

// parseYMLMap parses a YMLMAP into a Generator.
func parseYMLMap(m YML) models.Generator {
	var g models.Generator

	// define the generator options.
	g.GenFile = m.Generated["filepath"]
	g.GenPackage = m.Generated["package"]
	g.Imports = m.Import

	// define the generator functions.
	for name, function := range m.Functions {
		var gf models.Function
		gf.Name = name
		gf.Options.Error = function.Error

		// define the To types of the function.
		var gtt models.Type
		for toName, toType := range function.To {
			gtt.Name = toName
			gtt.Filepath = toType.Filepath
			gtt.Options = models.TypeOptions{
				Pointer:  toType.Pointer,
				Deepcopy: toType.Deepcopy,
			}
			gf.To = append(gf.To, gtt)
		}

		// define the From types of the function.
		var gtf models.Type
		for fromName, fromType := range function.From {
			gtf.Name = fromName
			gtf.Filepath = fromType.Filepath
			gtf.Options = models.TypeOptions{
				Pointer: fromType.Pointer,
			}

			// define the fields of the From type.
			for fieldName, fieldOptions := range fromType.Fields {
				var gfdto models.Field
				gfdto.Name = fieldOptions["to"]
				gfdto.Convert = fieldOptions["convert"]

				var gfdfrom models.Field
				gfdfrom.Name = fieldName
				gfdfrom.Convert = fieldOptions["convert"]
				gfdfrom.To = &gfdto

				gtt.Fields = append(gtt.Fields, gfdto)
				gtf.Fields = append(gtf.Fields, gfdfrom)
			}

			gf.From = append(gf.From, gtf)
		}

		g.Functions = append(g.Functions, gf)
	}
	return g
}
