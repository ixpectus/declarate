package vars

import (
	"sort"
	"strings"

	"github.com/ixpectus/declarate/contract"
)

type VarsCmd struct {
	Config *Config
	Vars   contract.Vars
	eval   contract.Evaluator
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
		keys := make([]string, 0, len(e.Config.Data))
		for k := range e.Config.Data {
			keys = append(keys, k)
		}
		sort.Slice(keys, func(i, j int) bool {
			iContains := strings.Contains(e.Config.Data[keys[i]], "{{")
			jContains := strings.Contains(e.Config.Data[keys[j]], "{{")
			if iContains && !jContains {
				return false
			}
			return true
		})
		for k := range keys {
			e.Vars.Set(keys[k], e.Config.Data[keys[k]])
		}

		keys = make([]string, 0, len(e.Config.DataPersistent))
		for k := range e.Config.DataPersistent {
			keys = append(keys, k)
		}
		sort.Slice(keys, func(i, j int) bool {
			iContains := strings.Contains(e.Config.DataPersistent[keys[i]], "{{")
			jContains := strings.Contains(e.Config.DataPersistent[keys[j]], "{{")
			if iContains && !jContains {
				return false
			}
			return true
		})
		for k := range keys {
			e.Vars.Set(keys[k], e.Config.DataPersistent[keys[k]])
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
