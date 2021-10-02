package models

// Generator represents a code generator.
type Generator struct {
	Loadpath  string            // The filepath the loader file is located in.
	Setpath   string            // The filepath the setup file is located in.
	Outpath   string            // The filepath the generated code is output to.
	Template  Template          // The template used to generate code.
	Imports   map[string]string // The imports to include in the generated file (map[packagealias]import).
	Functions []Function        // The functions to generate.
	Keep      []byte            // The code that is kept from the setup file (except for the header).
	Options   GeneratorOptions  // The custom options for the generator.
}

// GeneratorOptions represent options for a Generator.
type GeneratorOptions struct {
	Custom map[string]interface{} // The custom options of a generator.
}

// Template represets the template used to generate code.
type Template struct {
	Headpath string // The filepath to the template that generates header code.
	Funcpath string // The filepath to the template that generates function code.
}
