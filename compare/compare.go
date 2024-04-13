package compare

import (
	"fmt"
	"math"
	"reflect"
	"regexp"
	"sort"
	"strings"

	"github.com/fatih/color"
	conditionRule "github.com/ixpectus/declarate/condition"
	"github.com/ixpectus/declarate/contract"
	"github.com/ixpectus/declarate/eval"
	"github.com/ixpectus/declarate/tools"
)

type leafsMatchType int

const (
	pure leafsMatchType = iota
	regex
	condition
)

type ErrorSlice []error

func (x ErrorSlice) Len() int           { return len(x) }
func (x ErrorSlice) Less(i, j int) bool { return x[i].Error() < x[j].Error() }
func (x ErrorSlice) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

// Sort is a convenience method: x.Sort() calls Sort(x).
func (x ErrorSlice) Sort() { sort.Sort(x) }

var regexExprRx = regexp.MustCompile(`^\$matchRegexp\((.+)\)$`)

// compare compares values as plain text
// It can be compared several ways:
//   - Pure values: should be equal
//   - Regex: try to compile 'expected' as regex and match 'actual' with it
//     It activates on following syntax: $matchRegexp(%EXPECTED_VALUE%)
func (c *Comparer) compare(expected, actual interface{}, params contract.CompareParams) []error {
	errors := c.compareBranch("$", expected, actual, &params)
	sort.Sort(ErrorSlice(errors))
	return errors
}

func (c *Comparer) compareBranch(
	path string,
	expected, actual interface{},
	params *contract.CompareParams,
) []error {
	expectedType := getType(expected)

	actualType := getType(actual)

	var errors []error

	// compare types
	if leafMatchType(expected) != regex && leafMatchType(expected) != condition && expectedType != actualType {
		errors = append(errors, MakeError(path, "types do not match", expectedType, actualType))
		return errors
	}

	// compare scalars
	if isScalarType(actualType) && (params == nil || params.IgnoreValues == nil || !*params.IgnoreValues) {
		return c.compareLeafs(path, expected, actual)
	}

	// compare arrays
	if actualType == "array" {
		expectedArray := convertToArray(expected)
		actualArray := convertToArray(actual)

		if (params == nil || params.AllowArrayExtraItems == nil || !*params.AllowArrayExtraItems) && len(expectedArray) != len(actualArray) {
			errors = append(errors, MakeError(path, "array lengths do not match", len(expectedArray), len(actualArray)))
			return errors
		}

		if (params.IgnoreArraysOrdering != nil && *params.IgnoreArraysOrdering) ||
			(params.AllowArrayExtraItems != nil && *params.AllowArrayExtraItems) {
			expectedArray, actualArray = c.getUnmatchedArrays(expectedArray, actualArray, params)
		}
		if len(actualArray) < len(expectedArray) {
			errors = append(errors, MakeError(path, "array lengths do not match", len(expectedArray), len(actualArray)))
			return errors
		}

		// iterate over children
		for i, item := range expectedArray {
			subPath := fmt.Sprintf("%s[%d]", path, i)
			res := c.compareBranch(subPath, item, actualArray[i], params)
			errors = append(errors, res...)
			if params.FailFast && len(errors) != 0 {
				return errors
			}
		}
		if len(errors) > 0 {
			return errors
		}
	}

	// compare maps
	if actualType == "map" {
		expectedRef := reflect.ValueOf(expected)
		actualRef := reflect.ValueOf(actual)

		if (params.DisallowExtraFields != nil && *params.DisallowExtraFields) && expectedRef.Len() != actualRef.Len() {
			errors = append(errors, MakeError(path, "map lengths do not match", expectedRef.Len(), actualRef.Len()))
			return errors
		}

		for _, key := range expectedRef.MapKeys() {
			// check keys presence
			if ok := actualRef.MapIndex(key); !ok.IsValid() {
				errors = append(errors, MakeError(path, "key is missing", key.String(), "<missing>"))
				if params.FailFast {
					return errors
				}
				continue
			}

			// check values
			subPath := fmt.Sprintf("%s.%s", path, key.String())
			res := c.compareBranch(
				subPath,
				expectedRef.MapIndex(key).Interface(),
				actualRef.MapIndex(key).Interface(),
				params,
			)
			errors = append(errors, res...)
			if params.FailFast && len(errors) != 0 {
				return errors
			}
		}
	}

	return errors
}

func getType(value interface{}) string {
	if value == nil {
		return "nil"
	}
	rt := reflect.TypeOf(value)
	if rt.Kind() == reflect.Slice || rt.Kind() == reflect.Array {
		return "array"
	} else if rt.Kind() == reflect.Map {
		return "map"
	} else {
		return rt.String()
	}
}

func isScalarType(t string) bool {
	return !(t == "array" || t == "map")
}

func (c *Comparer) compareLeafs(path string, expected, actual interface{}) []error {
	var errors []error

	switch leafMatchType(expected) {
	case pure:
		errors = append(errors, c.comparePure(path, expected, actual)...)

	case regex:
		errors = append(errors, compareRegex(path, expected, actual)...)

	case condition:
		errors = append(errors, c.compareCondition(path, expected, actual)...)

	default:
		panic("unknown compare type")
	}

	return errors
}

