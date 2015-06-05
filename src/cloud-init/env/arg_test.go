package env

import "testing"

// Checks whether two slices contain the same elements.
func setEquals(lhs, rhs []string) bool {
	if len(lhs) != len(rhs) {
		return false
	}

	lhsSet := make(map[string]bool)
	for _, lval := range lhs {
		lhsSet[lval] = true
	}

	for _, rval := range rhs {
		if lhsSet[rval] != true {
			return false
		}
	}

	return true
}

func TestStructToArgsConversion(t *testing.T) {
	fakeUserInfo := struct {
		UserID          string   `flag:"uid"`
		Comment         string   `flag:"comment"`
		IsSystemAccount bool     `flag:"system"`
		Items           []string `flag:"items"`
	}{
		UserID:          "990",
		Comment:         "This is a Comment.",
		IsSystemAccount: true,
		Items:           []string{"item1", "item2", "item3"},
	}

	argSlice := GetArgumentFormOfStruct(fakeUserInfo)

	expectedArgs := []string{
		"--uid=990",
		"--system",
		"--comment=This is a Comment.",
		"--items=item1,item2,item3",
	}

	if setEquals(expectedArgs, argSlice) != true {
		t.Errorf("arg slice mismatch\nexpected:\n\t%v\ngot:\n\t%v\n", argSlice, expectedArgs)
	}
}
