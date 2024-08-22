package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"path"
	"reflect"
	"strings"
	"time"
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

func FilenameLastN(fileName string, n int) string {
	parts := strings.Split(fileName, "/")
	return strings.Join(parts[len(parts)-n:], "/")
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

func WaitStartAPI(host string, port string) error {
	connected := false
	for i := 0; i < 5; i++ {
		connected = CheckConnect(host, port)
		if connected {
			return nil
		}
		time.Sleep(5 * time.Millisecond)
	}
	return fmt.Errorf("server not running")
}

func CheckConnect(host string, port string) bool {
	conn, err := net.Dial("tcp", net.JoinHostPort(host, port))
	if err != nil {
		return false
	}
	if conn != nil {
		defer conn.Close()
		return true
	}
	return false
}

func FormatVariables(m map[string]string) string {
	vv := []string{}
	for k, v := range m {
		vv = append(vv, fmt.Sprintf("%v:%v", k, v))
	}
	return strings.Join(vv, "\n")
}
