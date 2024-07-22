package shell

import (
	"fmt"
	"strings"

	"github.com/dailymotion/allure-go"
	"github.com/ixpectus/declarate/contract"
)

type ShellCmd struct {
	Config       *Config
	Vars         contract.Vars
	report       contract.ReportAttachement
	responseBody string
	comparer     contract.Comparer
}

type extendedConfig struct {
	Shell *shellConfig `yaml:"shell,omitempty"`
}

type shellConfig struct {
	Cmd              string                 `yaml:"cmd,omitempty"`
	Response         *string                `yaml:"response,omitempty"`
	ComparisonParams contract.CompareParams `json:"comparisonParams" yaml:"comparisonParams"`
}

type Config struct {
	Cmd              string                 `yaml:"shell_cmd,omitempty"`
	Response         *string                `yaml:"shell_response,omitempty"`
	ComparisonParams contract.CompareParams `json:"comparisonParams" yaml:"comparisonParams"`
}

func (e *ShellCmd) SetVars(vv contract.Vars) {
	e.Vars = vv
}

func (e *ShellCmd) SetReport(r contract.ReportAttachement) {
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

func (ex *extendedConfig) isEmpty() bool {
	return ex == nil || ex.Shell == nil || ex.Shell.Cmd == ""
}

func (c *Config) isEmpty() bool {
	return c == nil || c.Cmd == ""
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
	if cfgExtended != nil && cfgExtended.Shell != nil {
		return &ShellCmd{
			comparer: u.comparer,
			Config: &Config{
				Cmd:              cfgExtended.Shell.Cmd,
				Response:         cfgExtended.Shell.Response,
				ComparisonParams: cfgExtended.Shell.ComparisonParams,
			},
		}, nil
	}
	return &ShellCmd{
		comparer: u.comparer,
		Config:   cfg,
	}, nil
}

func (e *ShellCmd) Do() error {
	if e.Config != nil && e.Config.Cmd != "" {
		e.Config.Cmd = e.Vars.Apply(e.Config.Cmd)
		res, err := e.run(e.Config.Cmd)
		if err != nil {
			return err
		}
		e.responseBody = strings.Join(res, "\n")
		if e.report != nil {
			e.report.AddAttachment("response", allure.TextPlain, []byte(e.responseBody))
		}

		return nil
	}

	return nil
}

func (e *ShellCmd) GetConfig() interface{} {
	return e.Config
}

func (e *ShellCmd) IsValid() error {
	return nil
}

func (e *ShellCmd) ResponseBody() *string {
	return &e.responseBody
}

// func (e *ShellCmd) Check() error {
// 	if e.Config.Response != nil {
// 		linesExpected := strings.Split(*e.Config.Response, "\n")
// 		linesGot := strings.Split(e.responseBody, "\n")

// 		if len(linesExpected) != len(linesGot) {
// 			errMsg := fmt.Sprintf("lines count differs, expected %v, got %v", len(linesExpected), len(linesGot))
// 			for k := range linesExpected {
// 				if len(linesGot) > k {
// 					if linesExpected[k] != linesGot[k] {
// 						errMsg += fmt.Sprintf("\nlines different at line %v, expected %v, got %v", k, linesExpected[k], linesExpected[k])
// 					}
// 				}
// 			}
// 			res := compare.MakeError(
// 				"",
// 				errMsg,
// 				*e.Config.Response,
// 				e.responseBody,
// 			)
// 			return res
// 		}
// 		errMsg := ""
// 		for k := range linesExpected {
// 			if len(linesGot) >= k {
// 				if strings.Trim(linesExpected[k], " ") != strings.Trim(linesGot[k], " ") {
// 					errMsg = fmt.Sprintf("\nlines different at line %v, expected `%v`, got `%v`", k, linesExpected[k], linesGot[k])
// 					break
// 				}
// 			}
// 		}
// 		if errMsg != "" {
// 			return &contract.TestError{
// 				Title:         "response body differs",
// 				Expected:      *e.Config.Response,
// 				Actual:        e.responseBody,
// 				Message:       errMsg,
// 				OriginalError: nil,
// 			}
// 		}
// 	}
// 	return nil
// }

func (e *ShellCmd) Check() error {
	if e.Config.Response != nil {
		var (
			err error
			errs []error
		)
		if e.Config.ComparisonParams.CompareJson != nil && *e.Config.ComparisonParams.CompareJson {
			errs, err = e.comparer.CompareJsonBody(*e.Config.Response, e.responseBody, e.Config.ComparisonParams)
			if err != nil {
				return fmt.Errorf("compare json failed: %w", err)
			}
		} else {
			errs = e.comparer.Compare(*e.Config.Response, e.responseBody, e.Config.ComparisonParams)
		}
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
