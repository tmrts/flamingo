package iptables_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	. "github.com/tmrts/flamingo/pkg/util/testutil"

	"github.com/tmrts/flamingo/pkg/sys"
	"github.com/tmrts/flamingo/pkg/sys/firewall"
	"github.com/tmrts/flamingo/pkg/sys/firewall/iptables"
)

func TestAddingFirewallRules(t *testing.T) {
	Convey("Given an action, a target chain and a package filtering rule", t, func() {
		rule := &iptables.Rule{
			Source:        []string{"192.168.1.1", "192.168.1.2"},
			Destination:   []string{"192.168.1.3", "192.168.1.4"},
			FromInterface: "eth0",
			ToInterface:   "eth1",
			Protocol:      "tcp",
			IsSyncPackage: true,
			Target:        iptables.DropTarget,
		}

		action := iptables.Append

		chain := firewall.Chain{
			Name:  "INPUT",
			Table: firewall.Filter,
		}

		Convey("The firewall manager implementation should perform the action", func() {
			exec := sys.NewStubExecutor("", nil)
			fwllmgr := iptables.Implementation{exec}

			err := fwllmgr.Perform(action, chain, rule)
			So(err, ShouldBeNil)

			So(<-exec.Exec, ShouldEqual, "iptables")
			So(<-exec.Args, ShouldBeSuperSetOf, []string{
				"--table=filter",
				"--append=INPUT",
				"--source=192.168.1.1,192.168.1.2",
				"--destination=192.168.1.3,192.168.1.4",
				"--in-interface=eth0",
				"--out-interface=eth1",
				"--protocol=tcp",
				"--jump=DROP",
				"--syn",
			})
		})
	})

	Convey("When performing an operation", t, func() {
		exec := sys.NewStubExecutor("", nil)
		fwllmgr := iptables.Implementation{exec}

		chain := firewall.Chain{
			Name:  "INPUT",
			Table: firewall.Filter,
		}

		rule := &iptables.Rule{
			Protocol: "tcp",
			Target:   iptables.DropTarget,
		}

		Convey("The xtables lock should be acquired to prevent multiple updates at the same time", func() {
			err := fwllmgr.Perform(iptables.Append, chain, rule)
			So(err, ShouldBeNil)

			So(<-exec.Exec, ShouldEqual, "iptables")
			So(<-exec.Args, ShouldContain, "--wait")
		})
	})
}
