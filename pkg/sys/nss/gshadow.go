package nss

import "strings"

type GroupShadowEntry struct {
	GroupName    string
	PasswordHash string
	Admins       []string
	Members      []string
}

// Parses NSS GroupShadowDatabase Entry.
func parseGroupShadowEntry(groupEntry string) *GroupShadowEntry {
	groupInfo := strings.Split(groupEntry, ":")

	groupAdmins := strings.Split(groupInfo[2], ",")
	groupMembers := strings.Split(groupInfo[3], ",")

	return &GroupShadowEntry{
		GroupName:    groupInfo[0],
		PasswordHash: groupInfo[1],
		Admins:       groupAdmins,
		Members:      groupMembers,
	}
}

// GetGroup queries the NSS Group Shadow Database.
func GetGroupShadowEntry(s Service, key string) (*GroupShadowEntry, error) {
	shadowEntry, err := s.GetEntryFrom(GroupShadowDatabase, key)
	if err != nil {
		return nil, err
	}

	return parseGroupShadowEntry(shadowEntry), nil
}
