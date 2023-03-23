package compare

import (
	"encoding/json"
	"errors"
	"fmt"
)

func CompareJsonBody(expectedBody string, realBody string, params CompareParams) ([]error, error) {
	// decode expected body
	var expected interface{}
	if err := json.Unmarshal([]byte(expectedBody), &expected); err != nil {
		return nil, fmt.Errorf(
			"invalid JSON in response for test : %s",
			err.Error(),
		)
	}

	// decode actual body
	var actual interface{}
	if err := json.Unmarshal([]byte(realBody), &actual); err != nil {
		return []error{errors.New("could not parse response")}, nil
	}

	return Compare(expected, actual, params), nil
}
