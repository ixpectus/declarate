package shell

import (
	"strings"

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
	Cmd            string            `yaml:"shell_cmd,omitempty"`
	VariablesToSet map[string]string `yaml:"variables_to_set,omitempty"`
}

type Config struct {
	Cmd            string
	VariablesToSet map[string]string
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
	cfg := &shellConfig{}
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
			},
		}, nil
	}
	return &ShellCmd{
		Config: &Config{
			Cmd:            cfg.Cmd,
			VariablesToSet: cfg.VariablesToSet,
		},
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
	return nil
}

func (e *ShellCmd) VariablesToSet() map[string]string {
	if e.Config != nil {
		return e.Config.VariablesToSet
	}
	return nil
}
