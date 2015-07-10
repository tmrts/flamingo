package rest

import "net/http"

var DefaultClient = ClientImplementation{http.DefaultClient}

type Client interface {
	Perform(string, ...Parameter) (*Response, error)
}

type ClientImplementation struct {
	*http.Client
}

func (c *ClientImplementation) Perform(r *Request) (*Response, error) {
	req := r.Normalize()

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	return &Response{resp}, nil
}