func (c *Comparer) compareCondition(path string, expected, actual interface{}) (errors []error) {
	expr, ok := expected.(string)
	if !ok {
		errors = append(errors, MakeError(path, "type mismatch", "string", reflect.TypeOf(expected)))
		return errors
	}
	expr = fmt.Sprintf("$(%s)", strings.TrimLeft(expr, "$"))

	if strings.Contains(expr, "()") {
		if tools.IsNumber(actual) {
			expr = strings.ReplaceAll(expr, "()", fmt.Sprintf("(%v)", actual))
		} else {
			expr = strings.ReplaceAll(expr, "()", fmt.Sprintf("(\"%v\")", actual))
		}
	} else if strings.Contains(expr, "))") {
		if tools.IsNumber(actual) {
			expr = strings.Replace(expr, "))", fmt.Sprintf(", %v))", actual), 1)
		} else {
			expr = strings.Replace(expr, "))", fmt.Sprintf(", \"%v\"))", actual), 1)
		}
	} else {
		if tools.IsNumber(actual) {
			expr = strings.Replace(expr, ")", fmt.Sprintf("(%v))", actual), 1)
		} else {
			expr = strings.Replace(expr, ")", fmt.Sprintf("(\"%v\"))", actual), 1)
		}
	}

	if !conditionRule.IsTrueNoWrap(c.vars, expr) {
		errors = append(errors, MakeError(path, "values do not match by condition", expected, actual))
	}

	return errors
}

func (c *Comparer) comparePure(path string, expected, actual interface{}) (errors []error) {
	expectedStr, expectedIsString := expected.(string)
	actualStr, actualIsString := actual.(string)

	tryLineByLine := expectedIsString && actualIsString && (c.defaultComparisonParams.LineByLine == nil || *c.defaultComparisonParams.LineByLine)
	linesExpected := strings.Split(expectedStr, "\n")
	linesGot := strings.Split(actualStr, "\n")

	if tryLineByLine && len(linesExpected) > 1 && len(linesGot) > 1 {
		if len(linesExpected) != len(linesGot) {
			errMsg := fmt.Sprintf("lines count differs, expected %v, got %v", len(linesExpected), len(linesGot))
			for k := range linesExpected {
				if len(linesGot) > k {
					if linesExpected[k] != linesGot[k] {
						errMsg += fmt.Sprintf("\nlines different at line %v, expected `%v`, got `%v`", k, linesExpected[k], linesGot[k])
					}
				}
			}
			res := MakeError(
				path,
				errMsg,
				expectedStr,
				actualStr,
			)
			return []error{res}
		}
		for k := range linesExpected {
			if len(linesGot) >= k {
				if strings.Trim(linesExpected[k], " ") != strings.Trim(linesGot[k], " ") {
					errMsg := fmt.Sprintf("\nlines different at line %v, expected `%v`, got `%v`", k, linesExpected[k], linesGot[k])
					return []error{MakeError(path, errMsg, expected, actual)}
				}
			}
		}
		return nil
	} else if expected != actual {
		errors = append(errors, MakeError(path, "values do not match", expected, actual))
	}

	return errors
}

func compareRegex(path string, expected, actual interface{}) (errors []error) {
	regexExpr, ok := expected.(string)
	if !ok {
		errors = append(errors, MakeError(path, "type mismatch", "string", reflect.TypeOf(expected)))
		return errors
	}

	rx, err := regexp.Compile(retrieveRegexStr(regexExpr))
	if err != nil {
		errors = append(errors, MakeError(path, "can not compile regex", nil, "error"))
		return errors
	}

	f, ok := actual.(float64)

	if ok && math.Mod(f, 1) == 0 {
		value := fmt.Sprintf("%v", int(f))
		if !rx.MatchString(value) {
			errors = append(errors, MakeError(path, "value does not match regex", expected, value))
			return errors
		}

		return nil
	}
	value := fmt.Sprintf("%v", actual)

	if !rx.MatchString(value) {
		errors = append(errors, MakeError(path, "value does not match regex", expected, actual))
		return errors
	}

	return nil
}

func retrieveRegexStr(expr string) string {
	if matches := regexExprRx.FindStringSubmatch(expr); matches != nil {
		return matches[1]
	}

	return ""
}

func leafMatchType(expected interface{}) leafsMatchType {
	val, ok := expected.(string)

	if !ok {
		return pure
	}

	if matches := regexExprRx.FindStringSubmatch(val); matches != nil {
		return regex
	}

	if eval.HasEvalResponse(val) {
		return condition
	}

	return pure
}

func MakeError(path, msg string, expected, actual interface{}) error {
	if path != "" {
		return fmt.Errorf(
			"at path %s %s:\nexpected:\n%s\nactual:\n%s",
			color.CyanString(path),
			msg,
			color.GreenString("%v", expected),
			color.RedString("%v", actual),
		)
	}
	return fmt.Errorf(
		"%s:\n     expected: \n%s\n       actual: \n%s",
		msg,
		color.GreenString("%v", expected),
		color.RedString("%v", actual),
	)
}

func convertToArray(array interface{}) []interface{} {
	ref := reflect.ValueOf(array)

	interfaceSlice := make([]interface{}, 0)
	for i := 0; i < ref.Len(); i++ {
		interfaceSlice = append(interfaceSlice, ref.Index(i).Interface())
	}
	return interfaceSlice
}

// For every elem in "expected" try to find elem in "actual". Returns arrays without matching.
func (c *Comparer) getUnmatchedArrays(expected, actual []interface{}, params *contract.CompareParams) ([]interface{}, []interface{}) {
	expectedError := make([]interface{}, 0)

	failfastParams := *params
	failfastParams.FailFast = true

	for _, expectedElem := range expected {
		found := false
		for i, actualElem := range actual {
			if len(c.compareBranch("", expectedElem, actualElem, &failfastParams)) == 0 {
				// expectedElem match actualElem
				found = true
				// remove actualElem from  actual
				if len(actual) != 1 {
					actual[i] = actual[len(actual)-1]
				}
				actual = actual[:len(actual)-1]
				break
			}
		}
		if !found {
			expectedError = append(expectedError, expectedElem)
			if params.FailFast {
				return expectedError, actual[0:1]
			}
		}
	}

	return expectedError, actual
}
