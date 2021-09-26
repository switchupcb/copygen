package loader

// YML represents the first level of the YML file.
type YML struct {
	Generated map[string]interface{} `yaml:"generated"`
	Import    []string               `yaml:"import"`
	Functions map[string]Function    `yaml:"functions"`
}

// Function represents function properties of the YML file.
type Function struct {
	To      map[string]To          `yaml:"to"`
	From    map[string]From        `yaml:"from"`
	Options map[string]interface{} `yaml:"options"`
}

// To represents to-type properties in the YML file.
type To struct {
	Package  string                 `yaml:"package"`
	Import   string                 `yaml:"import"`
	Pointer  bool                   `yaml:"pointer"`
	Depth    int                    `yaml:"depth"`
	Deepcopy string                 `yaml:"deepcopy"`
	Options  map[string]interface{} `yaml:"options"`
}

// From represents from-type properties in the YML file.
type From struct {
	Package  string                 `yaml:"package"`
	Import   string                 `yaml:"import"`
	Fields   map[string]Field       `yaml:"fields"`
	Pointer  bool                   `yaml:"pointer"`
	Depth    int                    `yaml:"depth"`
	Deepcopy string                 `yaml:"deepcopy"`
	Options  map[string]interface{} `yaml:"options"`
}

// Field represents the field properties of the YML file.
type Field struct {
	To       string                 `yaml:"to"`
	Convert  string                 `yaml:"convert"`
	Deepcopy string                 `yaml:"deepcopy"`
	Options  map[string]interface{} `yaml:"options"`
}
