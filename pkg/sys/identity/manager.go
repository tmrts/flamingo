package identity

import (
	"fmt"

	"github.com/tmrts/flamingo/pkg/sys"
	"github.com/tmrts/flamingo/pkg/util"
)

// User contains the directives to be used when creating a new user.
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

	// TODO(tmrts): add --default user flag?
	IsInactive                  bool   `yaml:"inactive" flag:"no-create-home"`
	CreateHomeDir               bool   `yaml:"no-create-home" flag:"no-create-home"`
	CreateUserGroupWithSameName bool   `flag:"user-group"`
	ExpireDate                  string `yaml:"expiredate" flag:"expiredate"`
	SELinuxUser                 string `yaml:"selinux-user" flag:"selinux-user"`
	DirectoryTemplate           string `flag:"skel"`
}

// Group contains the directives to be used when creating a new user group.
type Group struct {
	Name            string
	PasswordHash    string `flag:"password"`
	GID             string `flag:"gid"`
	IsSystemAccount bool   `flag:"system"`
}

// Group is an interface representing the ability to
// manipulate system users and groups.
type Manager interface {
	// CreateNewUser persists the supplied User to the system.
	CreateUser(User) error
	// CreateNewGroup persists the supplied Group to the system.
	CreateGroup(Group) error

	// SetUserPassword changes the given user's password to supplied password hash.
	SetUserPassword(string, string) error
	// SetGroupPassword changes the given group's password to supplied password hash.
	SetGroupPassword(string, string) error

	// AddUserToGroup adds the user with the given name to the group
	// with the given name's members.
	AddUserToGroup(string, string) error
}

type managerImplementation struct {
	Executor sys.Executor
}

func (mi *managerImplementation) CreateUser(usr User) error {
	args := util.GetArgumentFormOfStruct(usr)

	args = append(args, usr.Name)

	_, err := mi.Executor.Execute("useradd", args...)

	return err
}

func (mi *managerImplementation) SetUserPassword(userName, passwordHash string) error {
	passwordPair := fmt.Sprintf("%s:%s", userName, passwordHash)

	// TODO(tmrts): pass the password hash using a stdin pipe.
	_, err := mi.Executor.Execute("chpasswd", "-e", passwordPair)

	return err
}

func (mi *managerImplementation) CreateGroup(grp Group) error {
	args := util.GetArgumentFormOfStruct(grp)

	args = append(args, grp.Name)

	_, err := mi.Executor.Execute("groupadd", args...)

	return err
}

func (mi *managerImplementation) SetGroupPassword(groupName, passwordHash string) error {
	// TODO(tmrts): pass the password hash using a stdin pipe.
	_, err := mi.Executor.Execute("groupmod", groupName, "--password="+passwordHash)

	return err
}

func (mi *managerImplementation) AddUserToGroup(userName, groupName string) error {
	_, err := mi.Executor.Execute("usermod", userName, "--append", "--groups="+groupName)

	return err
}

// NewManager returns the Manager interface implementation that uses given
// sys.Executor to execute `shadow-utils` system package commands for
// performing system user/group manipulations.
func NewManager(e sys.Executor) Manager {
	return &managerImplementation{
		Executor: e,
	}
}
