package script

import (
	"fmt"
	"strings"

	"github.com/ixpectus/declarate/compare"
	"github.com/ixpectus/declarate/contract"
)

type ScriptCmd struct {
	Config       *Config
	Vars         contract.Vars
	responseBody string
	comparer     contract.Comparer
}

type extendedConfig struct {
	Script *scriptConfig `yaml:"script,omitempty"`
}

type scriptConfig struct {
	Path      string            `yaml:"path,omitempty"`
	Variables map[string]string `yaml:"variables,omitempty"`
	Response  *string           `yaml:"response,omitempty"`
}

type Config struct {
	Cmd       string            `yaml:"script_path,omitempty"`
	Variables map[string]string `yaml:"variables,omitempty"`
	Response  *string           `yaml:"script_response,omitempty"`
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
				Cmd:       cfgExtended.Script.Path,
				Variables: cfgExtended.Script.Variables,
				Response:  cfgExtended.Script.Response,
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
		res, err := Run(e.Config.Cmd)
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
		linesExpected := strings.Split(*e.Config.Response, "\n")
		linesGot := strings.Split(e.responseBody, "\n")
		if len(linesExpected) != len(linesGot) {
			errMsg := fmt.Sprintf("lines count differs, expected %v, got %v", len(linesExpected), len(linesGot))
			for k := range linesExpected {
				if len(linesGot) >= k {
					if linesExpected[k] != linesGot[k] {
						errMsg += fmt.Sprintf("\nlines different at line %v, expected %v, got %v", k, linesExpected[k], linesExpected[k])
					}
				}
			}
			res := compare.MakeError(
				"",
				errMsg,
				e.responseBody,
				*e.Config.Response,
			)
			return res
		}
		errMsg := ""
		for k := range linesExpected {
			if len(linesGot) >= k {
				if strings.Trim(linesExpected[k], " ") != strings.Trim(linesGot[k], " ") {
					errMsg = fmt.Sprintf("\nlines diffsserent at line %v, expected %v, got %v", k, linesExpected[k], linesExpected[k])
				}
			}
		}
		if errMsg != "" {
			return fmt.Errorf(errMsg)
		}
		// errors := e.comparer.Compare(
		// 	*e.Config.Response,
		// 	e.responseBody,
		// 	compare.CompareParams{},
		// )
		// if len(errors) > 0 {
		// 	return errors[0]
		// }
	}
	return nil
}

func (e *ScriptCmd) VariablesToSet() map[string]string {
	if e.Config != nil {
		return e.Config.Variables
	}
	return nil
}
