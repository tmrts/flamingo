package main

import (
	"env"
	"fmt"
)

func main() {
	argMap := map[string]string{
		"user-data": "/etc/file",
	}
	fmt.Println(env.FlattenArguments(argMap))
}
