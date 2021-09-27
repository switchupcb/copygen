package template

import (
	"github.com/switchupcb/copygen/cli/generator/interpreter"
	"github.com/switchupcb/copygen/cli/models"
)

// Header determines the func to generate header code.
func Header(gen *models.Generator) (string, error) {
	if gen.Template.Headpath == "" {
		return defaultHeader(gen), nil
	} else {
		return interpretHeader(gen)
	}
}

// defaultHeader creates the header of the generated file using the default method.
func defaultHeader(gen *models.Generator) string {
	var header string

	// package
	header += "// Code generated by github.com/switchupcb/copygen\n"
	header += "// DO NOT EDIT.\n"
	header += "package " + gen.Package + "\n"

	// imports
	header += "import (\n"
	for _, iprt := range gen.Imports {
		header += "\"" + iprt + "\"\n"
	}
	header += ")"
	return header
}

// interpretHeader creates the header of the generated file using an interpreted template file.
func interpretHeader(gen *models.Generator) (string, error) {
	fn, err := interpreter.InterpretFunc(gen.Loadpath, gen.Template.Headpath, "generator.Header")
	if err != nil {
		return "", err
	}

	// run the interpreted function.
	return fn(gen), nil
}