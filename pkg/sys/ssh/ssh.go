package ssh

import (
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/tmrts/flamingo/pkg/context"
	"github.com/tmrts/flamingo/pkg/file"
	"github.com/tmrts/flamingo/pkg/sys"
)

var (
	// DirPath is the path of the default SSH directory.
	DirPath = ".ssh"

	// AuthorizedKeysPath is the path of the default authorized_keys file for SSH.
	AuthorizedKeysPath = ".ssh/authorized_keys"
)

// Key is an SSH key public or private.
type Key []byte

// KeyPair is an SSH keypair, public and private.
type KeyPair struct {
	Public  Key
	Private Key
}

// Verify uses 'ssh-keygen' utility to verify an SSH key.
// It returns an error if a problem occurs or the key is invalid.
// The caller should diagnose the error for more information.
func Verify(key Key) error {
	tmpFile := &context.TempFile{
		Content: string(key),
	}

	var checkSSHValidity = func(f *os.File) error {
		_, err := sys.DefaultExecutor.Execute("ssh-keygen", "-l", "-f", f.Name())
		return err
	}

	errch := context.Using(tmpFile, checkSSHValidity)

	return <-errch
}

// InitializeFor checks whether .ssh directory exists in the given user's home
// directory and creates the directory if it doesn't exist.
func InitializeFor(usr *user.User) error {
	userSSHDirPath := filepath.Join(usr.HomeDir, DirPath)
	userAuthorizedKeysPath := filepath.Join(usr.HomeDir, AuthorizedKeysPath)

	userID, err := strconv.Atoi(usr.Uid)
	if err != nil {
		return err
	}

	groupID, err := strconv.Atoi(usr.Gid)
	if err != nil {
		return err
	}

	err = file.EnsureDirectoryExists(userSSHDirPath, 0700, userID, groupID)
	if err != nil {
		return err
	}

	return file.EnsureExists(userAuthorizedKeysPath, 0600, userID, groupID)
}

// AuthorizeKeys appends the given SSH public keys to the given os.File.
func AuthorizeKeys(authorizedKeysFile *os.File, publicKeys ...Key) error {
	var keys []string
	for _, key := range publicKeys {
		keys = append(keys, string(key))
	}

	_, err := authorizedKeysFile.WriteString(strings.Join(keys, "\n") + "\n")
	return err
}

// AuthorizeKeysFor initializes the SSH directory structure for the user and
// appends the given public keys to user's authorized keys file.
func AuthorizeKeysFor(usr *user.User, publicKeys []Key) error {
	err := InitializeFor(usr)
	if err != nil {
		return err
	}

	authorizationFile := &context.File{
		Path:        filepath.Join(usr.HomeDir, AuthorizedKeysPath),
		Permissions: os.FileMode(0600),
	}

	// TODO(tmrts): Check for duplicate keys.
	// TODO(tmrts): Add keys to /etc/ssh if keys exist.
	return <-context.Using(authorizationFile, func(f *os.File) error {
		return AuthorizeKeys(f, publicKeys...)
	})
}
