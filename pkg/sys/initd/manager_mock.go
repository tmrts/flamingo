package initd

type MockManager struct {
	Component chan Component
}

func NewMockManager() *MockManager {
	return &MockManager{
		Component: make(chan Component),
	}
}

func (mm *MockManager) ReloadDaemon() error {
	return nil
}

func (mm *MockManager) CreateComponent(name, contents string) error {
	return nil
}

func (mm *MockManager) Start(c Component) error {
	mm.Component <- c
	return nil
}

func (mm *MockManager) Stop(c Component) error {
	mm.Component <- c
	return nil
}

func (mm *MockManager) Disable(c Component) error {
	mm.Component <- c
	return nil
}

func (mm *MockManager) Extend(c Component) error {
	mm.Component <- c
	return nil
}

func (mm *MockManager) Install(c Component) error {
	mm.Component <- c
	return nil
}

func (mm *MockManager) Reload(c Component) error {
	mm.Component <- c
	return nil
}

func (mm *MockManager) Validate(c Component) error {
	mm.Component <- c
	return nil
}
