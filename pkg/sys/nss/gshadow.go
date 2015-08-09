package nss

import "strings"

// GroupShadowEntry is the representation of Name Switch Service
// 'gshadow' database entry fields.
type GroupShadowEntry struct {
	GroupName    string
	PasswordHash string
	Admins       []string
	Members      []string
}

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

// GetGroup queries the Name Switch Service 'gshadow' Database.
// It returns the parsed 'gshadow' entry belonging to given group key.
func GetGroupShadowEntry(s Service, groupKey string) (*GroupShadowEntry, error) {
	shadowEntry, err := s.GetEntryFrom(GroupShadowDatabase, groupKey)
	if err != nil {
		return nil, err
	}

	return parseGroupShadowEntry(shadowEntry), nil
}
