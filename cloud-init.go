package main

import (
	"env"
	"fmt"
	"os/user"
)

func main() {
	user.User
	argMap := map[string]string{
		"user-data": "/etc/file",
	}
	fmt.Println(env.FlattenArguments(argMap))
}
