package run

import (
	"errors"
	"fmt"

	"github.com/ixpectus/declarate/contract"
)

func (r *Runner) beforeTestStep(file string, conf *runConfig, lvl int) {
	if r.config.Wrapper != nil {
		cfg := &contract.RunConfig{
			Name:      conf.Name,
			Vars:      currentVars,
			Variables: conf.Variables,
			Commands:  conf.Commands,
		}
		r.config.Wrapper.BeforeTestStep(file, cfg, lvl)
		conf.Commands = cfg.Commands
	}
}

func (r *Runner) afterTestStep(file string, conf *runConfig, result Result, polling bool) {
	if r.config.Wrapper != nil {
		cfg := &contract.RunConfig{
			Name:      conf.Name,
			Vars:      currentVars,
			Variables: conf.Variables,
			Commands:  conf.Commands,
		}
		r.config.Wrapper.AfterTestStep(cfg,
			contract.Result{
				Err:      result.Err,
				Name:     conf.Name,
				Lvl:      result.Lvl,
				FileName: file,
				Response: result.Response,
			},
			polling,
		)
		conf.Commands = cfg.Commands
	}
}

func (r *Runner) beforeTest(file string, conf *runConfig, lvl int) {
	if r.config.Wrapper != nil {
		cfg := &contract.RunConfig{
			Name:      conf.Name,
			Vars:      currentVars,
			Variables: conf.Variables,
			Commands:  conf.Commands,
		}
		r.config.Wrapper.BeforeTest(file, cfg, lvl)
		conf.Commands = cfg.Commands
	}
}

func (r *Runner) afterTest(file string, conf runConfig, result Result) {
	if r.config.Wrapper != nil {
		cfg := &contract.RunConfig{
			Name:      conf.Name,
			Vars:      currentVars,
			Variables: conf.Variables,
			Commands:  conf.Commands,
		}
		r.config.Wrapper.AfterTest(cfg,
			contract.Result{
				Err:        result.Err,
				Name:       conf.Name,
				Lvl:        result.Lvl,
				FileName:   file,
				Response:   result.Response,
				PollResult: result.PollResult,
			},
		)
		conf.Commands = cfg.Commands
	}
}

func (r *Runner) outputErr(res Result) {
	var errTest *contract.TestError
	if errors.As(res.Err, &errTest) {
		r.output.Log(contract.Message{
			Filename: res.FileName,
			Name:     res.Name,
			Message:  res.Err.Error(),
			Title: fmt.Sprintf(
				"failed %v:%v\n%v",
				res.FileName,
				res.Name,
				errTest.Title,
			),
			Expected:            errTest.Expected,
			Actual:              errTest.Actual,
			Lvl:                 res.Lvl,
			Type:                contract.MessageTypeError,
			PollResult:          res.PollResult,
			PollConditionFailed: res.PollConditionFailed,
		})
		return
	}
	if res.Err != nil {
		r.output.Log(contract.Message{
			Filename:            res.FileName,
			Name:                res.Name,
			Message:             fmt.Sprintf("failed %v", res.Err),
			Lvl:                 res.Lvl,
			Type:                contract.MessageTypeError,
			PollConditionFailed: res.PollConditionFailed,
		})
	}
}
