package models

// Generator represents a code generator.
type Generator struct {
	Filepath  string     // The generated filepath.
	Loadpath  string     // The loader filepath.
	Package   string     // The generated package.
	Imports   []string   // The imports included in the generated file.
	Functions []Function // The functions to generate.
}
