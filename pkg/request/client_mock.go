package request

type MockClient struct {
	Client

	Method     chan string
	URL        chan string
	Parameters chan []Parameter

	OutResponse chan *Response
	OutError    error
}

func NewMockClient(l int) *MockClient {
	return &MockClient{
		Method:     make(chan string, 1),
		URL:        make(chan string, 1),
		Parameters: make(chan []Parameter, 1),

		OutResponse: make(chan *Response, l),
		OutError:    nil,
	}
}

func (mc *MockClient) Perform(method, url string, params ...Parameter) (*Response, error) {
	mc.Method <- method
	mc.URL <- url
	mc.Parameters <- params

	return <-mc.OutResponse, mc.OutError
}
