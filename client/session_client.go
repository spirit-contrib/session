package client

type SessionClient interface {
	Get(cookies map[string]string) (values map[string]interface{}, err error)
}
