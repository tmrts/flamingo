package nss

import (
	"fmt"
	"os/user"
	"strconv"
	"strings"

	"github.com/tmrts/flamingo/pkg/sys"
)

// DefaultService is the NSS gateway that interacts with the system
// using the default command execution process in sys.DefaultExecutor.
var DefaultService = New(sys.DefaultExecutor)

// Database represents NSS entry databases such as passwd, group, hosts, etc.
type Database string

const (
	UserDatabase        Database = "passwd"
	UserShadowDatabase  Database = "shadow"
	GroupDatabase       Database = "group"
	GroupShadowDatabase Database = "gshadow"

	HostsDatabase Database = "hosts"
)

// Service is a wrapper that can query the system NSS databases.
type Service interface {
	GetEntryFrom(Database, string) (string, error)
}

// Server is an NSS Service implementation
// It takes a sys.Executor to execute nss query commands.
type Server struct {
	Exec sys.Executor
}

// GetEntryFrom queries the given Name Service Switch libraries with the supplied database key.
// Returns the raw database entry/entries or a meaningful error message
// about the retrieval of the key.
func (s *Server) GetEntryFrom(db Database, key string) (string, error) {
	out, err := s.Exec.Execute("getent", string(db), key)
	if err != nil {
		switch err.Error() {
		case "exit status 1":
			return "", fmt.Errorf("Unknown database => %v", db)
		case "exit status 2":
			return "", fmt.Errorf("Key could not be found => %v", key)
		default:
			return "", err
		}
	}

	return out, nil
}

// New returns an NSS Service using the given sys.Executor
func New(e sys.Executor) *Server {
	return &Server{
		Exec: e,
	}
}

// GetIDsForUser searches the system for the given user name and returns
// the user ID and the primary group ID of that user.
func GetIDsForUser(username string) (userID int, groupID int, err error) {
	usr, err := GetUser(username)
	if err != nil {
		return
	}

	userID, err = strconv.Atoi(usr.Uid)
	if err != nil {
		return
	}

	groupID, err = strconv.Atoi(usr.Gid)
	if err != nil {
		return
	}

	return
}

// CurrentUser is an implementation of user.Current() since it's not
// implemented in some Linux distributions.
func GetCurrentUser() (*user.User, error) {
	usrname, err := sys.DefaultExecutor.Execute("whoami")
	if err != nil {
		return nil, err
	}

	usrname = strings.TrimRight(usrname, "\n")

	return GetUser(usrname)
}

// GetUser is an implementation of user.Lookup() since it's not
// implemented in some Linux distributions.
func GetUser(name string) (*user.User, error) {
	ent, err := GetPasswdEntry(DefaultService, name)
	if err != nil {
		return nil, err
	}

	return &user.User{
		Uid:      strconv.Itoa(ent.UID),
		Gid:      strconv.Itoa(ent.GID),
		Username: ent.UserName,
		HomeDir:  ent.HomeDir,
	}, nil
}
