package tests

import (
	"fmt"
	"strings"

	"github.com/ixpectus/declarate/commands/request"
	"github.com/ixpectus/declarate/contract"
)

type DebugWrapper struct{}

func NewDebugWrapper() *DebugWrapper {
	return new(DebugWrapper)
}

func (d *DebugWrapper) BeforeTest(file string, conf *contract.RunConfig, lvl int) {
	if strings.Contains(file, "yaml_wrapper") {
		for _, v := range conf.Commands {
			req, ok := v.(*request.Request)
			if ok && req.Config != nil {
				if req.Config.HeadersVal == nil {
					req.Config.HeadersVal = map[string]string{}
				}
				req.Config.HeadersVal["test-header"] = "www"
			}
		}
		fmt.Printf("before test %v\n", conf.Name)
	}
}

func (d *DebugWrapper) AfterTest(conf *contract.RunConfig, result contract.Result) {
	if strings.Contains(result.FileName, "yaml_wrapper") {
		fmt.Printf("after test %v\n", conf.Name)
	}
}

func (d *DebugWrapper) BeforeTestStep(file string, conf *contract.RunConfig, lvl int) {
	if strings.Contains(file, "yaml_wrapper") && lvl > 0 {
		for _, v := range conf.Commands {
			req, ok := v.(*request.Request)
			if ok && req.Config != nil {
				if req.Config.HeadersVal == nil {
					req.Config.HeadersVal = map[string]string{}
				}
				req.Config.HeadersVal["test-header"] = "www"
			}
		}
		fmt.Printf("before test step %v\n", conf.Name)
	}
}

func (d *DebugWrapper) AfterTestStep(conf *contract.RunConfig, result contract.Result) {
	if strings.Contains(result.FileName, "yaml_wrapper") && result.Lvl > 0 {
		fmt.Printf("after test step %v\n", conf.Name)
	}
}
