package sys

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
)

// Executor is an interface representing the ability to execute given
// a command and its arguments, and returning the output or the error.
type Executor interface {
	Execute(string, ...string) (string, error)
}

type linux struct{}

// Execute wraps the command execution pattern required in os/exec package.
// The command is executed with the supplied arguments and the output is returned
// to the caller. If an error occurs during the command execution, an error is
// returned with the messages read from the stderr of the executed command.
func (l *linux) Execute(executable string, args ...string) (string, error) {
	cmd := exec.Command(executable, args...)

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

// DefaultExecutor is the default Executor and is used by Execute.
var DefaultExecutor = &linux{}

// Execute executes the given command along with its arguments and the output
// is returned to the caller. If an error occurs during the command execution,
// an error is returned with the messages read from the stderr of the executed command.
func Execute(cmd string, args ...string) (string, error) {
	return DefaultExecutor.Execute(cmd, args...)
}
