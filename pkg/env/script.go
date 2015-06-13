package env

import "io/ioutil"

func WriteScriptTo(filename string, contents string) error {
	return ioutil.WriteFile(filename, []byte(contents), 0744)
}

func ExecuteScript(scriptPath string) (string, error) {
	return ExecuteCommand("bash", scriptPath)
}
