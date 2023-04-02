package contract

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
	VariablesToSet() map[string]string
	Check() error
	SetVars(vv Vars)
}

type Result struct {
	File   string
	Name   string
	Passed bool
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
)

type Message struct {
	Name     string
	Filename string
	Message  string
	Title    string
	Expected string
	Actual   string
	Lvl      int
	Type     MessageType
}

type Output interface {
	Log(message Message)
}

type Evaluator interface {
	Evaluate(s string) string
}
