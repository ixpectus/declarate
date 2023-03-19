package contract

type Vars interface {
	Set(k, val string)
	Get(k string) string
	Apply(text string) string
}

type Doer interface {
	Do() error
	ResponseBody() *string
	Check() error
	SetVars(vv Vars)
	FillData(unmarshal func(interface{}) error) error
}
