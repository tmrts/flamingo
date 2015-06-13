package env

func ExecuteScript(scriptPath string) (string, error) {
	return ExecuteCommand("bash", scriptPath)
}
