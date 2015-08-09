package nss

import (
	"strconv"
	"strings"
)

// PasswdEntry is the representation of Name Switch Service
// 'passwd' database entry fields.
type PasswdEntry struct {
	UserName        string
	UID             int
	GID             int
	GECOS           string
	HomeDir         string
	DefaultShell    string
	IsSystemAccount bool
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
func parsePasswdEntry(passwdEntry string) *PasswdEntry {
	userInfo := strings.Split(passwdEntry, ":")
	userID, _ := strconv.Atoi(userInfo[2])
	groupID, _ := strconv.Atoi(userInfo[3])

	return &PasswdEntry{
		UserName:        userInfo[0],
		UID:             userID,
		GID:             groupID,
		GECOS:           userInfo[4],
		HomeDir:         userInfo[5],
		DefaultShell:    userInfo[5],
		IsSystemAccount: userID < 1000,
	}
}

func GetPasswdEntry(s Service, key string) (*PasswdEntry, error) {
	ent, err := s.GetEntryFrom(UserDatabase, key)
	if err != nil {
		return nil, err
	}

	return parsePasswdEntry(ent), nil
}
