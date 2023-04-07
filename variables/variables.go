package variables

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/brianvoe/gofakeit"
	"github.com/ixpectus/declarate/contract"
	"github.com/maja42/goval"
	"github.com/tidwall/gjson"
)

var variableRx = regexp.MustCompile(`{{\s*\$(\w+)\s*}}`)

type Variables struct {
	data      map[string]string
	eval      contract.Evaluator
	functions map[string]goval.ExpressionFunction
}

func New(evaluator contract.Evaluator) *Variables {
	gofakeit.Seed(0)
	vv := &Variables{
		data: map[string]string{},
		eval: evaluator,
	}
	return vv
}

func (v *Variables) Set(k, val string) {
	if strings.ToUpper(k) == k {
		os.Setenv(k, val)
	}
	v.data[k] = val
}

func (v *Variables) Reset() {
	v.data = map[string]string{}
}

func (v *Variables) Get(k string) string {
	if v, ok := v.data[k]; ok {
		return v
	}
	return os.Getenv(k)
}

func (v *Variables) Apply(text string) string {
	used := usedVariables(text)
	for _, val := range used {
		text = strings.ReplaceAll(text, "{{$"+val+"}}", v.Get(val))
	}
	text = v.eval.Evaluate(text)
	return text
}

func usedVariables(str string) (res []string) {
	matches := variableRx.FindAllStringSubmatch(str, -1)
	for _, match := range matches {
		res = append(res, match[1])
	}
	return res
}

func FromJSON(vv map[string]string, body string) (map[string]string, error) {
	names := make([]string, 0, len(vv))
	paths := make([]string, 0, len(vv))
	for k, v := range vv {
		names = append(names, k)
		paths = append(paths, v)
	}
	vars := map[string]string{}
	results := gjson.GetMany(body, paths...)

	for n, res := range results {
		if !res.Exists() {
			return nil,
				fmt.Errorf("path '%s' doesn't exist in given json %s", paths[n], body)
		}
		vars[names[n]] = res.String()
	}

	return vars, nil
}
