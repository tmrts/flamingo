package request

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/tmrts/flamingo/pkg/flog"
)

var DefaultClient = ClientImplementation{http.DefaultClient}

type Client interface {
	Get(string, ...Parameter) (*Response, error)
}

type ClientImplementation struct {
	HTTPClient *http.Client
}

func (c *ClientImplementation) performRequest(r *Request) (*Response, error) {
	req := r.Normalize()

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	if !strings.HasPrefix(resp.Status, "20") {
		flog.Debug("Failed to get a proper HTTP response",
			flog.Fields{
				Error: err,
				Event: "performRequest",
			},
			flog.Details{
				"url":     r.URL,
				"method":  r.Method,
				"headers": r.Headers,
			},
		)
		return nil, fmt.Errorf("request: bad status code %v", resp.StatusCode)
	}

	return &Response{resp}, nil
}

func (c *ClientImplementation) Perform(method, url string, params ...Parameter) (*Response, error) {
	req := &Request{
		URL:     url,
		Method:  method,
		Headers: http.Header{},
	}

	for _, parametrize := range params {
		parametrize(req)
	}

	return c.performRequest(req)
}
