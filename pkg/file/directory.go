package file

import "os"

func EnsureDirectoryExists(dirname string, perms os.FileMode, userID int, groupID int) error {
	if err := os.MkdirAll(dirname, perms); err != nil {
		return err
	}

	if err := os.Chmod(dirname, perms); err != nil {
		return err
	}

	return os.Chown(dirname, userID, groupID)
}
