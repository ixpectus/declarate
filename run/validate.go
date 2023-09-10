package run

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

func (r *Runner) Validate(fileName string) error {
	file, err := os.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("file open: %w", err)
	}
	currentVars = r.config.Variables
	configs := []runConfig{}
	if err := yaml.Unmarshal(file, &configs); err != nil {
		return fmt.Errorf("unmarshall failed for file %s: %w", r.filenameShort(fileName), err)
	}
	for _, v := range configs {
		for _, c := range v.Commands {
			if err := c.IsValid(); err != nil {
				return fmt.Errorf("invalid command, %w", err)
			}
		}
	}

	return nil
}
