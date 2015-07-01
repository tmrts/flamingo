package nss

import (
	"strconv"
	"strings"
)

type GroupEntry struct {
	GroupName     string
	GID           int
	IsSystemGroup bool

	Members []string
}

// Parses NSS GroupDatabase Entry.
// Example:
//	 group:x:1:user1,user2 is turned into:
//	 Group {
//	 	Name:			 "group",
//	 	GID:			 "1",
//	 	Members:		 []string {"user1", "user2"},
//	 	IsSystemAccount: true,
//	 }
func parseGroupEntry(groupEntry string) *GroupEntry {
	groupInfo := strings.Split(groupEntry, ":")
	groupID, _ := strconv.Atoi(groupInfo[2])

	groupMembers := strings.Split(groupInfo[3], ",")

	return &GroupEntry{
		GroupName:     groupInfo[0],
		GID:           groupID,
		IsSystemGroup: groupID < 1000,
		Members:       groupMembers,
	}
}

// GetGroup queries the NSS Group Database.
func GetGroupEntry(s Service, key string) (*GroupEntry, error) {
	entry, err := s.GetEntryFrom(GroupDatabase, key)
	if err != nil {
		return nil, err
	}

	return parseGroupEntry(entry), nil
}
