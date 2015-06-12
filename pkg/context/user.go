package env

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
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

func CreateNewUser(usr User) error {
	args := GetArgumentFormOfStruct(usr)

	args = append(args, usr.Name)

	_, err := ExecuteCommand("useradd", args...)

	return err
}

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

func GetUser(key string) (*User, error) {
	entry, err := GetEntryFrom(UserDatabase, key)
	if err != nil {
		return nil, err
	}
	user := parsePasswdEntry(entry)

	return user, nil
}

func parseUserShadowEntry(shadowEntry string) *UserShadowEntry {
	userShadowInfo := strings.Split(shadowEntry, ":")

	return &UserShadowEntry{
		PasswordHash: userShadowInfo[1],
	}
}

func GetUserShadowEntry(key string) (*UserShadowEntry, error) {
	shadowEntry, err := GetEntryFrom(UserShadowDatabase, key)
	if err != nil {
		return nil, err
	}

	userShadowEntry := parseUserShadowEntry(shadowEntry)

	return userShadowEntry, nil
}

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
