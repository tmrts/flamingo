package context

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"

	"github.com/TamerTas/cloud-init/pkg/env"
	"github.com/TamerTas/cloud-init/pkg/utils"
)

type UserShadowEntry struct {
	PasswordHash string `flag:"password"`
}

type User struct {
	Name            string
	UID             string `flag:"uid"`
	Description     string `flag:"comment"`
	HomeDir         string `flag:"home"`
	DefaultShell    string `flag:"shell"`
	IsSystemAccount bool   `flag:"system"`

	PrimaryGroup    string   `flag:"gid"`
	SecondaryGroups []string `flag:"groups"`

	// TODO: Refactor flags out of User
	// User creation flags
	DontCreateHomeDir           bool   `flag:"no-create-home"`
	DirectoryTemplate           string `flag:"skel"`
	CreateUserGroupWithSameName string `flag:"user-group"`
}

// CreateNewUser adds the supplied User to the users.
func CreateNewUser(usr User) error {
	args := utils.GetArgumentFormOfStruct(usr)

	args = append(args, usr.Name)

	_, err := env.ExecuteCommand("useradd", args...)

	return err
}

// Parses NSS PasswdDatabase Entry.
// Example:
//	 user:x:1000:1000:A normal user.:/home/user:/bin/bash is turned into:
//	 User {
//	 	Name: "group",
//	 	UID: "1000",
//	 	PrimaryGroup: "1000",
//	 	Description: "A normal user",
//	 	HomeDir: "/home/user",
//	 	DefaultShell: "/bin/bash",
//	 	IsSystemAccount: false,
//	 }
func parsePasswdEntry(passwdEntry string) *User {
	userInfo := strings.Split(passwdEntry, ":")
	userID, _ := strconv.Atoi(userInfo[2])
	return &User{
		Name:            userInfo[0],
		UID:             userInfo[2],
		PrimaryGroup:    userInfo[3],
		Description:     userInfo[4],
		HomeDir:         userInfo[5],
		DefaultShell:    userInfo[5],
		IsSystemAccount: userID < 1000,
	}
}

// GetUser queries the NSS Passwd Database.
func GetUser(key string) (*User, error) {
	entry, err := env.GetEntryFrom(env.UserDatabase, key)
	if err != nil {
		return nil, err
	}
	user := parsePasswdEntry(entry)

	return user, nil
}

// Parses NSS UserShadowDatabase Entry.
func parseUserShadowEntry(shadowEntry string) *UserShadowEntry {
	userShadowInfo := strings.Split(shadowEntry, ":")

	return &UserShadowEntry{
		PasswordHash: userShadowInfo[1],
	}
}

// GetUserShadowEntry queries the NSS User Shadow Database.
func GetUserShadowEntry(key string) (*UserShadowEntry, error) {
	shadowEntry, err := env.GetEntryFrom(env.UserShadowDatabase, key)
	if err != nil {
		return nil, err
	}

	userShadowEntry := parseUserShadowEntry(shadowEntry)

	return userShadowEntry, nil
}

// SetPassword changes the user's password to the given password hash.
func (usr *User) SetPassword(passwordHash string) error {
	cmd := exec.Command("chpasswd", "-e")

	// Securely pass the password hash
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	passwordPair := fmt.Sprintf("%s:%s", usr.Name, passwordHash)

	_, err = stdin.Write([]byte(passwordPair))
	if err != nil {
		return err
	}

	stdin.Close()

	return cmd.Wait()
}
