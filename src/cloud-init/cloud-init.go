package main

import (
	"env"
	"fmt"
)

func main() {
	argMap := map[string]string{
		"user-data": "dummyUserData",
		"meta-data": "dummyMetaData",
	}

	fmt.Println(env.FlattenArguments(argMap))
}
