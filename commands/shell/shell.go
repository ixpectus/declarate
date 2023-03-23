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

type Config struct {
	Shell *ShellConfig `yaml:"shell,omitempty"`
}

type ShellConfig struct {
	Cmd            string            `json:"cmd,omitempty"`
	VariablesToSet map[string]string `yaml:"variables_to_set"`
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
	cfg := &Config{}
	if err := unmarshal(cfg); err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}
	return &ShellCmd{
		Config: cfg,
	}, nil
}

func (e *ShellCmd) Do() error {
	if e.Config != nil && e.Config.Shell != nil && e.Config.Shell.Cmd != "" {
		res, err := Run(e.Config.Shell.Cmd)
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
	if e.Config.Shell != nil {
		return e.Config.Shell.VariablesToSet
	}
	return nil
}
