package env

import (
	"os/user"
	"testing"
)

var (
	//newGroup := Group{
	//"newGroup",
	//dummyPassword,
	//}

	dummyUser = User{
		Name:    "dummyUser",
		HomeDir: "/home/dummyUser",
	}
)

func TestNewUser(t *testing.T) {
	err := NewUser(dummyUser)
	if err != nil {
		t.Fatalf("user creation failed -> %v", err)
	}

	newUser, err := user.Lookup("dummyUser")
	if err != nil {
		t.Fatalf("user look-up failed -> %v", err)
	}

	if newUser.Username != "dummyUser" {
		t.Fatalf("user name discrepancy -> %v", newUser.Username)
	}
}
