package initd

type mockManager struct {
	Component chan Component
}

func NewMockManager() Manager {
	return &mockManager{
		Component: make(chan Component),
	}
}

func (mm *mockManager) ReloadDaemon() error {
	return nil
}

func (mm *mockManager) CreateComponent(name, contents string) (Component, error) {
	return nil, nil
}

func (mm *mockManager) Start(c Component) error {
	mm.Component <- c
	return nil
}

func (mm *mockManager) Stop(c Component) error {
	mm.Component <- c
	return nil
}

func (mm *mockManager) Disable(c Component) error {
	mm.Component <- c
	return nil
}

func (mm *mockManager) Extend(c Component) error {
	mm.Component <- c
	return nil
}

func (mm *mockManager) Install(c Component) error {
	mm.Component <- c
	return nil
}

func (mm *mockManager) Reload(c Component) error {
	mm.Component <- c
	return nil
}

func (mm *mockManager) Validate(c Component) error {
	mm.Component <- c
	return nil
}
