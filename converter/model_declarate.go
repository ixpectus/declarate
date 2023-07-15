package converter

import "github.com/ixpectus/declarate/compare"

type DeclarateTest struct {
	Name             string                `yaml:"name,omitempty"`
	DbConn           string                `yaml:"db_conn,omitempty"`
	DbQuery          string                `yaml:"db_query,omitempty"`
	DbResponse       string                `yaml:"db_response,omitempty"`
	ComparisonParams compare.CompareParams `yaml:"comparisonParams,omitempty"`
	ScriptPath       string                `yaml:"script_path,omitempty,omitempty"`
	ScriptResponse   *string               `yaml:"script_response,omitempty,omitempty"`
	RequestTmpl      string                `yaml:"request,omitempty,flow"`
	RequestURL       string                `yaml:"path,omitempty" yaml:"path"`
	Method           string                `yaml:"method,omitempty"`
	ResponseTmpls    string                `yaml:"response,omitempty"`
	Steps            []DeclarateTest       `yaml:"steps,omitempty"`
	HeadersVal       map[string]string     `json:"headers,omitempty" yaml:"headers,omitempty"`
	Variables        map[string]string     `yaml:"variables,omitempty"`
	Poll             *Poll                 `yaml:"poll,omitempty"`
	Definition       *Definition           `yaml:"definition,omitempty"`
}

type Definition struct {
	Tags []string `yaml:"tags,omitempty"`
}

type response struct {
	Body   string `json:"body,omitempty"`
	Status int    `json:"status,omitempty"`
}
