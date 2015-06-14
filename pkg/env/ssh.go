package env

import (
	"os"

	"github.com/TamerTas/cloud-init/pkg/utils"
)

// VerifySSHKey uses ssh-keygen utility to verify an SSH key.
// It returns an error if a problem occurs or the key is invalid.
// The caller should diagnose the error for more information.
func VerifySSHKey(key []byte) error {
	tmpFile, err := utils.CreateTempFile(string(key))
	if err != nil {
		return err
	}

	tmpFileName := tmpFile.Name()
	defer os.Remove(tmpFileName)

	_, err := ExecuteCommand("ssh-keygen", "-l", "-f", tmpFileName)

	return err
}
