package run

import (
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
