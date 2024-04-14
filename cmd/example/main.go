package main

import (
	"flag"
	"log"
	"strings"

	"github.com/ixpectus/declarate/defaults"
	"github.com/ixpectus/declarate/report"
	"github.com/ixpectus/declarate/tests"
	"github.com/ixpectus/declarate/tools"
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
	flagFailFast = flag.Bool(
		"failFast", false, "fail after first fail test",
	)
	flagTags = flag.String(
		"tags",
		"",
		"tags for filter tags, example `-tags tag1,tag2`",
	)
	flagContinue = flag.Bool(
		"continue",
		false,
		"continue last test execution",
	)
	flagTests = flag.String(
		"tests",
		"",
		"test files, example `-tests config,db`",
	)
	flagWithProgressBar = flag.Bool(
		"progress_bar",
		false,
		"progress bar for poll interval",
	)
	flagOutput = flag.String(
		"output",
		"print",
		"test output log or print",
	)

	flagClearPersistent = flag.Bool(
		"clear",
		false,
		"clear persistent",
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
	if err := tools.WaitStartAPI("127.0.0.1", "8181"); err != nil {
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

	tags := []string{}
	filePathes := []string{}
	if *flagTags != "" {
		tags = strings.Split(*flagTags, ",")
	}
	if *flagTests != "" {
		filePathes = strings.Split(*flagTests, ",")
	}
	s := defaults.NewDefaultSuite(defaults.SuiteConfig{
		Dir:             *flagDir,
		NoColor:         true,
		DefaultDBConn:   "postgres://postgres@127.0.0.1:5440/?sslmode=disable",
		SkipTests:       coreTestsToSkip,
		ClearPersistent: *flagClearPersistent,
		DryRun:          *flagDryRun,
		WithProgresBar:  *flagWithProgressBar,
		DefaultHost:     "http://127.0.0.1:8181/",
		Wrapper:         tests.NewDebugWrapper(),
		Report:          report.NewEmptyReport(),
		Continue:        *flagContinue,
		Tags:            tags,
		Filepathes:      filePathes,
		AllPersistent:   true,
	})
	if err := s.Run(); err != nil {
		log.Println(err)
	}
}
