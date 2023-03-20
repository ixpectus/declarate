package variables

import (
	"os"
	"regexp"
	"strings"
)

var variableRx = regexp.MustCompile(`{{\s*\$(\w+)\s*}}`)

type Variables struct {
	data map[string]string
}

func New() *Variables {
	return &Variables{
		data: map[string]string{},
	}
}

func (v *Variables) Set(k, val string) {
	v.data[k] = val
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
	return text
}

func usedVariables(str string) (res []string) {
	matches := variableRx.FindAllStringSubmatch(str, -1)
	for _, match := range matches {
		res = append(res, match[1])
	}
	return res
}
