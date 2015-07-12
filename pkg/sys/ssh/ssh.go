package ssh

import (
	"os"
	"os/user"
	"strconv"
	"strings"

	"github.com/tmrts/flamingo/pkg/context"
	"github.com/tmrts/flamingo/pkg/file"
	"github.com/tmrts/flamingo/pkg/sys"
)

const (
	SSHDirPath         = ".ssh"
	AuthorizedKeysPath = ".ssh/authorized_keys"
)

type Key []byte

type KeyPair struct {
	Public, Private Key
}

// Verify uses ssh-keygen utility to verify an SSH key.
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

func InitializeFor(owner *user.User) error {
	userSSHDirPath := owner.HomeDir + "/" + SSHDirPath
	userAuthorizedKeysPath := owner.HomeDir + "/" + AuthorizedKeysPath

	userID, err := strconv.Atoi(owner.Uid)
	if err != nil {
		return err
	}

	groupID, err := strconv.Atoi(owner.Gid)
	if err != nil {
		return err
	}

	err = file.EnsureDirectoryExists(userSSHDirPath, 0700, userID, groupID)
	if err != nil {
		return err
	}

	return file.EnsureExists(userAuthorizedKeysPath, 0600, userID, groupID)
}

func AuthorizeKeys(authorizedKeysFile *os.File, publicKeys ...Key) error {
	var keys []string
	for _, key := range publicKeys {
		keys = append(keys, string(key))
	}

	_, err := authorizedKeysFile.WriteString(strings.Join(keys, "\n") + "\n")
	return err
}
