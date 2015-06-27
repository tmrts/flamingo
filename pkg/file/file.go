package file

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/tmrts/flamingo/pkg/util/rand"
)

type arguments struct {
	UID         int
	GID         int
	Flags       int
	Contents    string
	Permissions os.FileMode
}

type argument func(*arguments)

func UID(userID int) argument {
	return func(args *arguments) {
		args.UID = userID
	}
}

func GID(groupID int) argument {
	return func(args *arguments) {
		args.GID = groupID
	}
}

func Contents(c string) argument {
	return func(args *arguments) {
		args.Contents = c
	}
}

func Permissions(perms os.FileMode) argument {
	return func(args *arguments) {
		args.Permissions = perms
	}
}

func New(filepath string, setters ...argument) error {
	// Default arguments
	args := &arguments{
		UID:         os.Getuid(),
		GID:         os.Getgid(),
		Contents:    "",
		Permissions: 0666,
		Flags:       os.O_CREATE | os.O_EXCL | os.O_WRONLY,
	}

	for _, setter := range setters {
		setter(args)
	}

	f, err := os.OpenFile(filepath, args.Flags, args.Permissions)
	if err != nil {
		return err
	} else {
		defer f.Close()
	}

	if _, err := f.WriteString(args.Contents); err != nil {
		return err
	}

	return f.Chown(args.UID, args.GID)
}

func Temp(contents string, perms os.FileMode) (f *os.File, err error) {
	var (
		defaultTempDir = os.TempDir()
		tempFilePrefix = "flamingo-tempfile"
	)

	fname, err := UniqueName(defaultTempDir, tempFilePrefix)
	if err != nil {
		return
	}

	if err = New(fname, Contents(contents), Permissions(perms)); err != nil {
		return
	}

	return os.OpenFile(fname, os.O_RDWR, perms)
}

func WriteTo(filepath, contents string, args ...argument) error {
	args = append(args, Contents(contents))

	return New(filepath, args...)
}

func UniqueName(dir, prefix string) (string, error) {
	for i, retries := 0, 100; i < retries; i++ {
		suffix := rand.String(8)

		fname := filepath.Join(dir, prefix+suffix)

		if _, err := os.Lstat(fname); err != nil {
			if os.IsNotExist(err) {
				return fname, nil
			} else {
				return "", err
			}
		}
	}

	return "", fmt.Errorf("file name collision after every retry")
}

func EnsureExists(filename string, perms os.FileMode, userID int, groupID int) error {
	err := New(filename, Permissions(perms), UID(userID), GID(groupID))
	if !os.IsExist(err) {
		return err
	}

	return nil
}
