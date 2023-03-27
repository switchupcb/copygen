package models

// Generator represents a code generator.
type Generator struct {
	Functions []Function       // The functions to generate.
	Options   GeneratorOptions // The custom options for the generator.
	Setpath   string           // The filepath the setup file is located in.
	Outpath   string           // The filepath the generated code is output to.
	Tempath   string           // The filepath for the template used to generate code.
	Keep      []byte           // The code that is kept from the setup file.
}

// GeneratorOptions represents options for a Generator.
type GeneratorOptions struct {
	Custom  map[string]interface{} // The custom options of a generator.
	Matcher MatcherOptions         // The options for the matcher of a generator.
}

// MatcherOptions represents options for the Generator's matcher.
type MatcherOptions struct {
	CastDepth                    int  // The option that sets the maximum depth for automatic casting.
	Skip                         bool // The option that skips the matcher.
	AutoCast                     bool // The option that enables automatic casting.
	DisableAssignObjectInterface bool // The cast option feature flag that disables assignment of objects to interfaces.
	DisableAssertInterfaceObject bool // The cast option feature flag that disables assignment of interfaces to objects.
	DisableConvert               bool // The cast option feature flag that disables type conversion.
}
