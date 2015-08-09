package distro

import (
	"github.com/tmrts/flamingo/pkg/sys"
	"github.com/tmrts/flamingo/pkg/sys/firewall/iptables"
	"github.com/tmrts/flamingo/pkg/sys/identity"
	"github.com/tmrts/flamingo/pkg/sys/initd"
	"github.com/tmrts/flamingo/pkg/sys/initd/systemd"
	"github.com/tmrts/flamingo/pkg/sys/nss"
)

type Implementation struct {
	sys.Executor

	ID    identity.Manager
	NSS   nss.Service
	Initd initd.Manager

	// TODO(tmrts): Change to firewall.Manager interface
	Firewall *iptables.Implementation
}

// SetHostname changes the hostname of the given distribution.
func (imp *Implementation) SetHostname(hostname string) error {
	_, err := imp.Execute("hostnamectl", "set-hostname", hostname)

	return err
}

// SetTimeZone changes the hostname of the given distribution.
func (imp *Implementation) SetTimeZone(timezone string) error {
	_, err := imp.Execute("timedatectl", "set-timezone", timezone)

	return err
}

func CentOS(exec sys.Executor) *Implementation {
	return &Implementation{
		Executor: exec,

		ID:       identity.NewManager(exec),
		NSS:      &nss.Server{exec},
		Firewall: &iptables.Implementation{exec},
		//FileSystem: &file.System{exec},
		//Sysconf: conf.Manager hostname, timedate, etc. sysctl for CentOS
		Initd: &systemd.Implementation{
			UnitDir: "/etc/systemd/userexec",
			Exec:    exec,
		},
	}
}
