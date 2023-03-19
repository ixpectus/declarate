package run

import (
	"fmt"

	"github.com/ixpectus/declarate/contract"
)

func Run(cc []RunConfig, vv contract.Vars) {
	for _, v := range cc {
		runOne(v, 0, vv)
	}
}

func runOne(conf RunConfig, lvl int, vv contract.Vars) {
	prefix := ""
	for i := 0; i < lvl; i++ {
		prefix += " "
	}
	if conf.Name != "" {
		fmt.Printf(prefix+"run test with name %v\n", conf.Name)
	}
	for _, c := range conf.Commands {
		c.SetVars(vv)
		c.Do()
		if err := c.Check(); err != nil {
			fmt.Printf(prefix+"test failed %v\n", err)
		}
	}
	if len(conf.Steps) > 0 {
		for _, v := range conf.Steps {
			runOne(v, lvl+1, vv)
		}
	}
}
