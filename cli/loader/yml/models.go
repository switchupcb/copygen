// Package yml loads setup information from an external file.
package yml

// YML represents the first level of the YML file.
type YML struct {
	Generated Generated `yaml:"generated"`
	Templates Templates `yaml:"templates"`
}

// Generated represents generated properties of the YML file.
type Generated struct {
	Setup   string `yaml:"setup"`
	Output  string `yaml:"output"`
	Package string `yaml:"package"`
}

// Templates represent template properties of the YML file.
type Templates struct {
	Header   string `yaml:"header"`
	Function string `yaml:"function"`
}
