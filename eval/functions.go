package eval

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/brianvoe/gofakeit"
	"github.com/maja42/goval"
)

var defaultFunctions = map[string]goval.ExpressionFunction{
	"randName": func(args ...interface{}) (interface{}, error) {
		return gofakeit.Name(), nil
	},
	"someName": func(args ...interface{}) (interface{}, error) {
		return "Donny", nil
	},
	"getInput": func(args ...interface{}) (interface{}, error) {
		if len(args) == 0 {
			return nil, fmt.Errorf("method getInput accept one parameter")
		}
		return args[0], nil
	},
	"randFirstName": func(args ...interface{}) (interface{}, error) {
		return gofakeit.FirstName(), nil
	},
	"randInt32": func(args ...interface{}) (interface{}, error) {
		return gofakeit.Int32(), nil
	},
	"randInt64": func(args ...interface{}) (interface{}, error) {
		return gofakeit.Int64(), nil
	},
	"randFloat32": func(args ...interface{}) (interface{}, error) {
		return gofakeit.Float32(), nil
	},
	"randFloat64": func(args ...interface{}) (interface{}, error) {
		return gofakeit.Float64(), nil
	},
	"randLogin": func(args ...interface{}) (interface{}, error) {
		return fmt.Sprintf("%s%d", gofakeit.Word(), gofakeit.Int32()), nil
	},
	"notEmpty": func(args ...interface{}) (interface{}, error) {
		if len(args) == 0 {
			return false, nil
		}
		for _, v := range args {
			if empty(v) {
				return false, nil
			}
		}

		return true, nil
	},
	"any": func(args ...interface{}) (interface{}, error) {
		if len(args) == 0 {
			return false, nil
		}
		for _, v := range args {
			if empty(v) {
				return false, nil
			}
		}

		return true, nil
	},
	"intStrict": func(args ...interface{}) (interface{}, error) {
		if len(args) == 0 {
			return false, nil
		}
		_, ok := args[0].(int)
		return ok, nil
	},
	"num": func(args ...interface{}) (interface{}, error) {
		if len(args) == 0 {
			return false, nil
		}

		v := fmt.Sprintf("%v", args[0])
		_, err := strconv.ParseFloat(v, 64)
		return err == nil, nil
	},
	"empty": func(args ...interface{}) (interface{}, error) {
		if len(args) == 0 {
			return true, nil
		}
		for _, v := range args {
			if !empty(v) {
				return false, nil
			}
		}

		return true, nil
	},
	"oneOf": func(args ...interface{}) (interface{}, error) {
		if len(args) == 0 {
			return false, nil
		}
		if len(args) == 1 {
			return false, nil
		}
		val := args[len(args)-1]
		items := args[:len(args)-1]
		for _, v := range items {
			if v == val {
				return true, nil
			}
		}

		return false, nil
	},
}

func empty(val interface{}) bool {
	v := reflect.ValueOf(val)
	switch v.Kind() {
	case reflect.String, reflect.Array:
		return v.Len() == 0
	case reflect.Map, reflect.Slice:
		return v.Len() == 0 || v.IsNil()
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}

	return reflect.DeepEqual(val, reflect.Zero(v.Type()).Interface())
}
