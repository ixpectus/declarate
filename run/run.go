package run

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/ixpectus/declarate/contract"
	"github.com/ixpectus/declarate/variables"
	"gopkg.in/yaml.v2"
)

var (
	builders    []contract.CommandBuilder
	currentVars contract.Vars
)

type Runner struct {
	config RunnerConfig
	output contract.Output
}

type RunnerConfig struct {
	Variables contract.Vars
	Builders  []contract.CommandBuilder
	Output    contract.Output
	Wrapper   contract.TestWrapper
}

func New(c RunnerConfig) *Runner {
	builders = c.Builders
	return &Runner{
		config: c,
		output: c.Output,
	}
}

func (r *Runner) Run(fileName string) error {
	file, err := os.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("file open: %w", err)
	}
	currentVars = r.config.Variables
	configs := []runConfig{}
	if err := yaml.Unmarshal(file, &configs); err != nil {
		return fmt.Errorf("unmarshall failed for file %s: %w", fileName, err)
	}
	return r.run(configs, fileName)
}

func (r *Runner) run(
	cc []runConfig,
	fileName string,
) error {
	for _, v := range cc {
		if len(v.Commands) == 0 && len(v.Steps) == 0 {
			continue
		}
		r.beforeTest(fileName, v, 0)
		var err error
		var testResult *Result
		if v.Name != "" {
			r.output.Log(contract.Message{
				Name:    v.Name,
				Message: fmt.Sprintf("start  %v:%v", fileName, v.Name),
				Type:    contract.MessageTypeNotify,
			})
		}
		if len(v.Poll.PollInterval()) > 0 {
			testResult, err = r.runWithPollInterval(v, fileName)
		} else {
			testResult, err = r.runOne(v, 0, fileName)
		}
		testResult.FileName = fileName
		testResult.Name = v.Name
		if err != nil {
			return err
		}
		r.afterTest(v, *testResult)
		if testResult.Err != nil {
			r.outputErr(*testResult)
		} else {
			r.output.Log(contract.Message{
				Name:    v.Name,
				Message: fmt.Sprintf("passed %v:%v", fileName, v.Name),
				Type:    contract.MessageTypeSuccess,
			})
		}
	}
	return nil
}

func (r *Runner) beforeTest(file string, conf runConfig, lvl int) {
	if r.config.Wrapper != nil {
		r.config.Wrapper.BeforeTest(file, contract.RunConfig{
			Name:           conf.Name,
			Vars:           currentVars,
			VariablesToSet: conf.VariablesToSet,
			Commands:       conf.Commands,
		}, lvl)
	}
}

func (r *Runner) afterTest(conf runConfig, result Result) {
	if r.config.Wrapper != nil {
		r.config.Wrapper.AfterTest(contract.RunConfig{
			Name:           conf.Name,
			Vars:           currentVars,
			VariablesToSet: conf.VariablesToSet,
			Commands:       conf.Commands,
		},
			contract.Result{
				Err:      result.Err,
				Name:     result.Name,
				Lvl:      result.Lvl,
				FileName: result.FileName,
				Response: result.Response,
			},
		)
	}
}

func (r *Runner) runWithPollInterval(v runConfig, fileName string) (*Result, error) {
	var err error
	var testResult *Result
	for _, d := range v.Poll.PollInterval() {
		testResult, err = r.runOne(v, 0, fileName)
		if err != nil {
			return nil, err
		}
		if testResult.Err != nil {
			if v.Poll.ResponseRegexp != "" {
				rx, err := regexp.Compile(v.Poll.ResponseRegexp)
				if err != nil {
					break
				}
				if testResult.Response == nil || !rx.MatchString(*testResult.Response) {
					break
				}
			}
			r.output.Log(contract.Message{
				Name:    v.Name,
				Message: fmt.Sprintf("Sleep %v before next poll request", d),
				Type:    contract.MessageTypeNotify,
			})
			time.Sleep(d)
		} else {
			break
		}
	}
	return testResult, err
}

func (r *Runner) runOne(
	conf runConfig,
	lvl int,
	fileName string,
) (*Result, error) {
	var body *string
	for _, c := range conf.Commands {
		c.SetVars(currentVars)
		if err := c.Do(); err != nil {
			return &Result{
				Err:      err,
				Name:     conf.Name,
				Lvl:      lvl,
				FileName: fileName,
			}, nil
		}

		body = c.ResponseBody()
		if err := c.Check(); err != nil {
			return &Result{
				Err:      err,
				Name:     conf.Name,
				Lvl:      lvl,
				FileName: fileName,
				Response: body,
			}, nil
		}

		if body != nil {
			if c.VariablesToSet() != nil {
				varsToSet := c.VariablesToSet()
				jsonVars := map[string]string{}
				for k, v := range varsToSet {
					if v == "*" {
						currentVars.Set(k, *body)
					} else {
						jsonVars[k] = v
					}
				}
				if len(jsonVars) > 0 {
					vars, err := variables.FromJSON(jsonVars, *body)
					if err != nil {
						return &Result{
							Err:      err,
							Name:     conf.Name,
							Lvl:      lvl,
							FileName: fileName,
						}, nil
					}
					for k, v := range vars {
						currentVars.Set(k, v)
					}
				}
			}
		}
	}
	if len(conf.Steps) > 0 {
		for _, v := range conf.Steps {
			if v.Name != "" {
				r.output.Log(contract.Message{
					Name:    v.Name,
					Message: fmt.Sprintf("start  %v:%v", fileName, v.Name),
					Lvl:     lvl + 1,
					Type:    contract.MessageTypeNotify,
				})
			}
			testResult, err := r.runOne(v, lvl+1, fileName)
			if testResult.Err != nil {
				return testResult, nil
			}
			if err != nil {
				return nil, err
			}
			r.output.Log(contract.Message{
				Name:    v.Name,
				Lvl:     lvl + 1,
				Message: fmt.Sprintf("passed %v:%v", fileName, v.Name),
				Type:    contract.MessageTypeSuccess,
			})
		}
	}

	return &Result{
		Response: body,
	}, nil
}

func (r *Runner) outputErr(res Result) {
	var errTest *contract.TestError
	if errors.As(res.Err, &errTest) {
		r.output.Log(contract.Message{
			Name:    res.Name,
			Message: res.Err.Error(),
			Title: fmt.Sprintf(
				"failed %v:%v\n %v",
				res.FileName,
				res.Name,
				errTest.Title,
			),
			Expected: errTest.Expected,
			Actual:   errTest.Actual,
			Lvl:      res.Lvl,
			Type:     contract.MessageTypeError,
		})
		return
	}
	if res.Err != nil {
		r.output.Log(contract.Message{
			Name:    res.Name,
			Message: fmt.Sprintf("failed %v", res.Err),
			Lvl:     res.Lvl,
			Type:    contract.MessageTypeError,
		})
	}
}
