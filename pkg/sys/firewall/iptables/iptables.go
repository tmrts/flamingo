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
	AcceptTarget Target = "ACCEPT"
	ReturnTarget Target = "RETURN"
	DropTarget   Target = "DROP"
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
	Source      []string `flag:"source"`
	Destination []string `flag:"destination"`

	FromInterface string `flag:"in-interface"`
	ToInterface   string `flag:"out-interface"`

	IsSyncPackage bool `flag:"syn"`

	Protocol string `flag:"protocol"`

	Target Target `flag:"jump"`
}

func (r *Rule) FlagForm() []string {
	return util.GetFlagFormOfStruct(*r)
}

func (r *Rule) String() string {
	return strings.Join(r.FlagForm(), " ")
}

func (impl *Implementation) CheckDependencies() error {
	_, err := exec.LookPath("iptables")
	if err != nil {
		return fmt.Errorf("iptables: 'iptables' executable could not be found")
	}

	return nil
}

func (impl *Implementation) Perform(op Operation, c firewall.Chain, r *Rule) error {
	if r.Target == "" {
		return fmt.Errorf("iptables: rule doesnt have a target")
	}

	table := fmt.Sprintf("--table=%v", c.Table)
	action := fmt.Sprintf("--%v=%v", op, c.Name)

	args := []string{"--wait", table, action}

	args = append(args, r.FlagForm()...)

	_, err := impl.Execute("iptables", args...)

	return err
}
