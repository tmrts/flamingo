package env

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
)

func ExecuteCommand(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return "", err
	}

	if err := cmd.Start(); err != nil {
		return "", err
	}

	outBuf, err := ioutil.ReadAll(stdout)
	if err != nil {
		return "", err
	}
	out := string(outBuf)

	errBuf, err := ioutil.ReadAll(stderr)
	if err != nil {
		fmt.Println(errBuf)
		return "", err
	}
	errMsg := string(errBuf)

	if err := cmd.Wait(); err != nil {
		if errMsg != "" {
			return out, errors.New(errMsg)
		}

		return out, err
	}

	return out, nil
}
