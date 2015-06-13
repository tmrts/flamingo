package utils

import (
	"io/ioutil"
	"os"
)

func CreateTempFile(content string) (*os.File, error) {
	const (
		defaultTempDir      = ""
		temporaryFilePrefix = "testutils"
	)

	tmpFile, err := ioutil.TempFile(defaultTempDir, temporaryFilePrefix)
	if err != nil {
		return nil, err
	}

	_, err = tmpFile.WriteString(content)
	if err != nil {
		return nil, err
	}

	return tmpFile, nil
}
