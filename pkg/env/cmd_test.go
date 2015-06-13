package env

import "testing"

func TestExecutingACommandWithErrors(t *testing.T) {
	_, err := ExecuteCommand("", "")
	if err == nil {
		t.Fatalf("expected: error when executing an invalid command, got: no errors")
	}
}

func TestExecuteCommand(t *testing.T) {
	out, err := ExecuteCommand("echo", "-n", "fi fye fo fum")
	if err != nil {
		t.Fatalf("error executing echo command -> %v", err)
	}

	if out != "fi fye fo fum" {
		t.Errorf("wrong output for echo command -> expected: %v, got: %v", "this", err)
	}
}
