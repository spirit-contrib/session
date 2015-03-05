package client

import (
	"time"

	"github.com/gogap/spirit"
	"github.com/spirit-contrib/inlet_http"
	apiClient "github.com/spirit-contrib/inlet_http_api/client"
)

type HTTPAPISessionClient struct {
	client  apiClient.APIClient
	apiName string
}

func NewHTTPAPISessionClient(sessionServerUrl string, apiHeader string, apiName string, timeout time.Duration) SessionClient {
	cli := apiClient.NewHTTPAPIClient(sessionServerUrl, apiHeader, timeout)
	sessionCli := HTTPAPISessionClient{
		client:  cli,
		apiName: apiName,
	}
	return &sessionCli
}

func (p *HTTPAPISessionClient) Get(cookies map[string]string) (values map[string]interface{}, err error) {
	payload := spirit.Payload{}
	payload.SetContext(inlet_http.CTX_HTTP_COOKIES, cookies)

	cookieKeys := []string{}

	for k, _ := range cookies {
		cookieKeys = append(cookieKeys, k)
	}

	req := map[string]interface{}{"cookie_names": cookieKeys}
	payload.SetContent(req)

	v := make(map[string]interface{})

	err = p.client.Call(p.apiName, payload, &v)

	values = v

	return
}
