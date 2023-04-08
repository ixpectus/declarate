package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/fatih/color"
	"github.com/ixpectus/declarate/commands/db"
	"github.com/ixpectus/declarate/commands/echo"
	"github.com/ixpectus/declarate/commands/request"
	"github.com/ixpectus/declarate/commands/script"
	"github.com/ixpectus/declarate/commands/shell"
	"github.com/ixpectus/declarate/commands/vars"
	"github.com/ixpectus/declarate/compare"
	"github.com/ixpectus/declarate/contract"
	"github.com/ixpectus/declarate/output"
	"github.com/ixpectus/declarate/run"
)

func TestExample(t *testing.T) {
	os.Chdir("../")
	color.NoColor = false
	cmp := compare.New(compare.CompareParams{})
	runner = run.New(run.RunnerConfig{
		Variables: vv,
		Output:    &output.Output{},
		Builders: []contract.CommandBuilder{
			shell.NewUnmarshaller(),
			script.NewUnmarshaller(),
			&echo.Unmarshaller{},
			&vars.Unmarshaller{},
			request.NewUnmarshaller("http://localhost:8181/"),
			db.NewUnmarshaller("postgres://postgres@127.0.0.1:5440/?sslmode=disable", cmp),
		},
	})
	err := runner.Run("./tests/yaml_example/example.yaml")
	if err != nil {
		fmt.Println(err)
	}
}
