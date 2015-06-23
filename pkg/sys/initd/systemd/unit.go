package systemd

import "github.com/tmrts/flamingo/pkg/sys/initd"

type Unit struct {
	name string
	path string
}

func NewUnit(name, path string) initd.Component {
	return Unit{name, path}
}

func (u Unit) Name() string {
	return u.name
}

func (u Unit) Path() string {
	return u.path
}
