package loader

// YML represents the first level of the YML file.
type YML struct {
	Generated map[string]string   `yaml:"generated"`
	Import    []string            `yaml:"import"`
	Functions map[string]Function `yaml:"functions"`
}

// Function represents the function level of the YML file.
type Function struct {
	To    map[string]To   `yaml:"to"`
	From  map[string]From `yaml:"from"`
	Error bool            `yaml:"error"`
}

// To represents the to-type in the YML file.
type To struct {
	Filepath string `yaml:"filepath"`
	Pointer  bool   `yaml:"pointer"`
	Deepcopy bool   `yaml:"deepcopy"`
}

// From represents the from-type in the YML file.
type From struct {
	Filepath string                  `yaml:"filepath"`
	Pointer  bool                    `yaml:"pointer"`
	Fields   map[string]FieldOptions `yaml:"fields"`
}

// Field represents the field options of the YML file.
type FieldOptions map[string]string
