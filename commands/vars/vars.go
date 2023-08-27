package vars

import (
	"github.com/ixpectus/declarate/contract"
)

type VarsCmd struct {
	Config *Config
	Vars   contract.Vars
	eval   contract.Evaluator
}

type Config struct {
	Data           map[string]string `yaml:"variables,omitempty"`
	DataPersistent map[string]string `yaml:"variables_persistent,omitempty"`
}

func (e *VarsCmd) SetVars(vv contract.Vars) {
	e.Vars = vv
}

type Unmarshaller struct {
	eval contract.Evaluator
}

func NewUnmarshaller(eval contract.Evaluator) *Unmarshaller {
	return &Unmarshaller{
		eval: eval,
	}
}

func (u *Unmarshaller) Build(unmarshal func(interface{}) error) (contract.Doer, error) {
	cfg := &Config{}
	if err := unmarshal(cfg); err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}
	if cfg.Data == nil && cfg.DataPersistent == nil {
		return nil, nil
	}
	return &VarsCmd{
		Config: cfg,
	}, nil
}

func (e *VarsCmd) GetConfig() interface{} {
	return e.Config
}

func (e *VarsCmd) Do() error {
	if e.Config != nil {
		for k, v := range e.Config.Data {
			e.Vars.Set(k, v)
		}
		for k, v := range e.Config.DataPersistent {
			if err := e.Vars.SetPersistent(k, v); err != nil {
				panic(err)
			}
		}
	}

	return nil
}

func (e *VarsCmd) IsValid() error {
	return nil
}

func (e *VarsCmd) ResponseBody() *string {
	return nil
}

func (e *VarsCmd) Check() error {
	return nil
}

func (e *VarsCmd) VariablesToSet() map[string]string {
	return nil
}
