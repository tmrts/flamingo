package nss

import (
	"fmt"
	"os/user"
	"strconv"

	"github.com/tmrts/flamingo/pkg/sys"
)

type Database string

const (
	UserDatabase        Database = "passwd"
	UserShadowDatabase  Database = "shadow"
	GroupDatabase       Database = "group"
	GroupShadowDatabase Database = "gshadow"
)

type Service interface {
	GetEntryFrom(Database, string) (string, error)
}

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

func New(e sys.Executor) *Server {
	return &Server{
		Exec: e,
	}
}

func GetIDsForUser(username string) (userID int, groupID int, err error) {
	usr, err := user.Lookup(username)
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
