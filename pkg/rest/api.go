package rest

func Get(url string, params ...Parameter) (*Response, error) {
	return request("GET", url, params)
}

func Put(string, ...Parameter) (*Response, error) {
	return request("PUT", url, params)
}

func Post(string, ...Parameter) (*Response, error) {
	return request("POST", url, params)
}

func Head(string, ...Parameter) (*Response, error) {
	return request("HEAD", url, params)
}

func Delete(string, ...Parameter) (*Response, error) {
	return request("DELETE", url, params)
}

func Options(string, ...Parameter) (*Response, error) {
	return request("OPTIONS", url, params)
}
