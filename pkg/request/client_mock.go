package request

type MockClient struct {
	Client

	Method     chan string
	URL        chan []string
	Parameters chan []Parameter

	OutResponse *Response
	OutError    error
}

func NewMockExecutor(r *Response, err error) *MockExecutor {
	return &MockExecutor{
		Method:     make(chan string, 1),
		URL:        make(chan string, 1),
		Parameters: make(chan []Parameter, 1),

		OutResponse: r,
		OutError:    err,
	}
}

func (mc *MockClient) Perform(method, url string, params ...Parameter) (*Response, error) {
	mc.Method <- method
	mc.URL <- url
	mc.Parameters <- params

	return mc.OutResponse, mc.OutError
}
