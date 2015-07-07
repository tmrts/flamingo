package iptables

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/tmrts/flamingo/pkg/sys"
	"github.com/tmrts/flamingo/pkg/sys/firewall"
	"github.com/tmrts/flamingo/pkg/util"
)

type Target string

const (
	Accept Target = "ACCEPT"
	Drop   Target = "DROP"
	Return Target = "RETURN"
)

type Operation string

const (
	Append Operation = "append"
	Delete Operation = "delete"
	Insert Operation = "insert"
	Check  Operation = "check"
)

type Implementation struct {
	sys.Executor
}

type Rule struct {
	Source      string `flag:"source"`
	Destination string `flag:"source"`

	FromInterface string `flag:"in-interface"`
	ToInterface   string `flag:"out-interface"`

	IsSyncPackage bool `flag:"syn"`

	Protocol string `flag:"protocol"`

	Match string `flag:"match"`

	Target Target `flag:"jump"`
}

func ruleToFlag(r *Rule) string {
	flagForm := util.GetFlagFormOfStruct(*r)

	return strings.Join(flagForm, " ")
}

func (impl *Implementation) CheckDependencies() error {
	_, err := exec.LookPath("iptables")
	if err != nil {
		return fmt.Errorf("iptables: 'iptables' executable could not be found")
	}

	return nil
}

func (impl *Implementation) Perform(op Operation, c firewall.Chain, r *Rule) error {
	action := fmt.Sprintf("--table=%v --%v=%v", c.Table, op, c.Name)

	_, err := impl.Execute("iptables", action, ruleToFlag(r))

	return err
}
