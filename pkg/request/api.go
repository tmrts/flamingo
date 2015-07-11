package request

func Get(url string, params ...Parameter) (*Response, error) {
	return Perform("GET", url, params...)
}

func Put(url string, params ...Parameter) (*Response, error) {
	return Perform("PUT", url, params...)
}

func Post(url string, params ...Parameter) (*Response, error) {
	return Perform("POST", url, params...)
}

func Head(url string, params ...Parameter) (*Response, error) {
	return Perform("HEAD", url, params...)
}

func Delete(url string, params ...Parameter) (*Response, error) {
	return Perform("DELETE", url, params...)
}

func Options(url string, params ...Parameter) (*Response, error) {
	return Perform("OPTIONS", url, params...)
}
