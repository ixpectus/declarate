package run

import (
	"fmt"
	"os"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/ixpectus/declarate/compare"
	"github.com/ixpectus/declarate/condition"
	"github.com/ixpectus/declarate/contract"
	"github.com/ixpectus/declarate/report"
	"github.com/ixpectus/declarate/tools"
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
	Variables    contract.Vars
	Builders     []contract.CommandBuilder
	Output       contract.Output
	Wrapper      contract.TestWrapper
	T            *testing.T
	comparer     contract.Comparer
	pollComparer contract.Comparer
	Report       contract.Report
}

func New(c RunnerConfig) *Runner {
	builders = c.Builders
	if c.comparer == nil {
		c.comparer = compare.New(contract.CompareParams{
			IgnoreArraysOrdering: tools.To(true),
			DisallowExtraFields:  tools.To(false),
			AllowArrayExtraItems: tools.To(true),
		}, c.Variables)
	}
	if c.pollComparer == nil {
		c.pollComparer = compare.New(contract.CompareParams{
			IgnoreArraysOrdering: tools.To(true),
			DisallowExtraFields:  tools.To(false),
			AllowArrayExtraItems: tools.To(true),
		}, c.Variables)
	}
	if c.Report == nil {
		c.Report = report.NewEmptyReport()
	}
	c.Output.SetReport(c.Report)
	return &Runner{
		config: c,
		output: c.Output,
	}
}

func (r *Runner) Validate(fileName string) error {
	file, err := os.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("file open: %w", err)
	}
	currentVars = r.config.Variables
	configs := []runConfig{}
	if err := yaml.Unmarshal(file, &configs); err != nil {
		return fmt.Errorf("unmarshall failed for file %s: %w", r.filenameShort(fileName), err)
	}
	for _, v := range configs {
		for _, c := range v.Commands {
			if err := c.IsValid(); err != nil {
				return fmt.Errorf("invalid command, %w", err)
			}
		}
	}

	return nil
}

func (r *Runner) filenameShort(fileName string) string {
	parts := strings.Split(fileName, "/")
	if len(parts) > 4 {
		return path.Base(fileName)
	}
	return fileName
}

func (r *Runner) runFile(fileName string, t *testing.T) (bool, error) {
	file, err := os.ReadFile(fileName)
	if err != nil {
		return true, fmt.Errorf("file open: %w", err)
	}
	currentVars = r.config.Variables
	configs := []runConfig{}
	if err := yaml.Unmarshal(file, &configs); err != nil {
		return true, fmt.Errorf("unmarshall failed for file %s: %w", fileName, err)
	}

	for _, v := range configs {
		if len(v.Commands) == 0 && len(v.Steps) == 0 {
			continue
		}
		if v.Condition != "" && !condition.IsTrue(r.config.Variables, v.Condition) {
			r.output.Log(contract.Message{
				Filename: fileName,
				Name:     v.Name,
				Message:  fmt.Sprintf("skipped for file %s: %s", r.filenameShort(fileName), v.Name),
				Type:     contract.MessageTypeNotify,
			})
			continue
		}
		if t != nil {
			v.Name = currentVars.Apply(v.Name)
			var testResult *Result
			res := true
			var err error
			action := func() {
				testResult, err = r.run(v, fileName)
				if err != nil {
					r.output.Log(contract.Message{
						Filename:            fileName,
						Name:                v.Name,
						Message:             fmt.Sprintf("run failed for file %s: %s", r.filenameShort(fileName), err),
						Type:                contract.MessageTypeError,
						PollResult:          testResult.PollResult,
						PollConditionFailed: testResult.PollConditionFailed,
					})
					t.FailNow()
					res = false
				}
				if testResult.Err != nil {
					r.outputErr(*testResult)
					t.FailNow()
					res = false
				} else {
					r.output.Log(contract.Message{
						Filename:            fileName,
						Name:                v.Name,
						Message:             fmt.Sprintf("passed %v:%v", r.filenameShort(fileName), v.Name),
						Type:                contract.MessageTypeSuccess,
						PollResult:          testResult.PollResult,
						PollConditionFailed: testResult.PollConditionFailed,
					})
				}
			}
			r.config.Report.Step(report.ReportOptions{Description: v.Name}, action)
			if !res {
				return false, nil
			}
		} else {
			v.Name = currentVars.Apply(v.Name)
			testResult, err := r.run(v, fileName)
			if err != nil {
				r.output.Log(contract.Message{
					Filename:            fileName,
					Name:                v.Name,
					Message:             fmt.Sprintf("run failed for file %s: %s", r.filenameShort(fileName), err),
					Type:                contract.MessageTypeError,
					PollResult:          testResult.PollResult,
					PollConditionFailed: testResult.PollConditionFailed,
				})
				return false, nil
			}
			if testResult.Err != nil {
				r.outputErr(*testResult)
				return false, nil
			} else {
				r.output.Log(contract.Message{
					Filename:            fileName,
					Name:                v.Name,
					Message:             fmt.Sprintf("passed %v:%v", r.filenameShort(fileName), v.Name),
					Type:                contract.MessageTypeSuccess,
					PollResult:          testResult.PollResult,
					PollConditionFailed: testResult.PollConditionFailed,
				})
			}
		}
	}
	return false, nil
}

