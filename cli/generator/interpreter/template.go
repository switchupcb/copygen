package interpreter

import (
	"fmt"

	"github.com/switchupcb/copygen/cli/generator/templates"
	"github.com/switchupcb/copygen/cli/models"
)

// Header determines the func to generate header code.
func Header(gen models.Generator) (string, error) {
	if gen.Template.Headpath == "" {
		return templates.DefaultHeader(gen), nil
	}
	return "", fmt.Errorf("Templates are temporarily unsupported.")
	return interpretHeader(gen)
}

// Function determines the func to generate function code.
func Function(gen *models.Generator) (string, error) {
	var functions string

	// determine the method to analyze each function.
	if gen.Template.Funcpath == "" {
		for _, function := range gen.Functions {
			functions += templates.DefaultFunction(function) + "\n"
		}
		return functions, nil
	}
	return "", fmt.Errorf("Templates are temporarily unsupported.")
	return interpretFunction(gen)
}

// interpretHeader represents the interpreted header func that generates the header code.
func interpretHeader(gen models.Generator) (string, error) {
	v, err := interpretFunc(gen.Loadpath, gen.Template.Headpath, "templates.Header")
	if err != nil {
		return "", err
	}
	// fn := v.Interface().(func(models.Generator) string)
	// header := fn(gen)
	fmt.Println(v)
	return "", nil
}

// interpretFunction represents the interpreted function func that generates function code.
func interpretFunction(gen *models.Generator) (string, error) {
	v, err := interpretFunc(gen.Loadpath, gen.Template.Funcpath, "templates.Function")
	if err != nil {
		return "", err
	}
	// fn := v.Interface().(func(models.Generator) string)
	// header := fn(gen)
	fmt.Println(v)
	return "", nil
}
