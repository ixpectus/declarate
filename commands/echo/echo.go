package echo

import (
	"fmt"

	"github.com/ixpectus/declarate/contract"
)

type Echo struct {
	Config *Config
	Vars   contract.Vars
	Report contract.ReportAttachement
}

func (ex *extendedConfig) isEmpty() bool {
	return ex == nil || ex.Echo == nil || ex.Echo.Message == ""
}

func (c *config) isEmpty() bool {
	return c == nil || c.Message == ""
}

type Unmarshaller struct {
	host string
}

func (u *Unmarshaller) Build(unmarshal func(interface{}) error) (contract.Doer, error) {
	cfgExtended := &extendedConfig{}
	if err := unmarshal(cfgExtended); err != nil {
		return nil, err
	}
	cfgShort := &config{}
	if err := unmarshal(cfgShort); err != nil {
		return nil, err
	}
	if cfgExtended.isEmpty() && cfgShort.isEmpty() {
		return nil, nil
	}
	if !cfgExtended.isEmpty() {
		return &Echo{
			Config: &Config{
				Message:        cfgExtended.Echo.Message,
				Response:       cfgExtended.Echo.Response,
				VariablesToSet: cfgExtended.Echo.VariablesToSet,
			},
		}, nil
	}

	return &Echo{
		Config: &Config{
			Message:        cfgShort.Message,
			Response:       cfgShort.Response,
			VariablesToSet: cfgShort.VariablesToSet,
		},
	}, nil
}

type Config struct {
	Message        string
	Response       *string
	VariablesToSet map[string]string
}

func (c *Config) isEmpty() bool {
	return c == nil || c.Message == ""
}

type config struct {
	Message        string            `yaml:"echo_message,omitempty"`
	Response       *string           `yaml:"echo_response,omitempty"`
	VariablesToSet map[string]string `yaml:"variables"`
}

type extendedConfig struct {
	Echo *struct {
		Message        string            `yaml:"message,omitempty"`
		Response       *string           `yaml:"response,omitempty"`
		VariablesToSet map[string]string `yaml:"variables"`
	} `yaml:"echo,omitempty"`
}

func (e *Echo) SetVars(vv contract.Vars) {
	e.Vars = vv
}

func (e *Echo) SetReport(r contract.ReportAttachement) {
	e.Report = r
}

func (e *Echo) Do() error {
	if !e.Config.isEmpty() {
		e.Config.Message = e.Vars.Apply(e.Config.Message)
		fmt.Printf("\necho %v \n", e.Config.Message)
		return nil
	}
	return nil
}

func (e *Echo) GetConfig() interface{} {
	return e.Config
}

func (e *Echo) ResponseBody() *string {
	if !e.Config.isEmpty() {
		return e.Config.Response
	}
	return nil
}

func (e *Echo) ExpectedResponse() string {
	if !e.Config.isEmpty() {
		return e.Config.Message
	}
	return ""
}

func (e *Echo) IsValid() error {
	return nil
}

func (e *Echo) Variables() map[string]string {
	if e != nil && e.Config != nil {
		return e.Config.VariablesToSet
	}
	return nil
}

func (e *Echo) Check() error {
	if e != nil && !e.Config.isEmpty() {
		b := e.ResponseBody()
		if b != nil {
			if *b != e.Config.Message {
				return fmt.Errorf("expected %s, got %s", *b, e.Config.Message)
			}
		}
		return nil
	}
	return nil
}
