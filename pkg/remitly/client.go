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

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type resourceURI string

const (
	putLoadBalancers             = resourceURI("/loadbalancers/%s")
	getLoadBalancersInstances    = resourceURI("/loadbalancers/%s/instances")
	postLoadBalancersInstances   = resourceURI("/loadbalancers/%s/instances")
	deleteLoadBalancersInstances = resourceURI("/loadbalancers/%s/instances/%s")
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
	scheme, hostname string
	username         string
	hc               http.Client
}

func NewClient(cloudHost *url.URL, username string) Clienter {
	c := clientContext{
		scheme:   cloudHost.Scheme,
		hostname: cloudHost.Host,
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
	res, err := c.do(ctx, http.MethodPut, putLoadBalancers, []byte("{}"), name)
	if err != nil {
		return LoadBalancer{}, err
	}
	defer func() { _ = res.Body.Close() }()

	switch res.StatusCode {
	case http.StatusCreated:
		var lb LoadBalancer
		if err := json.NewDecoder(res.Body).Decode(&lb); err != nil {
			return lb, err
		}
		return lb, nil
	case http.StatusForbidden:
		return LoadBalancer{}, ErrForbidden
	case http.StatusNotFound:
		return LoadBalancer{}, ErrNotFound
	default:
		return LoadBalancer{}, errors.Wrapf(ErrUnknown, "http status code: '%d'", res.StatusCode)
	}
}

func (c *clientContext) GetInstances(ctx context.Context, lbName string) ([]Instance, error) {
	res, err := c.do(ctx, http.MethodGet, getLoadBalancersInstances, nil, lbName)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	switch res.StatusCode {
	case http.StatusOK:
		var instances []Instance
		if err := json.NewDecoder(res.Body).Decode(&instances); err != nil {
			return instances, err
		}
		return instances, nil
	case http.StatusForbidden:
		return nil, ErrForbidden
	case http.StatusNotFound:
		return nil, ErrNotFound
	default:
		return nil, errors.Wrapf(ErrUnknown, "http status code: '%d'", res.StatusCode)
	}
}

func (c *clientContext) CreateInstance(ctx context.Context, lbName, version string) (Instance, error) {
	p := CreateInstanceParams{Version: version}
	res, err := c.do(ctx, http.MethodPost, postLoadBalancersInstances, p, lbName)
	if err != nil {
		return Instance{}, err
	}
	defer func() { _ = res.Body.Close() }()

	switch res.StatusCode {
	case http.StatusCreated:
		var instance Instance
		if err := json.NewDecoder(res.Body).Decode(&instance); err != nil {
			return instance, err
		}
		return instance, nil
	case http.StatusForbidden:
		return Instance{}, ErrForbidden
	case http.StatusNotFound:
		return Instance{}, ErrNotFound
	default:
		return Instance{}, errors.Wrapf(ErrUnknown, "http status code: '%d'", res.StatusCode)
	}
}

func (c *clientContext) DeleteInstance(ctx context.Context, lbName, ID string) error {
	res, err := c.do(ctx, http.MethodDelete, deleteLoadBalancersInstances, nil, lbName, ID)
	if err != nil {
		return err
	}
	defer func() { _ = res.Body.Close() }()

	switch res.StatusCode {
	case http.StatusNoContent:
		return nil
	case http.StatusForbidden:
		return ErrForbidden
	case http.StatusNotFound:
		return ErrNotFound
	default:
		return errors.Wrapf(ErrUnknown, "http status code: '%d'", res.StatusCode)
	}
}

func (c *clientContext) do(ctx context.Context, method string, path resourceURI, body interface{}, args ...interface{}) (*http.Response, error) {
	url := url.URL{
		Scheme: c.scheme,
		Host:   c.hostname,
		Path:   fmt.Sprintf(string(path), args...),
	}

	var (
		req *http.Request
		err error
	)
	if strings.ToUpper(method) == http.MethodGet || body == nil {
		req, err = http.NewRequest(method, url.String(), nil)
	} else {
		b, err := json.Marshal(&body)
		if err != nil {
			return nil, err
		}
		req, err = http.NewRequest(method, url.String(), bytes.NewBuffer(b))
		req.Header.Add("Content-Type", "application/json")
	}
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", c.username)
	req.WithContext(ctx)

	now := time.Now()
	res, err := c.hc.Do(req)
	diff := time.Since(now)

	f := log.Fields{
		"milliseconds": diff.Milliseconds(),
		"method":       method,
		"url":          url.String(),
	}
	log.WithContext(ctx).WithFields(f).Trace("http call")
	return res, err
}
