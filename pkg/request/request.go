package request

import "net/http"

type Request struct {
	Method, URL string

	Headers http.Header
}

func (r *Request) Normalize() *http.Request {
	req, err := http.NewRequest(r.Method, r.URL, nil)
	if err != nil {
		panic(err)
	}

	req.Header = r.Headers

	return req
}

type Parameter func(*Request)

func Header(key string, values ...string) Parameter {
	return func(r *Request) {
		r.Headers.Del(key)

		for _, v := range values {
			r.Headers.Add(key, v)
		}
	}
}
