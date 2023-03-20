package vars

import (
	"github.com/ixpectus/declarate/contract"
)

type VarsCmd struct {
	Config *Config
	Vars   contract.Vars
}

type Config struct {
	Data map[string]string `yaml:"variables,omitempty"`
}

func (e *VarsCmd) SetVars(vv contract.Vars) {
	e.Vars = vv
}

func Build(unmarshal func(interface{}) error) (contract.Doer, error) {
	cfg := &Config{}
	if err := unmarshal(cfg); err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}
	return &VarsCmd{
		Config: cfg,
	}, nil
}

func (e *VarsCmd) Do() error {
	if e.Config != nil {
		for k, v := range e.Config.Data {
			e.Vars.Set(k, v)
		}
	}

	return nil
}

func (e *VarsCmd) ResponseBody() *string {
	return nil
}

func (e *VarsCmd) Check() error {
	return nil
}
