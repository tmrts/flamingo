package env

import (
	"strings"
	"testing"
)

func TestNSSErrors(t *testing.T) {
	_, err := GetEntryFrom("InvalidDatabase", "some entry")
	if strings.Contains(err.Error(), "Unknown database") == false {
		t.Errorf("wrong error msg retrieving entry from an invalid database -> %v", err)
	}

	_, err = GetEntryFrom(UserDatabase, "someHopefullyNonExistantUser")
	if strings.Contains(err.Error(), "Key could not be found") == false {
		t.Errorf("wrong error msg retrieving a non-existant entry -> %v", err)
	}
}
