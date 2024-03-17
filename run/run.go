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

// эта штука глобальная переменная, так как используется в run/config::UnmarshalYAML
var builders []contract.CommandBuilder // почему как глобальные переменные а не часть структуры
// потому что метод run вызывается из сьютов и там нужно сохранить переменные между вызовами различных частей сьюта
// но они ведь перетираются польностью

type Runner struct {
	config      RunnerConfig
	output      contract.Output
	currentVars contract.Vars
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

func (r *Runner) buildRunConfigs(fileName string) ([]runConfig, error) {
	file, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("file open: %w", err)
	}
	r.currentVars = r.config.Variables
	configs := []runConfig{}
	if err := yaml.Unmarshal(file, &configs); err != nil {
		return nil, fmt.Errorf("unmarshall failed for file %s: %w", fileName, err)
	}

	return configs, nil
}

func (r *Runner) Run(fileName string, t *testing.T) (bool, error) {
	configs, err := r.buildRunConfigs(fileName)
	if err != nil {
		return true, fmt.Errorf("unmarshall failed for file %s: %w", fileName, err)
	}
	for _, v := range configs {
		if len(v.Commands) == 0 && len(v.Steps) == 0 {
			// nothing to do
			continue
		}
		if v.Condition != "" && !condition.IsTrue(r.currentVars, v.Condition) {
			r.logSkip(v.Name, fileName, 0)
			continue
		}
		v.Name = r.currentVars.Apply(v.Name)
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
	// stores poll information, used for logs and reports
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
		if len(v.Poll.PollInterval())-1 == i { // last poll step
			isPolling = false
		}
		estimated := finish.Sub(time.Now())
		testResult, err = r.runOne(
			v,
			0,
			fileName,
			isPolling,
		)

		// unexpected test error run
		if err != nil {
			pollResult.Finish = time.Now()
			testResult.PollResult = &pollResult
			return nil, err
		}
		// test not passed
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
			break
			// if v.Poll.ResponseRegexp != "" || v.Poll.ResponseTmpls != nil {
			// 	break
			// }
		}
	}
	pollResult.Finish = time.Now()
	testResult.PollResult = &pollResult

	return testResult, err
}

func (r *Runner) setupCommand(cmd contract.Doer) contract.Doer {
	cmd.SetVars(r.currentVars)
	cmd.SetReport(r.config.Report)

	return cmd
}

func (r *Runner) fillVariablesByResponse(
	commandResponseBody *string,
	variablesToSet map[string]string,
) error {
	if commandResponseBody == nil || variablesToSet == nil {
		return nil
	}
	jsonVars := map[string]string{}
	for k, v := range variablesToSet {
		if v == "*" {
			r.currentVars.Set(k, *commandResponseBody)
		} else {
			jsonVars[k] = v
		}
	}
	if len(jsonVars) > 0 {
		vars, err := variables.FromJSON(jsonVars, *commandResponseBody, r.currentVars)
		if err != nil {
			return err
		}
		for k, v := range vars {
			r.currentVars.Set(k, v)
		}
	}

	return nil
}

func (r *Runner) fillPersistentVariablesByResponse(
	commandResponseBody *string,
	variablesToSet map[string]string,
) error {
	if commandResponseBody == nil || variablesToSet == nil {
		return nil
	}
	jsonVars := map[string]string{}
	for k, v := range variablesToSet {
		if v == "*" {
		} else {
			jsonVars[k] = v
		}
	}
	if len(jsonVars) > 0 {
		vars, err := variables.FromJSON(jsonVars, *commandResponseBody, r.currentVars)
		if err != nil {
			return err
		}
		for k, v := range vars {
			r.currentVars.SetPersistent(k, v)
		}
	}

	return nil
}

func (r *Runner) runCommand(cmd contract.Doer) (*string, error) {
	cmd = r.setupCommand(cmd)
	if err := cmd.Do(); err != nil {
		return nil, err
	}
	responseBody := cmd.ResponseBody()
	if err := cmd.Check(); err != nil {
		return responseBody, err
	}
	// fmt.Printf("\n>>> cmd variables %v <<< debug\n", cmd.Variables())
	// if err := r.fillVariablesByResponse(
	// 	responseBody,
	// 	cmd.Variables(),
	// ); err != nil {
	// 	return responseBody, err
	// }
	return responseBody, nil
}

func (r *Runner) runOne(
	conf runConfig,
	lvl int,
	fileName string,
	isPolling bool,
) (*Result, error) {
	var commandResponseBody *string
	var firstErrResult *Result
	conf.Name = r.currentVars.Apply(conf.Name)
	// fmt.Printf("\n>>> conf %v <<< debug\n", conf.Variables)
	// fmt.Printf("\n>>> %v <<< debug\n", len(conf.Commands))

	for _, command := range conf.Commands {
		r.beforeTestStep(fileName, &conf, lvl)
		var err error
		commandResponseBody, err = r.runCommand(command)
		if err != nil {
			res := &Result{
				Err:      err,
				Name:     conf.Name,
				Lvl:      lvl,
				FileName: fileName,
				Response: commandResponseBody,
			}
			r.afterTestStep(fileName, &conf, *res, isPolling)
			return res, nil
		}
	}
	// fmt.Printf("\n>>> %v <<< debug\n", commandResponseBody)

	if len(conf.Steps) > 0 {
		results := []string{}
		for _, stepRunConfig := range conf.Steps {
			if stepRunConfig.Condition != "" && !condition.IsTrue(r.config.Variables, stepRunConfig.Condition) {
				r.logSkip(stepRunConfig.Name, fileName, lvl+1)
				continue
			}
			if stepRunConfig.Name != "" && !isPolling {
				r.logStart(fileName, stepRunConfig, lvl+1)
			}
			var testResult *Result
			var err error
			action := func() {
				testResult, err = r.runOne(stepRunConfig, lvl+1, fileName, isPolling)
			}
			r.config.Report.Step(report.ReportOptions{Description: stepRunConfig.Name}, action)

			if testResult.Err != nil && isPolling {
				firstErrResult = testResult
				if testResult.Response != nil {
					results = append(results, *testResult.Response)
				} else {
					results = append(results, "")
				}
				continue
			}
			if testResult.Err != nil {
				r.afterTestStep(fileName, &conf, *testResult, isPolling)
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
			if !isPolling {
				r.logPass(stepRunConfig.Name, fileName, testResult, lvl+1)
			}
		}
		if len(results) > 0 {
			s := "[" + strings.Join(results, ", ") + "]"
			commandResponseBody = &s
		}
	}

	if err := r.fillVariablesByResponse(
		commandResponseBody,
		conf.Variables,
	); err != nil {
		res := &Result{
			Err:      err,
			Name:     conf.Name,
			Lvl:      lvl,
			FileName: fileName,
		}
		r.afterTestStep(fileName, &conf, *res, isPolling)
		return res, nil
	}
	if err := r.fillPersistentVariablesByResponse(commandResponseBody, conf.VariablesPersistent); err != nil {
		res := &Result{
			Err:      err,
			Name:     conf.Name,
			Lvl:      lvl,
			FileName: fileName,
		}
		r.afterTestStep(fileName, &conf, *res, isPolling)
		return res, nil
	}
	if firstErrResult != nil {
		firstErrResult.Response = commandResponseBody
		r.afterTestStep(fileName, &conf, *firstErrResult, isPolling)
		return firstErrResult, nil
	}

	res := &Result{
		Response: commandResponseBody,
		Lvl:      lvl,
		FileName: fileName,
	}

	r.afterTestStep(fileName, &conf, *res, isPolling)
	return res, nil
}
