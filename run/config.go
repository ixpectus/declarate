package run

import (
	"fmt"

	"github.com/ixpectus/declarate/contract"
)

type runConfig struct {
	Name                string      `yaml:"name,omitempty"`
	Steps               []runConfig `yaml:"steps,omitempty"`
	Vars                contract.Vars
	Variables           map[string]string `yaml:"variables"`
	VariablesPersistent map[string]string `yaml:"variables_persistent"`
	Commands            []contract.Doer
	Builders            []contract.CommandBuilder
	Poll                *Poll  `yaml:"poll,omitempty"`
	Condition           string `yaml:"condition,omitempty"`
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
			return fmt.Errorf("unmarshal build err: %w", err)
		}
		if b != nil {
			u.Commands = append(u.Commands, b)
		}
	}

	return nil
}
