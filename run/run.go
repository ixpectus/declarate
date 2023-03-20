package run

import (
	"fmt"
	"os"

	"github.com/ixpectus/declarate/contract"
	"github.com/k0kubun/pp"
	"gopkg.in/yaml.v2"
)

var (
	builders    []contract.CommandBuilder
	currentVars contract.Vars
)

type Runner struct {
	config RunnerConfig
}

type RunnerConfig struct {
	file      string
	variables contract.Vars
	builders  []contract.CommandBuilder
}

func New(c RunnerConfig) *Runner {
	builders = c.builders
	return &Runner{
		config: c,
	}
}

func (c *Runner) Run() error {
	file, err := os.ReadFile(c.config.file)
	if err != nil {
		return fmt.Errorf("file open: %w", err)
	}
	currentVars = c.config.variables
	configs := []runConfig{}
	pp.ColoringEnabled = false
	yaml.Unmarshal(file, &configs)
	run(configs)
	return nil
}

func run(
	cc []runConfig,
) {
	for _, v := range cc {
		runOne(v, 0)
	}
}

func runOne(
	conf runConfig,
	lvl int,
) {
	prefix := ""
	for i := 0; i < lvl; i++ {
		prefix += " "
	}
	if conf.Name != "" {
		fmt.Printf(prefix+"run test with name %v\n", conf.Name)
	}
	for _, c := range conf.Commands {
		c.SetVars(currentVars)
		c.Do()
		if err := c.Check(); err != nil {
			fmt.Printf(prefix+"test failed %v\n", err)
		}
	}
	if len(conf.Steps) > 0 {
		for _, v := range conf.Steps {
			runOne(v, lvl+1)
		}
	}
}
