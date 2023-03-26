package tests

import (
	"testing"

	"github.com/fatih/color"
	"github.com/ixpectus/declarate/commands/db"
	"github.com/ixpectus/declarate/commands/echo"
	"github.com/ixpectus/declarate/commands/request"
	"github.com/ixpectus/declarate/commands/vars"
	"github.com/ixpectus/declarate/contract"
	"github.com/ixpectus/declarate/output"
	"github.com/ixpectus/declarate/run"
	"github.com/ixpectus/declarate/suite"
	"github.com/ixpectus/declarate/variables"
)

var (
	vv     = variables.New()
	runner = run.New(run.RunnerConfig{
		Variables: vv,
		Output:    &output.Output{},
		Builders: []contract.CommandBuilder{
			&echo.Unmarshaller{},
			&vars.Unmarshaller{},
			request.NewUnmarshaller("http://localhost:8181/"),
			db.NewUnmarshaller("postgres://user:sdlfksdjflakdf@5.188.142.25:5432/dbaas_dev?sslmode=disable"),
		},
	})
)

func TestReq(t *testing.T) {
	color.NoColor = false
	runner.Run("./yaml/config_req.yaml")
	vv.Reset()
}

func TestDb(t *testing.T) {
	color.NoColor = false
	runner.Run("./yaml/config_db.yaml")
	vv.Reset()
}

func TestShell(t *testing.T) {
	color.NoColor = false
	runner.Run("./yaml/config_shell.yaml")
	vv.Reset()
}

func TestSuite(t *testing.T) {
	s := suite.New("./yaml", suite.RunConfig{
		RunAll:       false,
		Tags:         []string{},
		Filename:     []string{},
		SkipFilename: []string{},
		DryRun:       false,
		Variables:    vv,
		Output:       &output.Output{},
		Builders: []contract.CommandBuilder{
			&echo.Unmarshaller{},
			&vars.Unmarshaller{},
			request.NewUnmarshaller("http://localhost:8181/"),
			db.NewUnmarshaller("postgres://user:sdlfksdjflakdf@5.188.142.25:5432/dbaas_dev?sslmode=disable"),
		},
	})
	err := s.Run()
	if err != nil {
		t.Fail()
	}
}
