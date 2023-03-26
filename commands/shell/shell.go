package shell

import (
	"fmt"
	"strings"

	"github.com/ixpectus/declarate/compare"
	"github.com/ixpectus/declarate/contract"
)

type ShellCmd struct {
	Config       *Config
	Vars         contract.Vars
	responseBody string
}

type extendedConfig struct {
	Shell *shellConfig `yaml:"shell,omitempty"`
}

type shellConfig struct {
	Cmd            string            `yaml:"cmd,omitempty"`
	VariablesToSet map[string]string `yaml:"variables_to_set,omitempty"`
	Response       *string           `yaml:"response,omitempty"`
}

type Config struct {
	Cmd            string            `yaml:"shell_cmd,omitempty"`
	VariablesToSet map[string]string `yaml:"variables_to_set,omitempty"`
	Response       *string           `yaml:"shell_response,omitempty"`
}

func (e *ShellCmd) SetVars(vv contract.Vars) {
	e.Vars = vv
}

func NewUnmarshaller() *Unmarshaller {
	return &Unmarshaller{}
}

type Unmarshaller struct {
	host string
}

func (u *Unmarshaller) Build(unmarshal func(interface{}) error) (contract.Doer, error) {
	cfgExtended := &extendedConfig{}
	if err := unmarshal(cfgExtended); err != nil {
		return nil, err
	}
	cfg := &Config{}
	if err := unmarshal(cfg); err != nil {
		return nil, err
	}
	if cfg == nil && cfgExtended == nil {
		return nil, nil
	}
	if cfgExtended != nil && cfgExtended.Shell != nil {
		return &ShellCmd{
			Config: &Config{
				Cmd:            cfgExtended.Shell.Cmd,
				VariablesToSet: cfgExtended.Shell.VariablesToSet,
				Response:       cfgExtended.Shell.Response,
			},
		}, nil
	}
	return &ShellCmd{
		Config: cfg,
	}, nil
}

func (e *ShellCmd) Do() error {
	if e.Config != nil && e.Config.Cmd != "" {
		e.Config.Cmd = e.Vars.Apply(e.Config.Cmd)
		res, err := Run(e.Config.Cmd)
		if err != nil {
			return err
		}
		e.responseBody = strings.Join(res, "\n")

		return nil
	}

	return nil
}

func (e *ShellCmd) ResponseBody() *string {
	return &e.responseBody
}

func (e *ShellCmd) Check() error {
	if e.Config.Response != nil {
		linesExpected := strings.Split(*e.Config.Response, "\n")
		linesGot := strings.Split(e.responseBody, "\n")
		if len(linesExpected) != len(linesGot) {
			res := compare.MakeError(
				"",
				fmt.Sprintf("lines count differs, expected %v, got %v", len(linesExpected), len(linesGot)),
				e.responseBody,
				*e.Config.Response,
			)
			return res
		}
		errors := compare.Compare(
			*e.Config.Response,
			e.responseBody,
			compare.CompareParams{},
		)
		if len(errors) > 0 {
			return errors[0]
		}
	}
	return nil
}

func (e *ShellCmd) VariablesToSet() map[string]string {
	if e.Config != nil {
		return e.Config.VariablesToSet
	}
	return nil
}
