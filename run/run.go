package run

import (
	"errors"
	"fmt"
	"os"

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
	file      string
	variables contract.Vars
	builders  []contract.CommandBuilder
	output    contract.Output
}

func New(c RunnerConfig) *Runner {
	builders = c.builders
	return &Runner{
		config: c,
		output: c.output,
	}
}

func (r *Runner) Run() error {
	file, err := os.ReadFile(r.config.file)
	if err != nil {
		return fmt.Errorf("file open: %w", err)
	}
	currentVars = r.config.variables
	configs := []runConfig{}
	yaml.Unmarshal(file, &configs)
	r.run(configs)
	return nil
}

func (r *Runner) run(
	cc []runConfig,
) {
	for _, v := range cc {
		err := r.runOne(v, 0)
		if err != nil {
			fmt.Printf("\n>>> %v <<< debug\n", err)
		}
	}
}

func (r *Runner) runOne(
	conf runConfig,
	lvl int,
) error {
	prefix := ""
	for i := 0; i < lvl; i++ {
		prefix += " "
	}
	if conf.Name != "" {
		r.output.Log(contract.Message{
			Name:    conf.Name,
			Message: fmt.Sprintf("start  %v:%v", r.config.file, conf.Name),
			Lvl:     lvl,
			Type:    contract.MessageTypeNotify,
		})
	}
	for _, c := range conf.Commands {
		c.SetVars(currentVars)
		if err := c.Do(); err != nil {
			r.outputErr(err, conf, lvl)
			return err
		}
		if err := c.Check(); err != nil {
			r.outputErr(err, conf, lvl)
			return err
		}

		body := c.ResponseBody()
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
					r.outputErr(err, conf, lvl)
					if err != nil {
						return fmt.Errorf(prefix+"test failed %v", err)
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
			err := r.runOne(v, lvl+1)
			if err != nil {
				return err
			}
		}
	}

	r.output.Log(contract.Message{
		Name:    conf.Name,
		Message: fmt.Sprintf("passed %v:%v", r.config.file, conf.Name),
		Lvl:     lvl,
		Type:    contract.MessageTypeSuccess,
	})
	return nil
}

func (r *Runner) outputErr(err error, conf runConfig, lvl int) {
	var errTest *contract.TestError
	if errors.As(err, &errTest) {
		r.output.Log(contract.Message{
			Name:    conf.Name,
			Message: err.Error(),
			Title: fmt.Sprintf(
				"failed %v:%v\n %v",
				r.config.file,
				conf.Name,
				errTest.Title,
			),
			Expected: errTest.Expected,
			Actual:   errTest.Actual,
			Lvl:      lvl,
			Type:     contract.MessageTypeError,
		})
		return
	}
	if err != nil {
		r.output.Log(contract.Message{
			Name:    conf.Name,
			Message: fmt.Sprintf("failed %v", err),
			Lvl:     lvl,
			Type:    contract.MessageTypeError,
		})
	}
}
