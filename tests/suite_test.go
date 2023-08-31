package tests

import (
	"os"
	"testing"

	"github.com/ixpectus/declarate/commands/echo"
	"github.com/ixpectus/declarate/commands/script"
	"github.com/ixpectus/declarate/commands/shell"
	"github.com/ixpectus/declarate/commands/vars"
	"github.com/ixpectus/declarate/compare"
	"github.com/ixpectus/declarate/contract"
	"github.com/ixpectus/declarate/eval"
	"github.com/ixpectus/declarate/kv"
	"github.com/ixpectus/declarate/output"
	"github.com/ixpectus/declarate/suite"
	"github.com/ixpectus/declarate/variables"
)

func TestSuite(t *testing.T) {
	os.Chdir("../")
	evaluator := eval.NewEval(nil)
	vv := variables.New(evaluator, kv.New("persistent"))
	cmp := compare.New(contract.CompareParams{}, vv)
	// if output
	s := suite.New("./tests/suite", suite.RunConfig{
		TestRunWrapper: NewDebugWrapper(),
		Variables:      vv,
		Output: &output.Output{
			WithProgressBar: true,
		},
		Builders: []contract.CommandBuilder{
			&echo.Unmarshaller{},
			vars.NewUnmarshaller(evaluator),
			shell.NewUnmarshaller(cmp),
			script.NewUnmarshaller(cmp),
		},
		T: t,
	})
	if err := s.Run(); err != nil {
		t.Log(err)
		t.Fail()
	}
}
