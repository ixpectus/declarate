package tests

import (
	"testing"

	"github.com/fatih/color"
	"github.com/ixpectus/declarate/commands/db"
	"github.com/ixpectus/declarate/commands/echo"
	"github.com/ixpectus/declarate/commands/request"
	"github.com/ixpectus/declarate/commands/shell"
	"github.com/ixpectus/declarate/commands/vars"
	"github.com/ixpectus/declarate/contract"
	"github.com/ixpectus/declarate/output"
	"github.com/ixpectus/declarate/run"
	"github.com/ixpectus/declarate/variables"
)

func TestReq(t *testing.T) {
	color.NoColor = false
	runner := run.New(run.RunnerConfig{
		File:      "./yaml/config_req.yaml",
		Variables: variables.New(),
		Output:    &output.Output{},
		Builders: []contract.CommandBuilder{
			&echo.Unmarshaller{},
			&vars.Unmarshaller{},
			request.NewUnmarshaller("http://localhost:8181/"),
			db.NewUnmarshaller("postgres://user:sdlfksdjflakdf@5.188.142.25:5432/dbaas_dev?sslmode=disable"),
		},
	})
	runner.Run()
}

func TestDb(t *testing.T) {
	color.NoColor = false
	runner := run.New(run.RunnerConfig{
		File:      "./yaml/config_db.yaml",
		Variables: variables.New(),
		Output:    &output.Output{},
		Builders: []contract.CommandBuilder{
			&echo.Unmarshaller{},
			&vars.Unmarshaller{},
			request.NewUnmarshaller("http://localhost:8181/"),
			db.NewUnmarshaller("postgres://user:sdlfksdjflakdf@5.188.142.25:5432/dbaas_dev?sslmode=disable"),
		},
	})
	runner.Run()
}

func TestShell(t *testing.T) {
	color.NoColor = false
	runner := run.New(run.RunnerConfig{
		File:      "./yaml/config_shell.yaml",
		Variables: variables.New(),
		Output:    &output.Output{},
		Builders: []contract.CommandBuilder{
			&echo.Unmarshaller{},
			&vars.Unmarshaller{},
			request.NewUnmarshaller("http://localhost:8181/"),
			db.NewUnmarshaller("postgres://user:sdlfksdjflakdf@5.188.142.25:5432/dbaas_dev?sslmode=disable"),
			shell.NewUnmarshaller(),
		},
	})
	runner.Run()
}
