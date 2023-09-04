package tests

import (
	"flag"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/ixpectus/declarate/defaults"
	"github.com/ixpectus/declarate/output"
	"github.com/ixpectus/declarate/report"
	"github.com/ixpectus/declarate/tools"
)

var flagDir = flag.String(
	"dir", "./tests/yaml", "tests directory",
)

var flagClearPersistent = flag.Bool(
	"clear",
	false,
	"clear persistent",
)

var flagTests = flag.String(
	"tests",
	"",
	"test files, example `-tests config,db`",
)

func TestAllure(t *testing.T) {
	os.Chdir("../")
	go Handle()
	if err := tools.WaitStartAPI("127.0.0.1", "8181"); err != nil {
		log.Fatal(err)
	}
	filepathes := []string{}
	if *flagTests != "" {
		filepathes = strings.Split(*flagTests, ",")
	}
	s := defaults.NewDefaultSuite(defaults.SuiteConfig{
		Dir:             *flagDir,
		Report:          report.NewAllureReport("./allure-results"),
		DefaultDBConn:   "postgres://postgres@127.0.0.1:5440/?sslmode=disable",
		ClearPersistent: *flagClearPersistent,
		WithProgresBar:  true,
		Filepathes:      filepathes,
		Output: &output.Output{
			WithProgressBar: true,
		},
		DefaultHost: "http://localhost:8181/",
		T:           t,
	})
	if err := s.Run(); err != nil {
		t.Log(err)
		t.Fail()
	}
}
