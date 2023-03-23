package run

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
	"github.com/ixpectus/declarate/variables"
)

func TestSkelet(t *testing.T) {
	color.NoColor = false
	runner := New(RunnerConfig{
		file:      "./config_req.yaml",
		variables: variables.New(),
		output:    &output.Output{},
		builders: []contract.CommandBuilder{
			&echo.Unmarshaller{},
			&vars.Unmarshaller{},
			request.NewUnmarshaller("http://localhost:8181/"),
		},
	})
	runner.Run()
}

func TestDb(t *testing.T) {
	color.NoColor = false
	runner := New(RunnerConfig{
		file:      "./config_db.yaml",
		variables: variables.New(),
		output:    &output.Output{},
		builders: []contract.CommandBuilder{
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
	runner := New(RunnerConfig{
		file:      "./config_shell.yaml",
		variables: variables.New(),
		output:    &output.Output{},
		builders: []contract.CommandBuilder{
			&echo.Unmarshaller{},
			&vars.Unmarshaller{},
			request.NewUnmarshaller("http://localhost:8181/"),
			db.NewUnmarshaller("postgres://user:sdlfksdjflakdf@5.188.142.25:5432/dbaas_dev?sslmode=disable"),
			shell.NewUnmarshaller(),
		},
	})
	runner.Run()
}
