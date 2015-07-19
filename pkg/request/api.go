package request

func Get(url string, params ...Parameter) (*Response, error) {
	return DefaultClient.Perform("GET", url, params...)
}

func Put(url string, params ...Parameter) (*Response, error) {
	return DefaultClient.Perform("PUT", url, params...)
}

func Post(url string, params ...Parameter) (*Response, error) {
	return DefaultClient.Perform("POST", url, params...)
}

func Head(url string, params ...Parameter) (*Response, error) {
	return DefaultClient.Perform("HEAD", url, params...)
}

func Delete(url string, params ...Parameter) (*Response, error) {
	return DefaultClient.Perform("DELETE", url, params...)
}

func Options(url string, params ...Parameter) (*Response, error) {
	return DefaultClient.Perform("OPTIONS", url, params...)
}
