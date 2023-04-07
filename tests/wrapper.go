package tests

import (
	"strings"

	"github.com/ixpectus/declarate/contract"
	"github.com/k0kubun/pp"
)

type DebugWrapper struct{}

func NewDebugWrapper() *DebugWrapper {
	return new(DebugWrapper)
}

func (d *DebugWrapper) BeforeTest(file string, conf contract.RunConfig, lvl int) {
	if strings.Contains(file, "yaml_wrapper") {
		pp.Println(file, conf)
	}
}

func (d *DebugWrapper) AfterTest(conf contract.RunConfig, result contract.Result) {
	if strings.Contains(result.FileName, "yaml_wrapper") {
		pp.Println(conf, result)
	}
}
