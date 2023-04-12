package run

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ixpectus/declarate/compare"
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
	T         *testing.T
	comparer  contract.Comparer
}

func New(c RunnerConfig) *Runner {
	builders = c.Builders
	if c.comparer == nil {
		c.comparer = compare.New(compare.CompareParams{})
	}
	return &Runner{
		config: c,
		output: c.Output,
	}
}

func (r *Runner) Run(fileName string, t *testing.T) error {
	file, err := os.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("file open: %w", err)
	}
	currentVars = r.config.Variables
	configs := []runConfig{}
	if err := yaml.Unmarshal(file, &configs); err != nil {
		return fmt.Errorf("unmarshall failed for file %s: %w", fileName, err)
	}
	for _, v := range configs {
		if len(v.Commands) == 0 && len(v.Steps) == 0 {
			continue
		}
		if t != nil {
			v.Name = currentVars.Apply(v.Name)
			failed := t.Run(v.Name, func(t1 *testing.T) {
				testResult, err := r.run(v, fileName)
				if err != nil {
					r.output.Log(contract.Message{
						Name:    v.Name,
						Message: fmt.Sprintf("run failed for file %s: %s", fileName, err),
						Type:    contract.MessageTypeError,
					})
					t1.FailNow()
				}
				if testResult.Err != nil {
					r.outputErr(*testResult)
					t1.FailNow()
				} else {
					r.output.Log(contract.Message{
						Name:    v.Name,
						Message: fmt.Sprintf("passed %v:%v", fileName, v.Name),
						Type:    contract.MessageTypeSuccess,
					})
				}
			})
			if !failed {
				t.FailNow()
			}
		} else {
			v.Name = currentVars.Apply(v.Name)
			testResult, err := r.run(v, fileName)
			if err != nil {
				r.output.Log(contract.Message{
					Name:    v.Name,
					Message: fmt.Sprintf("run failed for file %s: %s", fileName, err),
					Type:    contract.MessageTypeError,
				})
			}
			if testResult.Err != nil {
				r.outputErr(*testResult)
				continue
			} else {
				r.output.Log(contract.Message{
					Name:    v.Name,
					Message: fmt.Sprintf("passed %v:%v", fileName, v.Name),
					Type:    contract.MessageTypeSuccess,
				})
			}
		}
	}
	return nil
}

func (r *Runner) run(
	v runConfig,
	fileName string,
) (*Result, error) {
	r.beforeTest(fileName, &v, 0)
	var (
		err        error
		testResult *Result
	)
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
		testResult, err = r.runOne(v, 0, fileName, false)
	}
	if err != nil {
		return testResult, fmt.Errorf("run test for file %s: %w", fileName, err)
	}
	r.afterTest(fileName, v, *testResult)
	return testResult, nil
}

func (r *Runner) runWithPollInterval(v runConfig, fileName string) (*Result, error) {
	var err error
	var testResult *Result
	v.Poll.comparer = r.config.comparer
	for i, d := range v.Poll.PollInterval() {
		isPolling := true
		if len(v.Poll.PollInterval())-1 == i {
			isPolling = false
		}
		testResult, err = r.runOne(v, 0, fileName, isPolling)
		if err != nil {
			return nil, err
		}
		if testResult.Err != nil {
			if v.Poll.ResponseRegexp != "" || v.Poll.ResponseTmpls != nil {
				if !v.Poll.pollContinue(testResult.Response) {
					break
				}
				// rx, err := regexp.Compile(v.Poll.ResponseRegexp)
				// if err != nil {
				// 	break
				// }
				// if testResult.Response == nil || !rx.MatchString(*testResult.Response) {
				// 	break
				// }
			}
			r.output.Log(contract.Message{
				Name:    v.Name,
				Message: fmt.Sprintf("poll %s:%s, wait %v", fileName, v.Name, d),
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
	polling bool,
) (*Result, error) {
	var body *string
	for _, c := range conf.Commands {
		c.SetVars(currentVars)
		conf.Name = currentVars.Apply(conf.Name)
		r.beforeTestStep(fileName, &conf, lvl)
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
			res := &Result{
				Err:      err,
				Name:     conf.Name,
				Lvl:      lvl,
				FileName: fileName,
				Response: body,
			}
			r.afterTestStep(fileName, &conf, *res, polling)
			return res, nil
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
						res := &Result{
							Err:      err,
							Name:     conf.Name,
							Lvl:      lvl,
							FileName: fileName,
						}
						r.afterTestStep(fileName, &conf, *res, polling)
						return res, nil
					}
					for k, v := range vars {
						currentVars.Set(k, v)
					}
				}
			}
		}
	}
	if len(conf.Steps) > 0 {
		results := []string{}
		for _, v := range conf.Steps {
			if v.Name != "" {
				r.output.Log(contract.Message{
					Name:    v.Name,
					Message: fmt.Sprintf("start  %v:%v", fileName, v.Name),
					Lvl:     lvl + 1,
					Type:    contract.MessageTypeNotify,
				})
			}
			testResult, err := r.runOne(v, lvl+1, fileName, polling)
			if testResult.Err != nil {
				return testResult, nil
			}
			if testResult.Response != nil {
				results = append(results, *testResult.Response)
			} else {
				results = append(results, "")
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
		if len(results) > 0 {
			s := "[" + strings.Join(results, ", ") + "]"
			body = &s
		}
	}

	if conf.VariablesToSet != nil {
		varsToSet := conf.VariablesToSet
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
				res := &Result{
					Err:      err,
					Name:     conf.Name,
					Lvl:      lvl,
					FileName: fileName,
				}
				r.afterTestStep(fileName, &conf, *res, polling)
				return res, nil
			}
			for k, v := range vars {
				currentVars.Set(k, v)
			}
		}
	}

	res := &Result{
		Response: body,
		Lvl:      lvl,
		FileName: fileName,
	}
	r.afterTestStep(fileName, &conf, *res, polling)
	return res, nil
}
