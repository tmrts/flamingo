package env

import (
	"os/user"
	"testing"
)

func TestNewUser(t *testing.T) {
	dummyPassword := Password{"hash"}

	//newGroup := Group{
	//"newGroup",
	//dummyPassword,
	//}

	dummyUser := User{
		Name:     "uname",
		Password: dummyPassword,
	}

	err := NewUser(dummyUser)
	if err != nil {
		t.Fatalf("User Creation Failed %v", err)
	}

	newUser, err := user.Lookup("uname")
	if err != nil {
		t.Fatalf("User Look-up Failed %v", err)
	}

	if newUser.Username != "uname" {
		t.Fatalf("user name discrepancy %v", newUser.Username)
	}
}
