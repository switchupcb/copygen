// Package config loads configuration data from an external file.
package config

import (
	"fmt"
	"os"

	"github.com/switchupcb/copygen/cli/models"
	"gopkg.in/yaml.v3"
)

// LoadYML loads a .yml configuration file into a Generator.
func LoadYML(filepath string) (*models.Generator, error) {
	file, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("the specified .yml filepath doesn't exist: %v\n%v", filepath, err)
	}

	var yml YML
	err = yaml.Unmarshal(file, &yml)
	if err != nil {
		return nil, err
	}

	gen := ParseYML(yml)
	gen.Loadpath = filepath
	return gen, nil
}

// ParseYML parses a YML into a Generator.
func ParseYML(yml YML) *models.Generator {
	gen := models.Generator{
		Setpath: yml.Generated.Setup,
		Outpath: yml.Generated.Output,
		Tempath: yml.Generated.Template,
		Options: models.GeneratorOptions{
			Custom: yml.Options,
		},
	}
	return &gen
}
