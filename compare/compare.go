package compare

import (
	"fmt"
	"reflect"
	"regexp"
	"sort"

	"github.com/fatih/color"
)

type CompareParams struct {
	IgnoreValues         *bool `json:"ignoreValues,omitempty" yaml:"ignoreValues,omitempty"`
	IgnoreArraysOrdering *bool `json:"ignoreArraysOrdering,omitempty" yaml:"ignoreArraysOrdering,omitempty"`
	DisallowExtraFields  *bool `json:"disallowExtraFields,omitempty" yaml:"disallowExtraFields,omitempty"`
	AllowArrayExtraItems *bool `json:"allowArrayExtraItems,omitempty" yaml:"allowArrayExtraItems,omitempty"`
	failFast             bool  // End compare operation after first error
}

type leafsMatchType int

const (
	pure leafsMatchType = iota
	regex
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
func compare(expected, actual interface{}, params CompareParams) []error {
	errors := compareBranch("$", expected, actual, &params)
	sort.Sort(ErrorSlice(errors))
	return errors
}

func compareBranch(
	path string,
	expected, actual interface{},
	params *CompareParams,
) []error {
	expectedType := getType(expected)
	actualType := getType(actual)
	var errors []error

	// compare types
	if leafMatchType(expected) != regex && expectedType != actualType {
		errors = append(errors, MakeError(path, "types do not match", expectedType, actualType))
		return errors
	}

	// compare scalars
	if isScalarType(actualType) && (params == nil || params.IgnoreValues == nil || !*params.IgnoreValues) {
		return compareLeafs(path, expected, actual)
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
			expectedArray, actualArray = getUnmatchedArrays(expectedArray, actualArray, params)
		}

		// iterate over children
		for i, item := range expectedArray {
			subPath := fmt.Sprintf("%s[%d]", path, i)
			res := compareBranch(subPath, item, actualArray[i], params)
			errors = append(errors, res...)
			if params.failFast && len(errors) != 0 {
				return errors
			}
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
				if params.failFast {
					return errors
				}
				continue
			}

			// check values
			subPath := fmt.Sprintf("%s.%s", path, key.String())
			res := compareBranch(
				subPath,
				expectedRef.MapIndex(key).Interface(),
				actualRef.MapIndex(key).Interface(),
				params,
			)
			errors = append(errors, res...)
			if params.failFast && len(errors) != 0 {
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

func compareLeafs(path string, expected, actual interface{}) []error {
	var errors []error

	switch leafMatchType(expected) {
	case pure:
		errors = append(errors, comparePure(path, expected, actual)...)

	case regex:
		errors = append(errors, compareRegex(path, expected, actual)...)

	default:
		panic("unknown compare type")
	}

	return errors
}

func comparePure(path string, expected, actual interface{}) (errors []error) {
	if expected != actual {
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

	value := fmt.Sprintf("%v", actual)

	rx, err := regexp.Compile(retrieveRegexStr(regexExpr))
	if err != nil {
		errors = append(errors, MakeError(path, "can not compile regex", nil, "error"))
		return errors
	}

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
func getUnmatchedArrays(expected, actual []interface{}, params *CompareParams) ([]interface{}, []interface{}) {
	expectedError := make([]interface{}, 0)

	failfastParams := *params
	failfastParams.failFast = true

	for _, expectedElem := range expected {
		found := false
		for i, actualElem := range actual {
			if len(compareBranch("", expectedElem, actualElem, &failfastParams)) == 0 {
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
			if params.failFast {
				return expectedError, actual[0:1]
			}
		}
	}

	return expectedError, actual
}
