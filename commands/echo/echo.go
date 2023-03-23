package echo

import (
	"fmt"

	"github.com/ixpectus/declarate/contract"
)

type Echo struct {
	Config *EchoConfig
	Vars   contract.Vars
}
type Unmarshaller struct {
	host string
}

func (u *Unmarshaller) Build(unmarshal func(interface{}) error) (contract.Doer, error) {
	cfg := &EchoConfig{}
	if err := unmarshal(cfg); err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}
	return &Echo{
		Config: cfg,
	}, nil
}

type EchoConfig struct {
	Echo *struct {
		Message        string            `yaml:"message,omitempty"`
		Response       *string           `yaml:"response,omitempty"`
		VariablesToSet map[string]string `yaml:"variables_to_set"`
	} `yaml:"echo,omitempty"`
}

func (e *Echo) SetVars(vv contract.Vars) {
	e.Vars = vv
}

func (e *Echo) Do() error {
	if e != nil && e.Config != nil && e.Config.Echo != nil {
		e.Config.Echo.Message = e.Vars.Apply(e.Config.Echo.Message)
		fmt.Printf("\necho %v \n", e.Config.Echo.Message)
	}
	return nil
}

func (e *Echo) ResponseBody() *string {
	if e != nil && e.Config.Echo != nil {
		return e.Config.Echo.Response
	}
	return nil
}

func (e *Echo) VariablesToSet() map[string]string {
	if e != nil && e.Config.Echo != nil {
		return e.Config.Echo.VariablesToSet
	}
	return nil
}

func (e *Echo) Check() error {
	if e != nil && e.Config.Echo != nil {
		b := e.ResponseBody()
		if b != nil {
			if *b != e.Config.Echo.Message {
				return fmt.Errorf("expected %s, got %s", *b, e.Config.Echo.Message)
			}
		}
	}
	return nil
}
