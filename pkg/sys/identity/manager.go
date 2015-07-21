package identity

import (
	"fmt"

	"github.com/tmrts/flamingo/pkg/sys"
	"github.com/tmrts/flamingo/pkg/util"
)

type User struct {
	Name            string `yaml:"name"`
	PasswordHash    string `yaml:"passwd" flag:"password"`
	GECOS           string `yaml:"gecos" flag:"comment"`
	HomeDir         string `yaml:"homedir" flag:"home"`
	DefaultShell    string `yaml:"shell" flag:"shell"`
	IsSystemAccount bool   `yaml:"system" flag:"system"`

	UID             string   `flag:"uid"`
	GID             string   `yaml:"primary-group" flag:"gid"`
	SecondaryGroups []string `yaml:"groups" flag:"groups"`

	// TODO(tmrts): add --default flag
	IsInactive                  bool   `yaml:"inactive" flag:"no-create-home"`
	CreateHomeDir               bool   `yaml:"no-create-home" flag:"no-create-home"`
	CreateUserGroupWithSameName bool   `flag:"user-group"`
	ExpireDate                  string `yaml:"expiredate" flag:"expiredate"`
	SELinuxUser                 string `yaml:"selinux-user" flag:"selinux-user"`
	DirectoryTemplate           string `flag:"skel"`
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

	AddUserToGroup(string, string) error
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

func (mi *ManagerImplementation) AddUserToGroup(userName, groupName string) error {
	_, err := mi.Executor.Execute("usermod", userName, "--append", "--groups="+groupName)

	return err
}
