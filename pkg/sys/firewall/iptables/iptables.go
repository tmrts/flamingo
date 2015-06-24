package iptables

import (
	"fmt"
	"os/exec"

	"github.com/tmrts/flamingo/pkg/sys"
	"github.com/tmrts/flamingo/pkg/sys/firewall"
)

type Target string

const (
	Accept Target = "ACCEPT"
	Drop   Target = "DROP"
	Return Target = "RETURN"
)

const (
	Append = "--append"
	Delete = "--delete"
	Insert = "--insert"
	Check  = "--check"
)

type Implementation struct {
	sys.Executor
}

func (impl *Implementation) CheckDependencies() error {
	_, err := exec.LookPath("iptables")
	if err != nil {
		return fmt.Errorf("iptables: 'iptables' executable could not be found")
	}

	return nil
}

func (impl *Implementation) AppendTo(c firewall.Chain, rule ...string) error {
	_, err := impl.Execute(Append, c.Name)

	return err
}

func (impl *Implementation) DeleteFrom(c firewall.Chain, rule ...string) error {
	_, err := impl.Execute(Delete, c.Name)

	return err
}

func (impl *Implementation) InsertTo(c firewall.Chain, rule ...string) error {
	_, err := impl.Execute(Insert, c.Name)

	return err
}

func (impl *Implementation) Check(c firewall.Chain, rule ...string) (bool, error) {
	_, err := impl.Execute(Check, c.Name)

	return err != nil, err
}
