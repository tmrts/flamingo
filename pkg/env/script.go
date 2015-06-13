package env

import "io/ioutil"

func HasShabang(filename string) (bool, error) {
	byteContent, _ := ioutil.ReadFile(filename)

	content := string(byteContent)

	return content[0] == '#' && content[1] == '!', nil
}

func WriteScriptTo(filename string, contents string) error {
	return ioutil.WriteFile(filename, []byte(contents), 0744)
}

func ExecuteScript(scriptPath string) (string, error) {
	return ExecuteCommand("bash", scriptPath)
}
