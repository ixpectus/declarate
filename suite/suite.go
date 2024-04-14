package suite

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/fatih/color"
	"github.com/ixpectus/declarate/condition"
	"github.com/ixpectus/declarate/contract"
	"github.com/ixpectus/declarate/report"
	"github.com/ixpectus/declarate/run"
	"github.com/ixpectus/declarate/tools"
	"github.com/recoilme/pudge"
	"gopkg.in/yaml.v2"
)

type SuiteConfig struct {
	Builders []contract.CommandBuilder
	Output   contract.Output
}

type RunConfig struct {
	NoColor           bool
	RunAll            bool
	FailFast          bool
	Tags              []string
	Filepathes        []string
	SkipFilename      []string
	DryRun            bool
	Continue          bool
	CleanRun          bool
	Report            contract.Report
	SuiteName         string
	SubSuiteName      string
	Variables         contract.Vars
	Builders          []contract.CommandBuilder
	Output            contract.Output
	TestRunWrapper    contract.TestWrapper
	T                 *testing.T
	PersistentStorage contract.Persistent
}

type Suite struct {
	Directory string
	Config    RunConfig
}

func New(directory string, cfg RunConfig) *Suite {
	if cfg.Report == nil {
		cfg.Report = report.NewEmptyReport()
		if cfg.Output != nil {
			cfg.Output.SetReport(cfg.Report)
		}
	}

	return &Suite{
		Directory: directory,
		Config:    cfg,
	}
}

func (s *Suite) testsDefinitions(tests []string) ([]testWithDefinition, error) {
	definitions := make([]testWithDefinition, 0, len(tests))
	for _, v := range tests {
		data, err := os.ReadFile(v)
		if err != nil {
			return nil, err
		}
		var testDefinitions []testDefinition
		err = yaml.Unmarshal(data, &testDefinitions)
		if err != nil {
			return nil, fmt.Errorf("parse test definitions from file %s: %w", v, err)
		}
		if len(testDefinitions) == 0 {
			continue
		}
		realDefinitions := []testWithDefinition{}
		for _, d := range testDefinitions {
			if d.Definition != nil {
				realDefinitions = append(realDefinitions, testWithDefinition{
					file:       v,
					definition: d,
				})
			}
		}
		if len(realDefinitions) > 1 {
			return nil, fmt.Errorf("test file %v should have only one definition", v)
		}
		if len(realDefinitions) > 0 {
			definitions = append(definitions, realDefinitions[0])
		}
	}
	return definitions, nil
}

const (
	keyToRun  = "to_run"
	keyRunned = "runned"
)

func (s *Suite) runnedTests() ([]string, error) {
	rawTests, err := s.Config.PersistentStorage.Get(keyRunned)
	if err != nil && !errors.Is(err, pudge.ErrKeyNotFound) {
		return nil, err
	}
	tests := []string{}
	if rawTests != "" {
		tests = strings.Split(rawTests, ",")
	}

	return tests, nil
}

func (s *Suite) addRunnedTest(testName string) error {
	if s.Config.PersistentStorage == nil {
		return nil
	}
	tests, err := s.runnedTests()
	if err != nil {
		return err
	}
	tests = append(tests, testName)
	s.Config.PersistentStorage.Set(keyRunned, strings.Join(tests, ","))

	return nil
}

