package variables

type persistent interface {
	Set(key string, value string) error
	Get(key string) (string, error)
	Reset() error
}
