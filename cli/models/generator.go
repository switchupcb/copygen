package models

import (
	"go/ast"
	"go/token"
)

// Generator represents a code generator.
type Generator struct {
	Functions []Function       // The functions to generate.
	Options   GeneratorOptions // The custom options for the generator.
	Loadpath  string           // The filepath the loader file is located in.
	Setpath   string           // The filepath the setup file is located in.
	Outpath   string           // The filepath the generated code is output to.
	Tempath   string           // The filepath for the template used to generate code.
	Keep      []byte           // The code that is kept from the setup file.

	ImportsByName   map[string]string // Map of imports to its alias.
	ImportsByPath   map[string]string // Map of imports to its alias.
	AlreadyImported map[string]bool   // Map of imports to its alias.
	// The fileset of the parser.
	Fileset *token.FileSet

	// The setup file as an Abstract Syntax Tree.
	SetupFile *ast.File
}

// GeneratorOptions represent options for a Generator.
type GeneratorOptions struct {
	Custom map[string]interface{} // The custom options of a generator.
}
