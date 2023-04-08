package compare

type Comparer struct {
	defaultComparisonParams CompareParams
}

func New(c CompareParams) *Comparer {
	return &Comparer{defaultComparisonParams: c}
}

func (c *Comparer) Compare(expected, actual interface{}, params CompareParams) []error {
	return compare(expected, actual, c.merge(params))
}

func (c *Comparer) CompareJsonBody(
	expectedBody string,
	realBody string,
	params CompareParams,
) ([]error, error) {
	return compareJsonBody(expectedBody, realBody, c.merge(params))
}

func (c *Comparer) merge(params CompareParams) CompareParams {
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
