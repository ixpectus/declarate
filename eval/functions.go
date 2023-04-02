package eval

import (
	"fmt"

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
}
