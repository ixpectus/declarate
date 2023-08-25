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
	comparer     contract.Comparer
}

type extendedConfig struct {
	Shell *shellConfig `yaml:"shell,omitempty"`
}

type shellConfig struct {
	Cmd       string            `yaml:"cmd,omitempty"`
	Variables map[string]string `yaml:"variables,omitempty"`
	Response  *string           `yaml:"response,omitempty"`
}

type Config struct {
	Cmd       string            `yaml:"shell_cmd,omitempty"`
	Variables map[string]string `yaml:"variables,omitempty"`
	Response  *string           `yaml:"shell_response,omitempty"`
}

func (e *ShellCmd) SetVars(vv contract.Vars) {
	e.Vars = vv
}

func NewUnmarshaller(comparer contract.Comparer) *Unmarshaller {
	return &Unmarshaller{
		comparer: comparer,
	}
}

type Unmarshaller struct {
	comparer contract.Comparer
}

func (ex *extendedConfig) isEmpty() bool {
	return ex == nil || ex.Shell == nil || ex.Shell.Cmd == ""
}

func (c *Config) isEmpty() bool {
	return c == nil || c.Cmd == ""
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
	if cfg.isEmpty() && cfgExtended.isEmpty() {
		return nil, nil
	}
	if cfgExtended != nil && cfgExtended.Shell != nil {
		return &ShellCmd{
			comparer: u.comparer,
			Config: &Config{
				Cmd:       cfgExtended.Shell.Cmd,
				Variables: cfgExtended.Shell.Variables,
				Response:  cfgExtended.Shell.Response,
			},
		}, nil
	}
	return &ShellCmd{
		comparer: u.comparer,
		Config:   cfg,
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

func (e *ShellCmd) GetConfig() interface{} {
	return e.Config
}

func (e *ShellCmd) IsValid() error {
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
				*e.Config.Response,
				e.responseBody,
			)
			return res
		}
		errors := e.comparer.Compare(
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
		return e.Config.Variables
	}
	return nil
}
