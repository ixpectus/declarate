package suite

type testDefinition struct {
	Definition *struct {
		Tags []string `yaml:"tags,omitempty"`
	} `yaml:"definition,omitempty"`
}

type testWithDefinition struct {
	file       string
	definition testDefinition
}

// func (r *TestRunner) testsWithTags(tests []string) []testWithTag {
// 	testWithTags := make([]testWithTag, 0, len(tests))
// 	for _, v := range tests {
// 		data, err := os.ReadFile(v)
// 		require.NoError(r.t, err)
// 		var testDefinitions []testDefinition
// 		err = yaml.Unmarshal(data, &testDefinitions)
// 		if len(testDefinitions) == 0 {
// 			continue
// 		}
// 		require.NoError(r.t, err)
// 		testWithTags = append(testWithTags, testWithTag{
// 			name: v,
// 			tags: testDefinitions[0].Tags,
// 		})
// 	}

// 	return testWithTags
// }

// func (r *TestRunner) filterTestsByTags(tests []string) []string {
// 	res := make([]string, 0, len(tests))
// 	testWithTags := r.testsWithTags(tests)
// 	testWithTags = tools.Filter(testWithTags, func(test testWithTag) bool {
// 		return !tools.Contains(test.tags, "skip")
// 	})
// 	if len(r.cfg.tags) == 0 {
// 		for _, v := range testWithTags {
// 			res = append(res, v.name)
// 		}
// 		return res
// 	}
// 	for _, v := range r.cfg.tags {
// 		for _, v1 := range testWithTags {
// 			if tools.Contains(v1.tags, v) && !tools.Contains(res, v1.name) {
// 				res = append(res, v1.name)
// 			}
// 		}
// 	}

// 	return res
// }

// type testDefinition struct {
// 	Tags []string `json:"tags" yaml:"tags"`
// }
