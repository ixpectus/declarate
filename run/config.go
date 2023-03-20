package run

import (
	"fmt"

	"github.com/ixpectus/declarate/contract"
)

type runConfig struct {
	Name     string      `yaml:"name,omitempty"`
	Steps    []runConfig `yaml:"steps,omitempty"`
	Commands []contract.Doer
	Builders []contract.CommandBuilder
	Vars     contract.Vars
}

func (u *runConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type raw runConfig
	if err := unmarshal((*raw)(u)); err != nil {
		return err
	}
	u.Commands = []contract.Doer{}
	for _, v := range builders {
		b, err := v(unmarshal)
		if err != nil {
			return fmt.Errorf("unmarshal build err: %w", err)
		}
		if b != nil {
			u.Commands = append(u.Commands, b)
		}
	}
	return nil
}
