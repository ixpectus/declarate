package converter

import "github.com/ixpectus/declarate/contract"

type DeclarateTest struct {
	Name             string                 `json:"name,omitempty" yaml:"name,omitempty"`
	DbConn           string                 `json:"db_conn,omitempty" yaml:"db_conn,omitempty"`
	DbQuery          string                 `json:"db_query,omitempty" yaml:"db_query,omitempty"`
	DbResponse       string                 `json:"db_response,omitempty" yaml:"db_response,omitempty"`
	ComparisonParams contract.CompareParams `json:"comparisonParams,omitempty" yaml:"comparisonParams,omitempty"`
	ScriptPath       string                 `json:"script_path,omitempty" yaml:"script_path,omitempty"`
	ScriptResponse   *string                `json:"script_response,omitempty" yaml:"script_response,omitempty"`
	RequestTmpl      string                 `json:"request,omitempty" yaml:"request,omitempty"`
	RequestURL       string                 `json:"path,omitempty" yaml:"path,omitempty" yaml:"path"`
	Method           string                 `json:"method,omitempty" yaml:"method,omitempty"`
	ResponseTmpls    string                 `json:"response,omitempty" yaml:"response,omitempty"`
	Steps            []DeclarateTest        `json:"steps,omitempty" yaml:"steps,omitempty"`
	HeadersVal       map[string]string      `json:"headers,omitempty" yaml:"headers,omitempty"`
	Variables        map[string]string      `json:"variables,omitempty" yaml:"variables,omitempty"`
	Poll             *Poll                  `json:"poll,omitempty" yaml:"poll,omitempty"`
	Definition       *Definition            `json:"definition,omitempty" yaml:"definition,omitempty"`
}

type Definition struct {
	Tags []string `yaml:"tags,omitempty"`
}

type response struct {
	Body   string `json:"body,omitempty"`
	Status int    `json:"status,omitempty"`
}
