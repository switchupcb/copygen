package models

// Generator represents a code generator.
type Generator struct {
	Loadpath  string           // The loader filepath.
	Filepath  string           // The generated filepath.
	Template  Template         // The template used to generate code.
	Package   string           // The generated package.
	Imports   []string         // The imports included in the generated file.
	Functions []Function       // The functions to generate.
	Options   GeneratorOptions // The custom options for the generator.
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
