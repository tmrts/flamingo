package rest

import "net/http"

func Get(url string, params ...Parameter) (*Response, error) {
	r := &Request{
		URL:     url,
		Method:  "GET",
		Headers: http.Header{},
	}

	for _, parametrize := range params {
		parametrize(r)
	}

	return DefaultClient.Perform(r)
}

//func Put(string, ...Parameter) (*Response, error)     {}
//func Post(string, ...Parameter) (*Response, error)    {}
//func Head(string, ...Parameter) (*Response, error)    {}
//func Delete(string, ...Parameter) (*Response, error)  {}
//func Options(string, ...Parameter) (*Response, error) {}
