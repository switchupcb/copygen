// Package config loads configuration data from an external file.
package config

// YML represents the first level of the YML file.
type YML struct {
	Options   map[string]interface{} `yaml:"custom"`
	Generated Generated              `yaml:"generated"`
}

// Generated represents generated properties of the YML file.
type Generated struct {
	Setup    string `yaml:"setup"`
	Output   string `yaml:"output"`
	Template string `yaml:"template"`
}
