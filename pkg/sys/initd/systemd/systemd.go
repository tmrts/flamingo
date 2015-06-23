package systemd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/tmrts/flamingo/pkg/context"
	"github.com/tmrts/flamingo/pkg/sys"
	"github.com/tmrts/flamingo/pkg/sys/initd"
)

const (
	DefaultUnitDir = "/etc/systemd/system"
)

// TODO(tmrts): Use systemd DBUS unix socket to interact with systemd,
//				instead of using systemctl commands.
type Implementation struct {
	UnitDir string
	Exec    sys.Executor
}

func (sysd *Implementation) ReloadDaemon() error {
	_, err := sysd.Exec.Execute("systemctl", "daemon-reload")

	return err
}

func (sysd *Implementation) CreateComponent(name, contents string) (initd.Component, error) {
	newUnit := &Unit{
		name: name,
		path: filepath.Join(sysd.UnitDir, name),
	}

	newFile := &context.NewFile{
		Path:        newUnit.Path(),
		Permissions: 0644,
	}

	errch := context.Using(newFile, func(f *os.File) error {
		_, err := f.WriteString(contents)

		return err
	})

	return newUnit, <-errch
}

func (sysd *Implementation) Start(c initd.Component) error {
	_, err := sysd.Exec.Execute("systemctl", "start", c.Name())

	return err
}

func (sysd *Implementation) Stop(c initd.Component) error {
	_, err := sysd.Exec.Execute("systemctl", "stop", c.Name())

	return err
}

func (sysd *Implementation) Install(c initd.Component) error {
	_, err := sysd.Exec.Execute("systemctl", "enable", "--system", c.Name())

	return err
}

func (sysd *Implementation) Disable(c initd.Component) error {
	_, err := sysd.Exec.Execute("systemctl", "disable", c.Name())

	return err
}

func (sysd *Implementation) Extend(c initd.Component) error {
	// unit.d file must be present
	_, err := sysd.Exec.Execute("systemctl", "edit", c.Name())

	return err
}

func (sysd *Implementation) Reload(c initd.Component) error {
	_, err := sysd.Exec.Execute("systemctl", "reenable", c.Name())

	return err
}

func (sysd *Implementation) Validate(c initd.Component) error {
	if !strings.HasPrefix(c.Path(), sysd.UnitDir) {
		return fmt.Errorf("systemd: unit file should be in %v", sysd.UnitDir)
	}

	return nil
}
