package env

import (
	"testing"

	"github.com/TamerTas/cloud-init/pkg/utils"
)

func TestScriptValidation(t *testing.T) {
	scriptContent := "#! /usr/bin/env bash\n"
	scriptContent += "echo -n Yeehaw\n"
	scriptContent += "\n"

	testScript, err := utils.CreateTempFile(scriptContent)
	if err != nil {
		t.Fatalf("error creating a temporary script file -> %v", err)
	}

	out, err := ExecuteScript(testScript.Name())
	if err != nil {
		t.Errorf("error executing a script -> %v", err)
	}

	if out != "Yeehaw" {
		t.Errorf("error executing a script expected: Yeehaw, got: %v", out)
	}

	testScript.Close()
}
