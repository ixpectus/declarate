package vars

import (
	"github.com/dailymotion/allure-go"
	"github.com/ixpectus/declarate/contract"
	"github.com/ixpectus/declarate/tools"
)

type VarsCmd struct {
	Config *Config
	Vars   contract.Vars
	report contract.ReportAttachement
}

type Config struct {
	Data           map[string]string `yaml:"variables,omitempty"`
	DataPersistent map[string]string `yaml:"variables_persistent,omitempty"`
}

func (e *VarsCmd) SetVars(vv contract.Vars) {
	e.Vars = vv
}

func (e *VarsCmd) SetReport(r contract.ReportAttachement) {
	e.report = r
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
		m := map[string]string{}
		for k, v := range e.Config.Data {
			m[k] = v
		}
		for k, v := range e.Config.DataPersistent {
			m[k] = v
		}
		res, _ := e.Vars.SetAll(m)
		if len(m) > 0 {
			e.report.AddAttachment("variables", allure.TextPlain, []byte(tools.FormatVariables(res)))
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
