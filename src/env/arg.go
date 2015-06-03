package env

import "fmt"

func FlattenArguments(argMap map[string]string) string {
	arguments := ""

	for flag, value := range argMap {
		arguments += fmt.Sprintf("--%v=%v ", flag, value)
	}

	return arguments
}
