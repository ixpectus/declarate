package tests

import (
	"testing"

	"github.com/fatih/color"
	"github.com/ixpectus/declarate/commands/db"
	"github.com/ixpectus/declarate/commands/echo"
	"github.com/ixpectus/declarate/commands/request"
	"github.com/ixpectus/declarate/commands/vars"
	"github.com/ixpectus/declarate/compare"
	"github.com/ixpectus/declarate/contract"
	"github.com/ixpectus/declarate/eval"
	"github.com/ixpectus/declarate/kv"
	"github.com/ixpectus/declarate/output"
	"github.com/ixpectus/declarate/run"
	"github.com/ixpectus/declarate/variables"
)

var (
	evaluator  = eval.NewEval(nil)
	vv         = variables.New(evaluator, kv.New("persistent", false))
	cmp        = compare.New(contract.CompareParams{}, vv)
	connLoader = db.NewPGLoader("postgres://postgres@127.0.0.1:5440/?sslmode=disable")
	runner     = run.New(run.RunnerConfig{
		Variables: vv,
		Output:    &output.Output{},
		Builders: []contract.CommandBuilder{
			&echo.Unmarshaller{},
			&vars.Unmarshaller{},
			request.NewUnmarshaller("http://localhost:8181/", cmp),
			db.NewUnmarshaller(connLoader, cmp),
		},
	})
)

func TestReq(t *testing.T) {
	color.NoColor = false
	runner.Run("./yaml/config_req.yaml", nil)
	vv.Reset()
}

func TestDb(t *testing.T) {
	color.NoColor = false
	runner.Run("./yaml/db.yaml", nil)
	vv.Reset()
}

func TestShell(t *testing.T) {
	color.NoColor = false
	runner.Run("./yaml/config_shell.yaml", nil)
	vv.Reset()
}
