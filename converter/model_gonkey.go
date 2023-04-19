package converter

import (
	"time"

	"github.com/ixpectus/declarate/compare"
)

type GonkeyTest struct {
	Name                     string                    `json:"name" yaml:"name"`
	Description              *string                   `json:"description" yaml:"description"`
	Status                   *string                   `json:"status" yaml:"status"`
	Variables                map[string]string         `json:"variables" yaml:"variables"`
	VariablesToSet           map[int]map[string]string `json:"variables_to_set" yaml:"variables_to_set"`
	DBVariablesToSet         map[string]string         `json:"db_variables_to_set" yaml:"db_variables_to_set"`
	Method                   string                    `json:"method" yaml:"method"`
	RequestURL               string                    `json:"path" yaml:"path"`
	QueryParams              *string                   `json:"query" yaml:"query"`
	DBConnectionString       string                    `json:"db_conn" yaml:"db_conn"`
	RequestTmpl              string                    `json:"request" yaml:"request"`
	ResponseTmpls            map[int]string            `json:"response" yaml:"response"`
	ResponseHeaders          map[int]map[string]string `json:"responseHeaders" yaml:"responseHeaders"`
	BeforeScriptParams       *scriptParams             `json:"beforeScript" yaml:"beforeScript"`
	AfterRequestScriptParams *scriptParams             `json:"afterRequestScript" yaml:"afterRequestScript"`
	HeadersVal               map[string]string         `json:"headers" yaml:"headers"`
	CookiesVal               map[string]string         `json:"cookies" yaml:"cookies"`
	ComparisonParams         compare.CompareParams     `json:"comparisonParams" yaml:"comparisonParams"`
	FixtureFiles             []string                  `json:"fixtures" yaml:"fixtures"`
	Tags                     []string                  `json:"tags" yaml:"tags"`
	MocksDefinition          map[string]interface{}    `json:"mocks" yaml:"mocks"`
	PauseValue               *int                      `json:"pause" yaml:"pause"`
	DbQueryTmpl              string                    `json:"dbQuery" yaml:"dbQuery"`
	DbResponseTmpl           []string                  `json:"dbResponse" yaml:"dbResponse"`
	DatabaseChecks           []DatabaseCheck           `json:"dbChecks" yaml:"dbChecks"`
	PollInterval             []time.Duration           `json:"pollInterval" yaml:"pollInterval"`
	Poll                     *Poll                     `json:"poll" yaml:"poll"`
}

type DatabaseCheck struct {
	DbQueryTmpl      string            `json:"dbQuery" yaml:"dbQuery"`
	DbResponseTmpl   []string          `json:"dbResponse" yaml:"dbResponse"`
	DBVariablesToSet map[string]string `json:"db_variables_to_set" yaml:"db_variables_to_set"`
}

type Poll struct {
	Duration           time.Duration `json:"duration,omitempty" yaml:"duration,omitempty"`
	Interval           time.Duration `json:"interval,omitempty" yaml:"interval,omitempty"`
	ResponseBodyRegexp string        `json:"response_body_regexp,omitempty" yaml:"response_body_regexp,omitempty"`
}

type scriptParams struct {
	PathTmpl string `json:"path" yaml:"path"`
	Timeout  int    `json:"timeout" yaml:"timeout"`
}
