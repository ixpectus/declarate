package main

import (
	"fmt"
	"sort"
	"strings"
)

// Function to format variables with ordered keys
func FormatVariablesWithOrderedKeys(m map[string]string) string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	vv := make([]string, len(keys))
	for i, k := range keys {
		vv[i] = fmt.Sprintf("%s:%s", k, m[k])
	}
	return strings.Join(vv, "\n")
}

func main() {
	inputMap := map[string]string{"c": "apple", "a": "banana", "d": "orange"}
	output := FormatVariablesWithOrderedKeys(inputMap)
	fmt.Println("Input Map:", inputMap)
	fmt.Println("Output:", output)
}
