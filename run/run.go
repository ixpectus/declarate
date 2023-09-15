package run

import (
	"fmt"
	"os"
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

func (r *Runner) Run(fileName string, t *testing.T) (bool, error) {
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
		if v.Condition != "" && !condition.IsTrue(currentVars, v.Condition) {
			r.logSkip(v.Name, fileName, 0)
			continue
		}
		v.Name = currentVars.Apply(v.Name)
		var testResult *Result
		res := true
		var err error
		action := func() {
			testResult, err = r.run(v, fileName)
			if err != nil {
				r.logRunFail(v.Name, fileName, err, testResult)
				if t != nil {
					t.FailNow()
				}
				res = false
			}
			if testResult.Err != nil {
				r.logErr(*testResult)
				if t != nil {
					t.FailNow()
				}
				res = false
			} else {
				r.logPass(v.Name, fileName, testResult, 0)
			}
		}
		r.config.Report.Step(report.ReportOptions{Description: v.Name}, action)
		if !res {
			return false, nil
		}

	}
	return false, nil
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
		r.logStart(fileName, v, 0)
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
			r.logPoll(fileName, v, pollInfo, d, estimated)
			time.Sleep(d)
		} else {
			if v.Poll.ResponseRegexp != "" || v.Poll.ResponseTmpls != nil {
				break
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
					vars, err := variables.FromJSON(jsonVars, *body, currentVars)
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
				r.logSkip(v.Name, fileName, lvl+1)
				continue
			}
			if v.Name != "" && !polling {
				r.logStart(fileName, v, lvl+1)
			}
			var testResult *Result
			var err error
			action := func() {
				testResult, err = r.runOne(v, lvl+1, fileName, polling)
			}
			r.config.Report.Step(report.ReportOptions{Description: v.Name}, action)

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
			if !polling {
				r.logPass(v.Name, fileName, testResult, lvl+1)
			}
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
			vars, err := variables.FromJSON(jsonVars, *body, currentVars)
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
			vars, err := variables.FromJSON(jsonVars, *body, currentVars)
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
