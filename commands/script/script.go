package script

import (
	"fmt"

	"github.com/ixpectus/declarate/contract"
)

type ScriptCmd struct {
	Config       *Config
	Vars         contract.Vars
	report       contract.ReportAttachement
	responseBody string
	comparer     contract.Comparer
}

type extendedConfig struct {
	Script *scriptConfig `yaml:"script,omitempty"`
}

type scriptConfig struct {
	Path     string  `yaml:"path,omitempty"`
	Response *string `yaml:"response,omitempty"`
}

type Config struct {
	Cmd      string  `yaml:"script_path,omitempty"`
	Response *string `yaml:"script_response,omitempty"`
}

func (ex *extendedConfig) isEmpty() bool {
	return ex == nil || ex.Script == nil || ex.Script.Path == ""
}

func (c *Config) isEmpty() bool {
	return c == nil || c.Cmd == ""
}

func (e *ScriptCmd) SetVars(vv contract.Vars) {
	e.Vars = vv
}

func (e *ScriptCmd) SetReport(r contract.ReportAttachement) {
	e.report = r
}

func NewUnmarshaller(comparer contract.Comparer) *Unmarshaller {
	return &Unmarshaller{
		comparer: comparer,
	}
}

type Unmarshaller struct {
	comparer contract.Comparer
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
	if cfgExtended != nil && cfgExtended.Script != nil {
		return &ScriptCmd{
			comparer: u.comparer,
			Config: &Config{
				Cmd:      cfgExtended.Script.Path,
				Response: cfgExtended.Script.Response,
			},
		}, nil
	}
	return &ScriptCmd{
		comparer: u.comparer,
		Config:   cfg,
	}, nil
}

func (e *ScriptCmd) GetConfig() interface{} {
	return e.Config
}

func (e *ScriptCmd) Do() error {
	if e.Config != nil && e.Config.Cmd != "" {
		e.Config.Cmd = e.Vars.Apply(e.Config.Cmd)
		res, err := e.run(e.Config.Cmd)
		if err != nil {
			return err
		}
		e.responseBody = res

		return nil
	}

	return nil
}

func (e *ScriptCmd) ResponseBody() *string {
	return &e.responseBody
}

func (e *ScriptCmd) IsValid() error {
	return nil
}

func (e *ScriptCmd) Check() error {
	if e.Config.Response != nil {
		errs := e.comparer.Compare(*e.Config.Response, e.responseBody, contract.CompareParams{})
		if len(errs) > 0 {
			msg := ""
			for _, v := range errs {
				msg += v.Error() + "\n"
			}
			return &contract.TestError{
				Title:         "response body differs",
				Expected:      *e.Config.Response,
				Actual:        e.responseBody,
				Message:       msg,
				OriginalError: fmt.Errorf("response body differs: %v", msg),
			}
		}
	}
	return nil
}
