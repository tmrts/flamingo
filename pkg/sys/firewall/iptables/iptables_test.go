package iptables_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/tmrts/flamingo/pkg/sys"
	"github.com/tmrts/flamingo/pkg/sys/firewall"
	"github.com/tmrts/flamingo/pkg/sys/firewall/iptables"
)

func TestAddingFirewallRules(t *testing.T) {
	Convey("Given a firewall manager implementation", t, func() {
		exec := sys.NewStubExecutor("", nil)
		fwllmgr := iptables.Implementation{exec}

		Convey("With a package filtering rule and a target chain", func() {
			chain := firewall.Chain{
				Name:  "INPUT",
				Table: firewall.Filter,
			}

			rule := &iptables.Rule{
				Source: "192.168.1.0/34",
			}

			Convey("It should append the rule to the chain", func() {
				err := fwllmgr.Perform(iptables.Append, chain, rule)

				So(<-exec.Exec, ShouldEqual, "iptables")
				So(<-exec.Args, ShouldEqual, "")
				So(err, ShouldBeNil)
			})

			Convey("It should delete the rule to the chain", func() {
				err := fwllmgr.Perform(iptables.Delete, chain, rule)

				So(<-exec.Exec, ShouldEqual, "iptables")
				So(<-exec.Args, ShouldEqual, "")
				So(err, ShouldBeNil)
			})

			Convey("It should insert the rule to the chain", func() {
				err := fwllmgr.Perform(iptables.Insert, chain, rule)

				So(<-exec.Exec, ShouldEqual, "iptables")
				So(<-exec.Args, ShouldEqual, "")
				So(err, ShouldBeNil)
			})

			Convey("It should check whether the rule exists in the chain", func() {
				err := fwllmgr.Perform(iptables.Check, chain, rule)

				So(<-exec.Exec, ShouldEqual, "iptables")
				So(<-exec.Args, ShouldEqual, "")
				So(err, ShouldBeNil)
			})
		})
	})
}
