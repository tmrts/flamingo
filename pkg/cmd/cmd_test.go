package cmd

import "testing"

func TestExecuteCommand(t *testing.T) {
	out, err := ExecuteCommand("echo", "test")
	if err != nil {
		t.Fatalf("error executing echo command -> %v", err)
	}

	if out != nil {
		t.Errorf("wrong output for echo command -> expected: %v, got: %v", "this", err)
	}
}

func TestExecutingACommandWithErrors(t *testing.T) {
	out, err := ExecuteCommand("", "")
	if err != nil {
		t.Fatalf("expected: error when executing an invalid command, got: no errors")
	}

	if out == "" {
		t.Errorf("expected: error message regarding the faulty executing of an invalid command, got: no error messages")
	}
}
