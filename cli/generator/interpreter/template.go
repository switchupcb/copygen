package interpreter

import (
	"fmt"

	"github.com/switchupcb/copygen/cli/generator/templates"
	"github.com/switchupcb/copygen/cli/models"
)

// Generate determines the func to generate function code.
func Generate(gen *models.Generator) (string, error) {
	var content string

	// determine the method to analyze each function.
	if gen.Tempath == "" {
		content += templates.Generate(gen) + "\n"

		return content, nil
	}

	return content, fmt.Errorf("templates are temporarily unsupported")
}

// interpretFunction represents the interpreted function func that generates function code.
func interpretFunction(gen *models.Generator) error {
	v, err := interpretFunc(gen.Loadpath, gen.Tempath, "templates.Generate")
	if err != nil {
		return err
	}

	fmt.Println(v)

	return nil
}
