package suite

type testDefinition struct {
	Definition *struct {
		Tags        []string `yaml:"tags,omitempty"`
		Condition   string   `yaml:"condition,omitempty"`
		Description string   `yaml:"description,omitempty"`
		ID          string   `yaml:"id,omitempty"`
	} `yaml:"definition,omitempty"`
}

type testWithDefinition struct {
	file       string
	definition testDefinition
}
