package request

import "net/http"

var DefaultClient = ClientImplementation{http.DefaultClient}

type Client interface {
	Perform(string, ...Parameter) (*Response, error)
}

type ClientImplementation struct {
	*http.Client
}

func (c *ClientImplementation) perform(r *Request) (*Response, error) {
	req := r.Normalize()

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	return &Response{resp}, nil
}

func Perform(method, url string, params ...Parameter) (*Response, error) {
	r := &Request{
		URL:     url,
		Method:  method,
		Headers: http.Header{},
	}

	for _, parametrize := range params {
		parametrize(r)
	}

	return DefaultClient.perform(r)
}
