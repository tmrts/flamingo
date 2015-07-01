package nss

import "strings"

type ShadowEntry struct {
	UserName      string
	PasswordHash  string
	ChangedAgo    string
	ChangedBefore string
	ExpiryPeriod  string
	DisabledSince string
}

// Parses NSS UserShadowDatabase Entry.
func parseShadowEntry(shadowEntry string) *ShadowEntry {
	userShadowInfo := strings.Split(shadowEntry, ":")

	return &ShadowEntry{
		UserName:     userShadowInfo[0],
		PasswordHash: userShadowInfo[1],
	}
}

// ShadowEntry queries the NSS User Shadow Database.
func GetShadowEntry(s Service, key string) (*ShadowEntry, error) {
	ent, err := s.GetEntryFrom(UserShadowDatabase, key)
	if err != nil {
		return nil, err
	}

	return parseShadowEntry(ent), nil
}
