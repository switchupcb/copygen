// Package config loads configuration data from an external file.
package config

// YML represents the first level of the YML file.
type YML struct {
	Options   map[string]interface{} `yaml:"custom"`
	Generated Generated              `yaml:"generated"`
	Matcher   Matcher                `yaml:"matcher"`
}

// Generated represents generated properties of the YML file.
type Generated struct {
	Setup    string `yaml:"setup"`
	Output   string `yaml:"output"`
	Template string `yaml:"template"`
}

// Matcher represents matcher properties of the YML file.
type Matcher struct {
	Skip bool `yaml:"skip"`
	Cast Cast `yaml:"cast"`
}

// Cast represents matcher cast properties of the YML file.
type Cast struct {
	Depth    int      `yaml:"depth"`
	Enabled  bool     `yaml:"enabled"`
	Disabled Disabled `yaml:"disabled"`
}

// Disabled represents matcher cast feature flags of the YML file.
type Disabled struct {
	AssignObjectInterface bool `yaml:"assignObjectInterface"`
	AssertInterfaceObject bool `yaml:"assertInterfaceObject"`
	Convert               bool `yaml:"convert"`
}
