package nss

type MockService struct {
	DB  chan Database
	Key chan string

	OutStr chan string
	OutErr chan error
}

func NewMockService() *MockService {
	return &MockService{
		DB:  make(chan Database, 1),
		Key: make(chan string, 1),

		OutStr: make(chan string, 1),
		OutErr: make(chan error, 1),
	}
}

func (ms *MockService) GetEntryFrom(db Database, key string) (string, error) {
	ms.DB <- db
	ms.Key <- key

	return <-ms.OutStr, <-ms.OutErr
}
