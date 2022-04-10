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
		return nil, fmt.Errorf("the specified .yml filepath doesn't exist: %v\n%w", filepath, err)
	}

	var yml YML
	if err := yaml.Unmarshal(file, &yml); err != nil {
		return nil, fmt.Errorf("an error occurred unmarshalling the .yml file\n%w", err)
	}

	gen := ParseYML(yml)

	gen.Loadpath = filepath

	return gen, nil
}

// ParseYML parses a YML into a Generator.
func ParseYML(yml YML) *models.Generator {
	return &models.Generator{
		Setpath: yml.Generated.Setup,
		Outpath: yml.Generated.Output,
		Tempath: yml.Generated.Template,
		Options: models.GeneratorOptions{
			Custom: yml.Options,
		},
	}
}
