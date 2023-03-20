package contract

type Vars interface {
	Set(k, val string)
	Get(k string) string
	Apply(text string) string
}

type CommandBuilder func(unmarshal func(interface{}) error) (Doer, error)

type Doer interface {
	Do() error
	ResponseBody() *string
	Check() error
	SetVars(vv Vars)
}
