package run

import (
	"testing"

	"github.com/ixpectus/declarate/commands/echo"
	"github.com/ixpectus/declarate/commands/vars"
	"github.com/ixpectus/declarate/contract"
	"github.com/ixpectus/declarate/variables"
)

func TestSkelet(t *testing.T) {
	runner := New(RunnerConfig{
		file:      "./config.yaml",
		variables: variables.New(),
		builders: []contract.CommandBuilder{
			echo.Build,
			vars.Build,
		},
	})
	runner.Run()

	// configsEcho := []RunConfigEcho{}
	// pp.ColoringEnabled = false
	// yaml.Unmarshal(file, &configsEcho)
	// pp.Println(configsEcho)
}
