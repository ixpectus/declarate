package run

import (
	"github.com/ixpectus/declarate/commands/echo"
	"github.com/ixpectus/declarate/commands/vars"
	"github.com/ixpectus/declarate/contract"
)

type RunConfig struct {
	Name     string      `yaml:"name,omitempty"`
	Steps    []RunConfig `yaml:"steps,omitempty"`
	Commands []contract.Doer
	Vars     contract.Vars
}

func (u *RunConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type raw RunConfig
	if err := unmarshal((*raw)(u)); err != nil {
		return err
	}
	e := &echo.Echo{}
	if err := e.FillData(unmarshal); err != nil {
		return err
	}
	u.Commands = append(u.Commands, e)
	vv := &vars.VarsCmd{}
	if err := vv.FillData(unmarshal); err != nil {
		return err
	}
	u.Commands = append(u.Commands, vv)

	return nil
}
