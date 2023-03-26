package tests

import (
	"fmt"
	"os"
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
)

func TestExample(t *testing.T) {
	os.Chdir("../")
	color.NoColor = false
	runner = run.New(run.RunnerConfig{
		Variables: vv,
		Output:    &output.Output{},
		Builders: []contract.CommandBuilder{
			shell.NewUnmarshaller(),
			&echo.Unmarshaller{},
			&vars.Unmarshaller{},
			request.NewUnmarshaller("http://localhost:8181/"),
			db.NewUnmarshaller("postgres://user:sdlfksdjflakdf@5.188.142.25:5432/dbaas_dev?sslmode=disable"),
		},
	})
	err := runner.Run("./tests/yaml_example/example.yaml")
	if err != nil {
		fmt.Println(err)
	}
}