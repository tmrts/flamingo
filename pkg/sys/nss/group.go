package nss

import (
	"strconv"
	"strings"
)

// GroupEntry is the representation of Name Switch Service
// 'group' database entry fields.
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

// GetGroup queries the Name Switch Service 'group' Database for
// a given group key. Group key is usually the name of that group.
// It returns the parsed group entry.
func GetGroupEntry(s Service, groupKey string) (*GroupEntry, error) {
	entry, err := s.GetEntryFrom(GroupDatabase, groupKey)
	if err != nil {
		return nil, err
	}

	return parseGroupEntry(entry), nil
}
