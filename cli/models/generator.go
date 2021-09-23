package models

// Generator represents a code generator.
type Generator struct {
	// Filepath Fields
	Filepath string   // The generated filepath.
	Loadpath string   // The loader filepath.
	Template Template // The template used to generate code.

	// Code Generation Fields
	Package   string     // The generated package.
	Imports   []string   // The imports included in the generated file.
	Functions []Function // The functions to generate.
}

type Template struct {
	Headpath string // The filepath to the template that generates header code.
	Funcpath string // The filepath to the template that generates function code.
}
