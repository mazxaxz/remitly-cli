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
	// CreateLoadBalancer creates load balancer with given name
	CreateLoadBalancer(ctx context.Context, name string) (LoadBalancer, error)
	// GetInstances returns array of instances by load balancer name
	GetInstances(ctx context.Context, lbName string) ([]Instance, error)
	// CreateInstance within load balancer scope
	CreateInstance(ctx context.Context, lbName, version string) (Instance, error)
	// DeleteInstance deletes instance by ID within load balancer scope
	DeleteInstance(ctx context.Context, lbName, ID string) error
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

func (c *clientContext) CreateLoadBalancer(ctx context.Context, name string) (LoadBalancer, error) {
	return LoadBalancer{}, nil
}

func (c *clientContext) GetInstances(ctx context.Context, lbName string) ([]Instance, error) {
	return nil, nil
}

func (c *clientContext) CreateInstance(ctx context.Context, lbName, version string) (Instance, error) {
	return Instance{}, nil
}

func (c *clientContext) DeleteInstance(ctx context.Context, lbName, ID string) error {
	return nil
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
