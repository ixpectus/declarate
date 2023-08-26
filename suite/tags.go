package suite

type testDefinition struct {
	Definition *struct {
		Tags      []string `yaml:"tags,omitempty"`
		Condition string   `yaml:"condition,omitempty"`
	} `yaml:"definition,omitempty"`
}

type testWithDefinition struct {
	file       string
	definition testDefinition
}
