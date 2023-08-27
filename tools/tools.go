package tools

import (
	"bytes"
	"encoding/json"
	"path"
	"strings"
)

func Filter[T any](slice []T, f func(T) bool) []T {
	var res []T
	for _, e := range slice {
		if f(e) {
			res = append(res, e)
		}
	}

	return res
}

func Contains[T comparable](elems []T, v T) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}

	return false
}

func Intersect[T comparable](a, b []T) []T {
	set := make([]T, 0)
	for _, v := range a {
		if Contains(b, v) {
			set = append(set, v)
		}
	}

	return set
}

func To[T any](t T) *T {
	return &t
}

func JSONPrettyPrint(in string) string {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(in), "", "\t")
	if err != nil {
		return in
	}
	return out.String()
}

func JSONRemarshal(in string) (string, error) {
	var ifce interface{}
	err := json.Unmarshal([]byte(in), &ifce)
	if err != nil {
		return "", err
	}
	res, err := json.Marshal(ifce)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

func FilenameShort(fileName string) string {
	parts := strings.Split(fileName, "/")
	if len(parts) > 4 {
		return path.Base(fileName)
	}
	return fileName
}
