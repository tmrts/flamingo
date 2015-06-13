package env

import (
	"testing"

	"github.com/TamerTas/cloud-init/pkg/utils"
)

func TestScriptValidation(t *testing.T) {
	testCases := []struct {
		Case     string
		Content  string
		IsScript bool
	}{
		{
			Case:     "fileWithoutShabang",
			Content:  "Some stuff\n",
			IsScript: false,
		},
		{
			Case:     "fileWithInvalidShabang",
			Content:  "# !\nSome command\n",
			IsScript: false,
		},
		{
			Case:     "emptyScriptWithShabang",
			Content:  "#! /usr/bin/env bash\n\n",
			IsScript: true,
		},
		{
			Case:     "scriptWithShabang",
			Content:  "#! /usr/bin/env bash\nprintf ${PWD}\n",
			IsScript: true,
		},
	}

	for _, v := range testCases {
		tmpFile, err := utils.CreateTempFile(v.Content)
		if err != nil {
			t.Fatalf("error creating a temporary file -> %v", err)
		}

		isScript, err := HasShabang(tmpFile.Name())
		if err != nil {
			t.Errorf("error validating a script -> %v", err)
		}

		if v.IsScript != isScript {
			t.Errorf("wrong result from script validation case %v -> expected: %v, got: %v", v.Case, v.IsScript, isScript)
		}
	}
}

func TestScriptExecution(t *testing.T) {
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
