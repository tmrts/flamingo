package sys

type MockExecutor struct {
	Exec chan string
	Args chan []string

	OutStr chan string
	OutErr chan error
}

func NewMockExecutor() *MockExecutor {
	return &MockExecutor{
		Exec: make(chan string, 1),
		Args: make(chan []string, 1),

		OutStr: make(chan string, 1),
		OutErr: make(chan error, 1),
	}
}

func (me *MockExecutor) Execute(exec string, args ...string) (string, error) {
	me.Exec <- exec
	me.Args <- args

	return <-me.OutStr, <-me.OutErr
}

type StubExecutor struct {
	Exec chan string
	Args chan []string

	out string
	err error
}

func NewStubExecutor(out string, err error) *StubExecutor {
	return &StubExecutor{
		Exec: make(chan string, 1),
		Args: make(chan []string, 1),

		out: out,
		err: err,
	}
}

func (se *StubExecutor) Execute(exec string, args ...string) (string, error) {
	se.Exec <- exec
	se.Args <- args

	return se.out, se.err
}
