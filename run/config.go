package run

import (
	"fmt"

	"github.com/ixpectus/declarate/contract"
)

type runConfig struct {
	Name           string      `yaml:"name,omitempty"`
	Steps          []runConfig `yaml:"steps,omitempty"`
	Vars           contract.Vars
	VariablesToSet map[string]string `yaml:"variables_to_set"`
	Commands       []contract.Doer
	Builders       []contract.CommandBuilder
}

func (u *runConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type raw runConfig
	if err := unmarshal((*raw)(u)); err != nil {
		return err
	}
	u.Commands = []contract.Doer{}
	for _, v := range builders {
		b, err := v.Build(unmarshal)
		if err != nil {
			fmt.Printf("\n>>> %v <<< debug\n", err)
			return fmt.Errorf("unmarshal build err: %w", err)
		}
		if b != nil {
			u.Commands = append(u.Commands, b)
		}
	}
	return nil
}