func (r *Runner) Run(fileName string, t *testing.T) (bool, error) {
	if t != nil {
		return r.runFile(fileName, t)
	}
	return r.runFile(fileName, t)
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
			Filename:       fileName,
			Name:           v.Name,
			HasNestedSteps: len(v.Steps) > 0,
			HasPoll:        len(v.Poll.PollInterval()) > 0,
			Message:        fmt.Sprintf("start %v:%v", fileName, v.Name),
			Type:           contract.MessageTypeNotify,
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
	v.Poll.comparer = r.config.pollComparer
	start := time.Now()
	finish := start
	for _, d := range v.Poll.PollInterval() {
		finish = finish.Add(d)
	}

	pollInfo := contract.PollInfo{
		Start:  start,
		Finish: finish,
	}
	pollResult := contract.PollResult{
		Start:         start,
		PlannedFinish: finish,
	}
	for i, d := range v.Poll.PollInterval() {
		isPolling := true
		if len(v.Poll.PollInterval())-1 == i {
			isPolling = false
		}
		estimated := finish.Sub(time.Now())
		testResult, err = r.runOne(v, 0, fileName, isPolling)
		if err != nil {
			pollResult.Finish = time.Now()
			testResult.PollResult = &pollResult
			return nil, err
		}
		if testResult.Err != nil {
			if v.Poll.ResponseRegexp != "" || v.Poll.ResponseTmpls != nil {
				res, errs, _ := v.Poll.pollContinue(testResult.Response)
				if !res {
					if len(errs) > 0 {
						testResult.PollConditionFailed = true
						testResult.Err = errs[0]
					}
					break
				}
			}
			r.output.Log(contract.Message{
				Filename: r.filenameShort(fileName),
				Name:     v.Name,
				Poll:     &pollInfo,
				Message: fmt.Sprintf(
					"poll %s:%s, wait %v, estimated %v",
					r.filenameShort(fileName),
					v.Name,
					d,
					estimated.Truncate(time.Second),
				),
				Type: contract.MessageTypeNotify,
			})
			time.Sleep(d)
		} else {
			if v.Poll.ResponseRegexp != "" || v.Poll.ResponseTmpls != nil {
				res, errs, _ := v.Poll.pollContinue(testResult.Response)
				if !res {
					if len(errs) > 0 {
						testResult.PollConditionFailed = true
						testResult.Err = errs[0]
					}
					break
				}
			}
		}
	}
	pollResult.Finish = time.Now()
	testResult.PollResult = &pollResult

	return testResult, err
}

func (r *Runner) runOne(
	conf runConfig,
	lvl int,
	fileName string,
	polling bool,
) (*Result, error) {
	var body *string
	var firstErrResult *Result
	for _, c := range conf.Commands {
		c.SetVars(currentVars)
		c.SetReport(r.config.Report)
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
			if v.Condition != "" && !condition.IsTrue(r.config.Variables, v.Condition) {
				r.output.Log(contract.Message{
					Filename: fileName,
					Name:     v.Name,
					Message:  fmt.Sprintf("skipped %s: %s", r.filenameShort(fileName), v.Name),
					Lvl:      lvl + 1,
					Type:     contract.MessageTypeNotify,
				})
				continue
			}
			if v.Name != "" {
				r.output.Log(contract.Message{
					Filename:       fileName,
					Name:           v.Name,
					Message:        fmt.Sprintf("start %v:%v", fileName, v.Name),
					HasNestedSteps: len(v.Steps) > 0,
					Lvl:            lvl + 1,
					Type:           contract.MessageTypeNotify,
				})
			}
			testResult, err := r.runOne(v, lvl+1, fileName, polling)
			if testResult.Err != nil && polling {
				firstErrResult = testResult
				if testResult.Response != nil {
					results = append(results, *testResult.Response)
				} else {
					results = append(results, "")
				}
				continue
			}
			if testResult.Err != nil {
				r.afterTestStep(fileName, &conf, *testResult, polling)
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
				Filename:   fileName,
				Name:       v.Name,
				Lvl:        lvl + 1,
				Message:    fmt.Sprintf("passed %v:%v", r.filenameShort(fileName), v.Name),
				Type:       contract.MessageTypeSuccess,
				PollResult: testResult.PollResult,
			})
		}
		if len(results) > 0 {
			s := "[" + strings.Join(results, ", ") + "]"
			body = &s
		}
	}

	if conf.Variables != nil {
		varsToSet := conf.Variables
		jsonVars := map[string]string{}
		for k, v := range varsToSet {
			if v == "*" {
				currentVars.Set(k, *body)
			} else {
				jsonVars[k] = v
			}
		}
		if len(jsonVars) > 0 && body != nil {
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
	if conf.VariablesPersistent != nil {
		varsToSet := conf.VariablesPersistent
		jsonVars := map[string]string{}
		for k, v := range varsToSet {
			if v == "*" {
				currentVars.SetPersistent(k, *body)
			} else {
				jsonVars[k] = v
			}
		}
		if len(jsonVars) > 0 && body != nil {
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
				currentVars.SetPersistent(k, v)
			}
		}
	}
	if firstErrResult != nil {
		firstErrResult.Response = body
		r.afterTestStep(fileName, &conf, *firstErrResult, polling)
		return firstErrResult, nil
	}

	res := &Result{
		Response: body,
		Lvl:      lvl,
		FileName: fileName,
	}
	r.afterTestStep(fileName, &conf, *res, polling)
	return res, nil
}
