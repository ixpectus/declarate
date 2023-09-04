package defaults

import (
	"testing"

	"github.com/ixpectus/declarate/commands/db"
	"github.com/ixpectus/declarate/commands/echo"
	"github.com/ixpectus/declarate/commands/request"
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

type SuiteConfig struct {
	Dir             string
	SkipTests       []string
	DryRun          bool
	ClearPersistent bool
	WithProgresBar  bool
	DefaultDBConn   string
	DefaultHost     string
	Tags            []string
	Filepathes      []string
	NoColor         bool
	Wrapper         contract.TestWrapper
	T               *testing.T
	Output          contract.Output
	Report          contract.Report
	FailFast        bool
}

func NewDefaultSuite(conf SuiteConfig) *suite.Suite {
	evaluator := eval.NewEval(nil)
	vv := variables.New(evaluator, kv.New("persistent", conf.ClearPersistent))
	cmp := compare.New(contract.CompareParams{}, vv)
	connLoader := db.NewPGLoader(conf.DefaultDBConn)

	var out contract.Output
	out = &output.OutputPrintln{
		WithProgressBar: conf.WithProgresBar,
	}
	if conf.Output != nil {
		out = conf.Output
	}
	s := suite.New(conf.Dir, suite.RunConfig{
		RunAll:         false,
		NoColor:        conf.NoColor,
		SkipFilename:   conf.SkipTests,
		TestRunWrapper: conf.Wrapper,
		DryRun:         conf.DryRun,
		Variables:      vv,
		Tags:           conf.Tags,
		T:              conf.T,
		Filepathes:     conf.Filepathes,
		Report:         conf.Report,
		Output:         out,
		Builders: []contract.CommandBuilder{
			&echo.Unmarshaller{},
			vars.NewUnmarshaller(evaluator),
			shell.NewUnmarshaller(cmp),
			script.NewUnmarshaller(cmp),
			request.NewUnmarshaller(conf.DefaultHost, cmp),
			db.NewUnmarshaller(connLoader, cmp),
		},
	})

	return s
}
