package env

import "testing"

func TestPasswdDatabaseEntryRetrieval(t *testing.T) {
	root, err := GetUser("root")
	if err != nil {
		t.Fatalf("user entry retrieval error -> %v", err)
	}

	if root.Name != "root" {
		t.Fatalf("wrong user name for root account")
	}

	if root.UID != "0" {
		t.Fatalf("wrong user ID for root account")
	}

	if root.PrimaryGroup != "0" {
		t.Fatalf("wrong group ID for root account")
	}

	if root.HomeDir != "/root" {
		t.Fatalf("wrong home directory for root account")
	}

	if root.IsSystemAccount != true {
		t.Fatalf("wrong account type for root account")
	}
}

func TestShadowDatabaseEntryRetrieval(t *testing.T) {
	root, err := GetUserShadowEntry("root")
	if err != nil {
		t.Fatalf("user entry retrieval error -> %v", err)
	}

	if root.PasswordHash != "locked" {
		t.Fatalf("wrong user name for root account")
	}
}

func TestCreateNewUser(t *testing.T) {
	dummyUser := User{
		Name:              "dummyUser",
		HomeDir:           "/etc",
		Description:       "This is a description.",
		DontCreateHomeDir: true,
	}

	err := CreateNewUser(dummyUser)
	if err != nil {
		t.Fatalf("user creation failed -> %v", err)
	}

	newUser, err := GetUser("dummyUser")
	if err != nil {
		t.Fatalf("user look-up failed -> %v", err)
	}

	if newUser.Name != "dummyUser" {
		t.Errorf("user name discrepancy -> %v", newUser.Name)
	}

	if newUser.Description != "This is a description." {
		t.Errorf("user name discrepancy -> %v", newUser.Name)
	}
}

func TestSetUserPassword(t *testing.T) {
	root, err := GetUser("root")
	if err != nil {
		t.Fatalf("user look-up failed -> %v", err)
	}

	root.SetPassword("PASSWORD_HASH")

	rootShadowEntry, err := GetUserShadowEntry("root")
	if err != nil {
		t.Fatalf("shadow entry retrieval failed -> %v", err)
	}

	if rootShadowEntry.PasswordHash != "PASSWORD_HASH" {
		t.Fatal("setting user password failed for root")
	}
}
