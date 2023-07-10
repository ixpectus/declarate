package contract

import (
	"database/sql"
	"time"

	"github.com/ixpectus/declarate/compare"
)

type Vars interface {
	Set(k, val string)
	Get(k string) string
	Apply(text string) string
}

type CommandBuilder interface {
	Build(unmarshal func(interface{}) error) (Doer, error)
}

type Doer interface {
	Do() error
	ResponseBody() *string
	IsValid() error
	VariablesToSet() map[string]string
	GetConfig() interface{}
	Check() error
	SetVars(vv Vars)
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

type MessageType string

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
	Name       string
	Filename   string
	Message    string
	Title      string
	Expected   string
	Actual     string
	Lvl        int
	Type       MessageType
	Poll       *PollInfo
	PollResult *PollResult
}

type Output interface {
	Log(message Message)
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
	Name           string      `yaml:"name,omitempty"`
	Steps          []RunConfig `yaml:"steps,omitempty"`
	Vars           Vars
	VariablesToSet map[string]string `yaml:"variables_to_set"`
	Commands       []Doer
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
	Compare(expected, actual interface{}, params compare.CompareParams) []error
	CompareJsonBody(
		expectedBody string,
		realBody string,
		params compare.CompareParams,
	) ([]error, error)
}

type DBConnectLoader interface {
	Get(string) (*sql.DB, error)
}
