package tools

import (
	"bytes"
	"encoding/json"
	"path"
	"reflect"
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

func IsNumber(val interface{}) bool {
	v := reflect.ValueOf(val)
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return true
	case reflect.Float32, reflect.Float64:
		return true
	}
	return false
}
