package distro

import (
	"github.com/tmrts/flamingo/pkg/sys"
	"github.com/tmrts/flamingo/pkg/sys/firewall/iptables"
	"github.com/tmrts/flamingo/pkg/sys/identity"
	"github.com/tmrts/flamingo/pkg/sys/initd/systemd"
	"github.com/tmrts/flamingo/pkg/sys/nss"
)

// CentOS returns the distribution implementation of CentOS
// operating system that uses the given sys.Executor.
func CentOS(exec sys.Executor) ContextConsumer {
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
