package suite

import (
	"fmt"
	"os"
	"strings"

	"github.com/ixpectus/declarate/contract"
	"github.com/ixpectus/declarate/run"
)

type SuiteConfig struct {
	Builders []contract.CommandBuilder
	Output   contract.Output
}

type RunConfig struct {
	RunAll       bool
	Tags         []string
	Filename     []string
	SkipFilename []string
	DryRun       bool
	Variables    contract.Vars
	Builders     []contract.CommandBuilder
	Output       contract.Output
}

type Suite struct {
	Directory string
	Config    RunConfig
}

func New(directory string, cfg RunConfig) *Suite {
	return &Suite{
		Directory: directory,
		Config:    cfg,
	}
}

func (s *Suite) Run() error {
	tests, err := s.allTests(s.Directory)
	if err != nil {
		return fmt.Errorf("run: %w", err)
	}
	if s.Config.DryRun {
		fmt.Println(fmt.Sprintf("tests to run\n%s", strings.Join(tests, "\n")))
		return nil
	}

	runner := run.New(run.RunnerConfig{
		Variables: s.Config.Variables,
		Output:    s.Config.Output,
		Builders:  s.Config.Builders,
	},
	)
	for _, v := range tests {
		_ = runner.Run(v)
	}
	return nil
}

func (r *Suite) allTests(testPath string) ([]string, error) {
	stat, err := os.Stat(testPath)
	if err != nil {
		return nil, err
	}
	if !stat.IsDir() {
		return []string{testPath}, nil
	}
	res := []string{}
	files, err := os.ReadDir(testPath)
	if err != nil {
		return nil, fmt.Errorf("load all tests: %w", err)
	}
	for _, v := range files {
		foundSkipped := false
		for _, vv := range r.Config.SkipFilename {
			if vv == v.Name() || vv+".yaml" == v.Name() {
				foundSkipped = true
				break
			}
		}
		if foundSkipped {
			continue
		}
		if v.IsDir() {
			nested, err := r.allTests(testPath + "/" + v.Name())
			if err != nil {
				return nil, err
			}
			res = append(res, nested...)
		} else {
			res = append(res, fmt.Sprintf("%s/%s", testPath, v.Name()))
		}
	}

	return res, nil
}
