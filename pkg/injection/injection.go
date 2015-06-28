package injection

type Restore func()

type Injector interface {
	Inject(interface{}) Restore
}

type Dependency struct {
	implementation interface{}
}

type ContextManager struct {
	Dependency  Injector
	Replacement interface{}

	restore Restore
}

func (icm *ContextManager) Enter() error {
	icm.restore = icm.Dependency.Inject(icm.Replacement)

	return nil
}

func (icm *ContextManager) Exit() error {
	icm.restore()

	return nil
}

func (icm *ContextManager) Use(fn interface{}) error {
	return fn.(func() error)()
}
