package compare

import (
	"github.com/ixpectus/declarate/contract"
)

type Comparer struct {
	defaultComparisonParams contract.CompareParams
	vars                    contract.Vars
}

func New(c contract.CompareParams, vars contract.Vars) *Comparer {
	return &Comparer{defaultComparisonParams: c, vars: vars}
}

func (c *Comparer) Compare(expected, actual interface{}, params contract.CompareParams) []error {
	return c.compare(expected, actual, c.merge(params))
}

func (c *Comparer) CompareJsonBody(
	expectedBody string,
	realBody string,
	params contract.CompareParams,
) ([]error, error) {
	return c.compareJsonBody(expectedBody, realBody, c.merge(params))
}

func (c *Comparer) merge(params contract.CompareParams) contract.CompareParams {
	mergedParams := c.defaultComparisonParams
	if params.AllowArrayExtraItems != nil {
		mergedParams.AllowArrayExtraItems = params.AllowArrayExtraItems
	}
	if params.DisallowExtraFields != nil {
		mergedParams.DisallowExtraFields = params.DisallowExtraFields
	}
	if params.IgnoreArraysOrdering != nil {
		mergedParams.IgnoreArraysOrdering = params.IgnoreArraysOrdering
	}
	if params.IgnoreValues != nil {
		mergedParams.IgnoreValues = params.IgnoreValues
	}

	return mergedParams
}
