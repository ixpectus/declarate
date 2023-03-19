package run

import (
	"os"
	"testing"

	"github.com/ixpectus/declarate/variables"
	"github.com/k0kubun/pp"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

func TestSkelet(t *testing.T) {
	file, err := os.ReadFile("./config.yaml")
	require.NoError(t, err)
	vars := variables.New()
	configs := []RunConfig{}
	pp.ColoringEnabled = false
	yaml.Unmarshal(file, &configs)
	Run(configs, vars)

	// configsEcho := []RunConfigEcho{}
	// pp.ColoringEnabled = false
	// yaml.Unmarshal(file, &configsEcho)
	// pp.Println(configsEcho)
}
