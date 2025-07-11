package variables

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/brianvoe/gofakeit"
	"github.com/tidwall/gjson"

	"github.com/ixpectus/declarate/contract"
)

var VariableRx = regexp.MustCompile(`{{\s*\$(\w+)\s*}}`)

type Variables struct {
	data          map[string]string
	eval          contract.Evaluator
	persistent    persistent
	allPersistent bool
}

func New(
	evaluator contract.Evaluator,
	persistent persistent,
	allPersistent bool,
) *Variables {
	gofakeit.Seed(0)
	vv := &Variables{
		data:          map[string]string{},
		eval:          evaluator,
		persistent:    persistent,
		allPersistent: allPersistent,
	}

	return vv
}

func (v *Variables) Set(k, val string) error {
	if v.allPersistent {
		return v.SetPersistent(k, val)
	}
	val = v.Apply(val)
	val = v.eval.Evaluate(val)
	if strings.ToUpper(k) == k {
		os.Setenv(k, val)
	}
	v.data[k] = val

	return nil
}

func reorderMap(mapData map[string]string) []string {
	cmpFn := func(a, b string) bool {
		countA := strings.Count(a, "{{") + strings.Count(a, "$(")
		countB := strings.Count(b, "{{") + strings.Count(b, "$(")
		if countA != countB {
			return countA > countB
		} else if countA == 0 && countB == 0 {
			return false
		}

		return true
	}

	keys := make([]string, len(mapData))
	i := 0
	for key := range mapData {
		keys[i] = key
		i++
	}
	sort.Slice(keys, func(i, j int) bool {
		return !cmpFn(mapData[keys[i]], mapData[keys[j]])
	})

	return keys
}

func (v *Variables) SetAll(m map[string]string) (map[string]string, error) {
	res := map[string]string{}
	keys := reorderMap(m)
	for _, k := range keys {
		val := m[k]
		if err := v.Set(k, val); err != nil {
			return nil, fmt.Errorf("set key %v, value %v: %w", k, val, err)
		}
		if !v.allPersistent {
			res[k] = v.data[k]
		} else {
			val, _ := v.persistent.Get(k)
			res[k] = val
		}
	}

	return res, nil
}

func (v *Variables) SetPersistent(k, val string) error {
	val = v.Apply(val)
	val = v.eval.Evaluate(val)
	if strings.ToUpper(k) == k {
		os.Setenv(k, val)
	}

	return v.persistent.Set(k, val)
}

func (v *Variables) Reset() {
	v.data = map[string]string{}
}

func (v *Variables) Get(k string) string {
	if v, ok := v.data[k]; ok {
		return v
	}
	if v.persistent != nil {
		res, _ := v.persistent.Get(k)
		if res != "" {
			return res
		}
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
	matches := VariableRx.FindAllStringSubmatch(str, -1)
	for _, match := range matches {
		res = append(res, match[1])
	}
	return res
}

func FromJSON(vv map[string]string, body string, existedVars contract.Vars) (map[string]string, error) {
	names := make([]string, 0, len(vv))
	paths := make([]string, 0, len(vv))
	for k, v := range vv {
		names = append(names, existedVars.Apply(k))
		paths = append(paths, existedVars.Apply(v))
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
