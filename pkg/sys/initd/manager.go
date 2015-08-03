// initd provides interfaces and implementations for system initialization daemons
package initd

// Manager interface wraps initd gateway implementations
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
