package context

import "testing"

func TestCreateNewGroup(t *testing.T) {
	fakeGroup := Group{
		Name:            "fakeGroup",
		IsSystemAccount: true,
	}

	err := CreateNewGroup(fakeGroup)
	if err != nil {
		t.Fatalf("group creation failed -> %v", err)
	}

	newGroup, err := GetGroup(fakeGroup.Name)
	if err != nil {
		t.Fatalf("error retrieving new group from db -> %v", err)
	}

	if newGroup.Name != "fakeGroup" {
		t.Errorf("group name discrepancy -> %v", err)
	}

	fakeGroup.SetPassword("PASSWORD_HASH")

	groupShadowEntry, err := GetGroupShadowEntry("fakeGroup")
	if err != nil {
		t.Fatalf("shadow entry retrieval failed -> %v", err)
	}

	if groupShadowEntry.PasswordHash != "PASSWORD_HASH" {
		t.Fatal("setting user password failed for fakeGroup")
	}
}
