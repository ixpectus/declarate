package eval

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/maja42/goval"
)

type Eval struct {
	functions map[string]goval.ExpressionFunction
}

var evalRe = regexp.MustCompile(`\$\((.+)\)`)

func NewEval(functions map[string]goval.ExpressionFunction) *Eval {
	e := &Eval{}
	e.functions = defaultFunctions
	for k, v := range functions {
		e.functions[k] = v
	}

	return e
}

func (e *Eval) Evaluate(s string) string {
	used := usedEval(s)
	for _, val := range used {
		eval := goval.NewEvaluator()
		result, err := eval.Evaluate(val, nil, e.functions)
		if err == nil {
			s = strings.ReplaceAll(s, fmt.Sprintf("{{$(%s)}}", val), fmt.Sprintf("%v", result))
		}
	}
	return s
}

func usedEval(str string) (res []string) {
	matches := evalRe.FindAllStringSubmatch(str, -1)
	for _, match := range matches {
		res = append(res, match[1])
	}
	return res
}