func (s *Suite) Run() error {
	if s.Config.CleanRun {
		s.Config.PersistentStorage.Reset()
	}
	color.NoColor = s.Config.NoColor

	allTests, err := s.AllTests(s.Directory)
	if err != nil {
		return fmt.Errorf("run: %w", err)
	}
	tests, err := s.filterTestsByTags(allTests)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("filter tests by tags: %w", err)
	}
	if !s.Config.Continue && s.Config.PersistentStorage != nil {
		s.Config.PersistentStorage.Set(keyToRun, "")
		s.Config.PersistentStorage.Set(keyRunned, "")
	}
	if s.Config.PersistentStorage != nil {
		s.Config.PersistentStorage.Set(keyToRun, strings.Join(tests, ","))
	}

	if len(s.Config.Filepathes) > 0 {
		if len(s.Config.Tags) > 0 {
			tests = s.filterTestsByPathes(allTests, tests)
		} else {
			tests = s.filterTestsByPathes(allTests, []string{})
		}
	}
	if s.Config.Continue {
		runned, _ := s.runnedTests()
		tests = s.filterTestsByAlreadyRun(tests, runned)
	}

	runner := run.New(run.RunnerConfig{
		Variables: s.Config.Variables,
		Output:    s.Config.Output,
		Builders:  s.Config.Builders,
		Report:    s.Config.Report,
		Wrapper:   s.Config.TestRunWrapper,
		T:         s.Config.T,
	},
	)

	if s.Config.DryRun {
		fmt.Println(fmt.Sprintf("tests to run\n%s", strings.Join(tests, "\n")))
		if err := s.validate(tests, runner); err != nil {
			return err
		}
		return nil
	}
	if err := s.validate(tests, runner); err != nil {
		return err
	}
	failed := false
	for _, v := range tests {
		definitions, err := s.testsDefinitions([]string{v})
		if err != nil {
			log.Println(err)
			s.Config.T.Fail()
		}
		var description string
		if len(definitions) > 0 {
			if definitions[0].definition.Definition.Condition != "" {
				if !condition.IsTrue(
					s.Config.Variables,
					definitions[0].definition.Definition.Condition,
				) {
					log.Println(fmt.Sprintf("test %s skipped by condition", v))
					continue
				}
			}
			description = definitions[0].definition.Definition.Description
		}

		if s.Config.T != nil {
			s.Config.T.Run(tools.FilenameLastN(v, 2), func(t *testing.T) {
				if s.Config.T.Failed() && !failed {
					failed = true
				}
				if failed && s.Config.FailFast {
					t.Skip()
				}
				action := func() {
					failed, err = runner.Run(v, t)
					if err != nil {
						log.Println(err)
						t.Fail()
					}
				}
				s.Config.Report.Test(
					t,
					action,
					report.ReportOptions{
						Description: description,
						Suite:       s.Config.SuiteName,
						Epic:        s.Config.SuiteName + "-epic",
						SubSuite:    s.Config.SubSuiteName,
						Tags:        s.Config.Tags,
					},
				)
				if !s.Config.T.Failed() && !failed {
					s.addRunnedTest(v)
				}
			})
		} else {
			failed, err = runner.Run(v, nil)
			if err != nil {
				log.Println(err)
				if s.Config.FailFast {
					return err
				}
			}
			if !failed {
				s.addRunnedTest(v)
			}
		}
	}
	return nil
}

func (s *Suite) validate(tests []string, runner *run.Runner) error {
	hasInvalid := false
	for _, v := range tests {
		err := runner.Validate(v)
		if err != nil {
			log.Println(fmt.Sprintf("invalid test `%s` description\n  %v", v, err))
			hasInvalid = true
		}
	}
	if hasInvalid {
		if s.Config.T != nil {
			s.Config.T.FailNow()
		} else {
			return fmt.Errorf("tests validation failed")
		}
	}
	return nil
}

func (s *Suite) filterTestsByPathes(
	allTests []string,
	selectedTests []string,
) []string {
	for _, v := range s.Config.Filepathes {
		fullName := s.Directory + "/" + v
		for _, v1 := range allTests {
			if strings.Contains(v1, fullName) && !tools.Contains(selectedTests, v1) {
				selectedTests = append(selectedTests, v1)
			}
		}
	}

	return selectedTests
}

func (s *Suite) filterTestsByAlreadyRun(
	allTests []string,
	alreadyRun []string,
) []string {
	res := []string{}
	for _, v := range allTests {
		found := false
		for _, v1 := range alreadyRun {
			if v1 == v {
				found = true
				break
			}
		}
		if !found {
			res = append(res, v)
		}
	}

	return res
}

func (r *Suite) filterTestsByTags(tests []string) ([]string, error) {
	if len(r.Config.Tags) == 0 {
		return tests, nil
	}
	res := make([]string, 0, len(tests))
	definitions, err := r.testsDefinitions(tests)
	if err != nil {
		return nil, err
	}
	definitions = tools.Filter(definitions, func(test testWithDefinition) bool {
		return !tools.Contains(test.definition.Definition.Tags, "skip")
	})
	for _, v := range r.Config.Tags {
		for _, v1 := range definitions {
			if tools.Contains(v1.definition.Definition.Tags, v) && !tools.Contains(res, v1.file) {
				res = append(res, v1.file)
			}
		}
	}

	return res, nil
}

func (r *Suite) AllTests(testPath string) ([]string, error) {
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
			nested, err := r.AllTests(testPath + "/" + v.Name())
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
