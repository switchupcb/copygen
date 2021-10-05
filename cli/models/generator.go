package models

// Generator represents a code generator.
type Generator struct {
	Loadpath  string           // The filepath the loader file is located in.
	Setpath   string           // The filepath the setup file is located in.
	Outpath   string           // The filepath the generated code is output to.
	Tempath   string           // The filepath for thetemplate used to generate code.
	Keep      []byte           // The code that is kept from the setup file.
	Functions []Function       // The functions to generate.
	Options   GeneratorOptions // The custom options for the generator.
}

// GeneratorOptions represent options for a Generator.
type GeneratorOptions struct {
	Custom map[string]interface{} // The custom options of a generator.
}
