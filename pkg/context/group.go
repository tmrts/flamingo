package context

import (
	"strconv"
	"strings"

	"github.com/TamerTas/cloud-init/pkg/env"
	"github.com/TamerTas/cloud-init/pkg/utils"
)

type GroupShadowEntry struct {
	Name           string
	PasswordHash   string
	Administrators []string
	Members        []string
}

type Group struct {
	Name            string
	GID             string `flag:"gid"`
	IsSystemAccount bool   `flag:"system"`

	Members []string
}

// CreateNewGroup adds the supplied Group to the system groups.
func CreateNewGroup(grp Group) error {
	args := utils.GetArgumentFormOfStruct(grp)

	args = append(args, grp.Name)

	_, err := env.ExecuteCommand("groupadd", args...)

	return err
}

// Parses NSS GroupDatabase Entry.
// Example:
//	 group:x:1:user1,user2 is turned into:
//	 Group {
//	 	Name: "group",
//	 	GID: 1,
//	 	Members: user1, user2,
//	 	IsSystemAccount: true,
//	 }

func parseGroupEntry(groupEntry string) *Group {
	groupInfo := strings.Split(groupEntry, ":")
	groupID, _ := strconv.Atoi(groupInfo[2])

	groupMembers := strings.Split(groupInfo[3], ",")

	return &Group{
		Name:            groupInfo[0],
		GID:             groupInfo[2],
		Members:         groupMembers,
		IsSystemAccount: groupID < 1000,
	}
}

// GetGroup queries the NSS Group Database.
func GetGroup(key string) (*Group, error) {
	entry, err := env.GetEntryFrom(env.GroupDatabase, key)
	if err != nil {
		return nil, err
	}
	group := parseGroupEntry(entry)

	return group, nil
}

// Parses NSS GroupShadowDatabase Entry.
func parseGroupShadowEntry(groupEntry string) *GroupShadowEntry {
	groupInfo := strings.Split(groupEntry, ":")

	groupAdmins := strings.Split(groupInfo[2], ",")
	groupMembers := strings.Split(groupInfo[3], ",")

	return &GroupShadowEntry{
		Name:           groupInfo[0],
		PasswordHash:   groupInfo[1],
		Administrators: groupAdmins,
		Members:        groupMembers,
	}
}

// GetGroup queries the NSS Group Shadow Database.
func GetGroupShadowEntry(key string) (*GroupShadowEntry, error) {
	shadowEntry, err := env.GetEntryFrom(env.GroupShadowDatabase, key)
	if err != nil {
		return nil, err
	}

	groupShadowEntry := parseGroupShadowEntry(shadowEntry)

	return groupShadowEntry, nil
}

// SetPassword changes the group's password to the given password hash.
func (grp *Group) SetPassword(passwordHash string) error {
	_, err := env.ExecuteCommand("groupmod", grp.Name, "--password="+passwordHash)

	return err
}
