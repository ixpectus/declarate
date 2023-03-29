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
	"github.com/ixpectus/declarate/commands/vars"
	"github.com/ixpectus/declarate/contract"
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
		"tags for filter tags, example `tags tag1,tag2`",
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
	vv := variables.New()
	s := suite.New(*flagDir, suite.RunConfig{
		RunAll:       false,
		Filename:     []string{},
		SkipFilename: coreTestsToSkip,
		DryRun:       *flagDryRun,
		Variables:    vv,
		Output:       &output.OutputPrintln{},
		Builders: []contract.CommandBuilder{
			&echo.Unmarshaller{},
			&vars.Unmarshaller{},
			request.NewUnmarshaller("http://localhost:8181/"),
			db.NewUnmarshaller("postgres://postgres@127.0.0.1:5440/?sslmode=disable"),
		},
	})

	if *flagTags != "" {
		s.Config.Tags = strings.Split(*flagTags, ",")
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
