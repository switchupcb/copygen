package interpreter

import (
	"fmt"

	"github.com/switchupcb/copygen/cli/generator/template"
	"github.com/switchupcb/copygen/cli/models"
)

// Header determines the func to generate header code.
func Header(gen *models.Generator) (string, error) {
	if gen.Template.Headpath == "" {
		return template.DefaultHeader(gen), nil
	} else {
		return "", fmt.Errorf("Templates are temporarily unsupported.")
		// return interpretHeader(gen)
	}
}

// Function determines the func to generate function code.
func Function(gen *models.Generator) (string, error) {
	var functions string

	// determine the method to analyze each function.
	if gen.Template.Funcpath == "" {
		for _, function := range gen.Functions {
			functions += template.DefaultFunction(&function) + "\n"
		}
		return functions, nil
	}
	return "", fmt.Errorf("Templates are temporarily unsupported.")

	fn, err := interpretFunction(gen)
	if err != nil {
		return "", err
	}
	for _, function := range gen.Functions {
		functions += fn(&function) + "\n"
	}
	return functions, nil
}

// interpretHeader creates the header of the generated file using an interpreted template file.
func interpretHeader(gen *models.Generator) (string, error) {
	fn, err := interpretFunc(gen.Loadpath, gen.Template.Headpath, "generator.Header")
	if err != nil {
		return "", err
	}

	// run the interpreted function.
	return fn(gen), nil
}

// interpretFunction creates the header of the generated file using an interpreted template file.
func interpretFunction(gen *models.Generator) (func(f *models.Function) string, error) {
	fn, err := interpretFunc(gen.Loadpath, gen.Template.Funcpath, "generator.Function")
	if err != nil {
		return nil, err
	}

	// run the interpreted function.
	return func(function *models.Function) string {
		return fn(function)
	}, nil
}
