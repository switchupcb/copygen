package loader

// YML represents the first level of the YML file.
type YML struct {
	Generated map[string]interface{} `yaml:"generated"`
	Import    []string               `yaml:"import"`
	Functions map[string]Function    `yaml:"functions"`
}

// Function represents the function level of the YML file.
type Function struct {
	To      map[string]To          `yaml:"to"`
	From    map[string]From        `yaml:"from"`
	Options map[string]interface{} `yaml:"options"`
}

// To represents the to-type in the YML file.
type To struct {
	Package string                 `yaml:"package"`
	Pointer bool                   `yaml:"pointer"`
	Options map[string]interface{} `yaml:"options"`
}

// From represents the from-type in the YML file.
type From struct {
	Package string                 `yaml:"package"`
	Pointer bool                   `yaml:"pointer"`
	Fields  map[string]Field       `yaml:"fields"`
	Options map[string]interface{} `yaml:"options"`
}

// Field represents the field options of the YML file.
type Field struct {
	To      string                 `yaml:"to"`
	Convert string                 `yaml:"convert"`
	Options map[string]interface{} `yaml:"options"`
}
