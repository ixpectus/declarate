package compare

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ixpectus/declarate/contract"
)

func (c *Comparer) compareJsonBody(expectedBody string, realBody string, params contract.CompareParams) ([]error, error) {
	// decode expected body
	var expected interface{}
	if err := json.Unmarshal([]byte(expectedBody), &expected); err != nil {
		return nil, fmt.Errorf(
			"invalid JSON in response for test, json %s err : %s",
			expectedBody,
			err.Error(),
		)
	}

	// decode actual body
	var actual interface{}
	if err := json.Unmarshal([]byte(realBody), &actual); err != nil {
		return []error{errors.New("could not parse response")}, nil
	}

	return c.compare(expected, actual, params), nil
}
