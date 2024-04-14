package contract

import (
	"database/sql"
	"time"
)

type Vars interface {
	Set(k, val string) error
	Get(k string) string
	Apply(text string) string
	SetPersistent(k, val string) error
}

type CommandBuilder interface {
	Build(unmarshal func(interface{}) error) (Doer, error)
}

type Doer interface {
	Do() error
	ResponseBody() *string
	IsValid() error
	GetConfig() interface{}
	Check() error
	SetVars(vv Vars)
	SetReport(r ReportAttachement)
}

type TestError struct {
	Title         string
	Expected      string
	Actual        string
	Message       string
	OriginalError error
}

func (e *TestError) Error() string {
	return e.Message
}

func (e *TestError) Unwrap() error {
	return e.OriginalError
}

type (
	MessageType string
	ActionType  string
)

var (
	MessageTypeSuccess MessageType = "success"
	MessageTypeError   MessageType = "error"
	MessageTypeNotify  MessageType = "notify"
	MessageTypePoll    MessageType = "poll"
)

type PollInfo struct {
	Start     time.Time
	Finish    time.Time
	Estimated time.Duration
}

type PollResult struct {
	Start         time.Time
	Finish        time.Time
	PlannedFinish time.Time
}

type Message struct {
	Name                string
	Filename            string
	Message             string
	ActionType          string
	HasNestedSteps      bool
	HasPoll             bool
	Title               string
	Expected            string
	Actual              string
	Lvl                 int
	Type                MessageType
	Poll                *PollInfo
	PollResult          *PollResult
	PollConditionFailed bool
}

type Output interface {
	Log(message Message)
	SetReport(r Report)
}

type Evaluator interface {
	Evaluate(s string) string
}

type TestWrapper interface {
	BeforeTest(file string, conf *RunConfig, lvl int)
	AfterTest(conf *RunConfig, result Result)
	BeforeTestStep(file string, conf *RunConfig, lvl int)
	AfterTestStep(conf *RunConfig, result Result, isPolling bool)
}

type RunConfig struct {
	Name      string      `yaml:"name,omitempty"`
	Steps     []RunConfig `yaml:"steps,omitempty"`
	Vars      Vars
	Variables map[string]string `yaml:"variables"`
	Commands  []Doer
}

type Result struct {
	Err        error
	Name       string
	Lvl        int
	FileName   string
	Response   *string
	PollResult *PollResult
}

type Comparer interface {
	Compare(expected, actual interface{}, params CompareParams) []error
	CompareJsonBody(
		expectedBody string,
		realBody string,
		params CompareParams,
	) ([]error, error)
}

type DBConnectLoader interface {
	Get(string) (*sql.DB, error)
}

type CompareParams struct {
	IgnoreValues         *bool `json:"ignoreValues,omitempty" yaml:"ignoreValues,omitempty"`
	IgnoreArraysOrdering *bool `json:"ignoreArraysOrdering,omitempty" yaml:"ignoreArraysOrdering,omitempty"`
	DisallowExtraFields  *bool `json:"disallowExtraFields,omitempty" yaml:"disallowExtraFields,omitempty"`
	AllowArrayExtraItems *bool `json:"allowArrayExtraItems,omitempty" yaml:"allowArrayExtraItems,omitempty"`
	LineByLine           *bool `json:"lineByLine,omitempty" yaml:"lineByLine,omitempty"`
	FailFast             bool  `json:"failFast,omitempty" yaml:"failFast,omitempty"` // End compare operation after first error
}

type Persistent interface {
	Set(key string, value string) error
	Get(key string) (string, error)
	Reset() error
}
