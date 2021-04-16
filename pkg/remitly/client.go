package remitly

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Clienter interface {
}

type clientContext struct {
	baseURL  *url.URL
	username string
	hc       http.Client
}

func NewClient(cloudHost *url.URL, username string) Clienter {
	c := clientContext{
		baseURL:  cloudHost,
		username: username,
		hc: http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:    10,
				IdleConnTimeout: 30 * time.Second,
			},
		},
	}
	return &c
}

func (c *clientContext) do(ctx context.Context, method, path string, body interface{}, args ...interface{}) (*http.Response, error) {
	url := fmt.Sprintf(c.baseURL.String()+path, args)

	var (
		req *http.Request
		err error
	)
	if strings.ToUpper(method) == http.MethodGet || body == nil {
		req, err = http.NewRequest(method, url, nil)
	} else {
		b, err := json.Marshal(&body)
		if err != nil {
			return nil, err
		}
		req, err = http.NewRequest(method, url, bytes.NewBuffer(b))
	}
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", c.username)
	req.WithContext(ctx)

	return c.hc.Do(req)
}
