package echo

import (
	"fmt"

	"github.com/ixpectus/declarate/contract"
)

type Echo struct {
	Config *EchoConfig
	Vars   contract.Vars
}

func (e *Echo) FillData(unmarshal func(interface{}) error) error {
	cfg := &EchoConfig{}
	if err := unmarshal(cfg); err != nil {
		return err
	}
	e.Config = cfg
	return nil
}

type EchoConfig struct {
	Echo *struct {
		Message  string  `yaml:"message,omitempty"`
		Response *string `yaml:"response,omitempty"`
	} `yaml:"echo,omitempty"`
}

func (e *Echo) SetVars(vv contract.Vars) {
	e.Vars = vv
}

func (e *Echo) Do() error {
	if e != nil && e.Config != nil && e.Config.Echo != nil {
		e.Config.Echo.Message = e.Vars.Apply(e.Config.Echo.Message)
		fmt.Printf("\n>>> %v <<< debug\n", e.Config.Echo.Message)
	}
	return nil
}

func (e *Echo) ResponseBody() *string {
	if e != nil && e.Config.Echo != nil {
		return e.Config.Echo.Response
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
