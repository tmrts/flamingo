package identity

import (
	"fmt"

	"github.com/tmrts/flamingo/pkg/sys"
	"github.com/tmrts/flamingo/pkg/util"
)

type User struct {
	Name            string
	PasswordHash    string `flag:"password"`
	GECOS           string `flag:"comment"`
	HomeDir         string `flag:"home"`
	DefaultShell    string `flag:"shell"`
	IsSystemAccount bool   `flag:"system"`

	UID             string   `flag:"uid"`
	GID             string   `flag:"gid"`
	SecondaryGroups []string `flag:"groups"`

	CreateHomeDir               bool   `flag:"no-create-home"`
	DirectoryTemplate           string `flag:"skel"`
	CreateUserGroupWithSameName string `flag:"user-group"`
}

type Group struct {
	Name            string
	PasswordHash    string `flag:"password"`
	GID             string `flag:"gid"`
	IsSystemAccount bool   `flag:"system"`
}

type Manager interface {
	CreateUser(User) error
	CreateGroup(Group) error

	SetUserPassword(string, string) error
	SetGroupPassword(string, string) error
}

type ManagerImplementation struct {
	Executor sys.Executor
}

func (mi *ManagerImplementation) CreateUser(usr User) error {
	args := util.GetArgumentFormOfStruct(usr)

	args = append(args, usr.Name)

	_, err := mi.Executor.Execute("useradd", args...)

	return err
}

func (mi *ManagerImplementation) SetUserPassword(userName, passwordHash string) error {
	passwordPair := fmt.Sprintf("%s:%s", userName, passwordHash)

	// TODO(tmrts): pass the password hash using a stdin pipe.
	_, err := mi.Executor.Execute("chpasswd", "-e", passwordPair)

	return err
}

// CreateNewGroup adds the supplied Group to the system groups.
func (mi *ManagerImplementation) CreateGroup(grp Group) error {
	args := util.GetArgumentFormOfStruct(grp)

	args = append(args, grp.Name)

	_, err := mi.Executor.Execute("groupadd", args...)

	return err
}

func (mi *ManagerImplementation) SetGroupPassword(groupName, passwordHash string) error {
	// TODO(tmrts): pass the password hash using a stdin pipe.
	_, err := mi.Executor.Execute("groupmod", groupName, "--password="+passwordHash)

	return err
}
