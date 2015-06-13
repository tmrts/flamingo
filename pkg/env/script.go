package env

import "io/ioutil"

func WriteScriptTo(fileName string, contents string) error {
	return ioutil.WriteFile(fileName, []byte(contents), 0744)
}

func ExecuteScript(scriptPath string) (string, error) {
	return ExecuteCommand("bash", scriptPath)
}
