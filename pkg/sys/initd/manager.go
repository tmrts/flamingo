package initd

type Manager interface {
	ReloadDaemon() error

	CreateComponent(name, contents string) (Component, error)

	Start(Component) error
	Stop(Component) error
	Install(Component) error
	Disable(Component) error
	Extend(Component) error
	Reload(Component) error
	Validate(Component) error
}
