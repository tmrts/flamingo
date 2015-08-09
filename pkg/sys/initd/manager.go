// initd provides interfaces and implementations for system initialization daemons
package initd

// Component represents a file that can be used with
// an initd.Manager (e.g. systemd unit, upstart script, etc.)
type Component interface {
	Name() string
	Path() string
}

// Manager interface represents a system Manager that manipulates
// the system boot events such as systemd, sysV-init, or upstart.
type Manager interface {
	// ReloadDaemon reloads the system Manager daemon.
	ReloadDaemon() error

	// CreateComponent creates a component with the given name
	// and contents and returns an object representing the created component.
	CreateComponent(name, contents string) (Component, error)

	Validate(Component) error
	Install(Component) error
	Disable(Component) error
	Start(Component) error
	Reload(Component) error
	Stop(Component) error
	Extend(Component) error
}
