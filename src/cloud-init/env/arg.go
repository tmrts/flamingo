package env

import "fmt"

func FlattenArguments(argMap map[string]string) []string {
	args := []string{}

	for flag, value := range argMap {
		argument := fmt.Sprintf(fmt.Sprintf("--%v=%v", flag, value))

		args = append(args, argument)
	}

	return args
}
