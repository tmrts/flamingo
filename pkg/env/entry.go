package env

import (
	"errors"
	"fmt"

	"github.com/TamerTas/cloud-init/pkg/cmd"
)

const (
	UserDatabase        = "passwd"
	UserShadowDatabase  = "shadow"
	GroupDatabase       = "group"
	GroupShadowDatabase = "gshadow"
)

// Queries the given Name Service Switch libraries with the supplied database key.
// Returns the raw database entry/entries or a meaningful error message
// about the retrieval of the key.
func GetEntryFrom(database, key string) (string, error) {
	out, err := cmd.ExecuteCommand("getent", database, key)
	if err != nil {
		switch err.Error() {
		case "exit status 1":
			return "", errors.New(fmt.Sprintf("Unknown database => %v", database))
		case "exit status 2":
			return "", errors.New(fmt.Sprintf("Key could not be found => %v", key))
		default:
			return "", err
		}
	}

	return out, nil
}
