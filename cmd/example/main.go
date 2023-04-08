package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/ixpectus/declarate/commands/db"
	"github.com/ixpectus/declarate/commands/echo"
	"github.com/ixpectus/declarate/commands/request"
	"github.com/ixpectus/declarate/commands/script"
	"github.com/ixpectus/declarate/commands/shell"
	"github.com/ixpectus/declarate/commands/vars"
	"github.com/ixpectus/declarate/compare"
	"github.com/ixpectus/declarate/contract"
	"github.com/ixpectus/declarate/eval"
	"github.com/ixpectus/declarate/output"
	"github.com/ixpectus/declarate/suite"
	"github.com/ixpectus/declarate/tests"
	"github.com/ixpectus/declarate/variables"
)

var (
	coreTestsToRun  stringList
	coreTestsToSkip stringList
	flagDir         = flag.String(
		"dir", "./tests/yaml", "tests directory",
	)
	flagDryRun = flag.Bool(
		"dryRun", false, "show tests for run, don't run them",
	)
	flagTags = flag.String(
		"tags",
		"",
		"tags for filter tags, example `-tags tag1,tag2`",
	)

	flagTests = flag.String(
		"tests",
		"",
		"test files, example `-tests config,db`",
	)
)

type stringList []string

func (f *stringList) String() string {
	return ""
}

func (f *stringList) Set(value string) error {
	*f = append(*f, value)
	return nil
}

func main() {
	go tests.Handle()
	color.NoColor = true
	if err := waitStartAPI("127.0.0.1", "8181"); err != nil {
		log.Fatal(err)
	}
	flag.Var(
		&coreTestsToRun,
		"t",
		"test files to run",
	)
	flag.Var(
		&coreTestsToSkip,
		"s",
		"test to skip",
	)
	flag.Parse()
	evaluator := eval.NewEval(nil)
	vv := variables.New(evaluator)
	cmp := compare.New(compare.CompareParams{})
	s := suite.New(*flagDir, suite.RunConfig{
		RunAll:         false,
		Filepathes:     []string{},
		SkipFilename:   coreTestsToSkip,
		TestRunWrapper: tests.NewDebugWrapper(),
		DryRun:         *flagDryRun,
		Variables:      vv,
		Output:         &output.OutputPrintln{},
		Builders: []contract.CommandBuilder{
			&echo.Unmarshaller{},
			vars.NewUnmarshaller(evaluator),
			shell.NewUnmarshaller(cmp),
			script.NewUnmarshaller(cmp),
			request.NewUnmarshaller("http://localhost:8181/", cmp),
			db.NewUnmarshaller("postgres://postgres@127.0.0.1:5440/?sslmode=disable", cmp),
		},
	})

	if *flagTags != "" {
		s.Config.Tags = strings.Split(*flagTags, ",")
	}
	if *flagTests != "" {
		s.Config.Filepathes = strings.Split(*flagTests, ",")
	}
	if err := s.Run(); err != nil {
		log.Println(err)
	}
}

func waitStartAPI(host string, port string) error {
	connected := false
	for i := 0; i < 5; i++ {
		connected = checkConnect(host, port)
		if connected {
			return nil
		}
		time.Sleep(5 * time.Millisecond)
	}
	return fmt.Errorf("server not running")
}

func checkConnect(host string, port string) bool {
	conn, err := net.Dial("tcp", net.JoinHostPort(host, port))
	if err != nil {
		return false
	}
	if conn != nil {
		defer conn.Close()
		return true
	}
	return false
}
