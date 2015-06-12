package main

import "fmt"

func main() {
	argMap := map[string]string{
		"user-data": "dummyUserData",
		"meta-data": "dummyMetaData",
	}

	fmt.Println(argMap)
}
