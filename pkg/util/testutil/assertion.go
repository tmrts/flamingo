package testutil

import (
	"fmt"
	"reflect"
)

// Checks whether two slices contain the same elements.
func ShouldSetEqual(actual interface{}, expected ...interface{}) (msg string) {
	actualSlice, expectedSlice := actual.([]string), expected[0].([]string)

	msg = fmt.Sprintf("Expected:\t%q\nActual:\t%q\n(Should have the same elements)", expectedSlice, actualSlice)

	if len(actualSlice) != len(expectedSlice) {
		return
	}

	actualSet := make(map[string]bool)
	for _, a := range actualSlice {
		actualSet[a] = true
	}

	for _, e := range expectedSlice {
		if actualSet[e] != true {
			return
		}
	}

	expectedSet := make(map[string]bool)
	for _, e := range expectedSlice {
		expectedSet[e] = true
	}

	for _, a := range actualSlice {
		if expectedSet[a] != true {
			return
		}
	}

	return ""
}

func ShouldConsistOf(actual interface{}, expected ...interface{}) (msg string) {
	msg = fmt.Sprintf("Expected:\n%q\nActual:\n%q\n(Should consist of the given elements)", expected, actual)

	switch actual.(type) {
	case []int:
		for i, v := range actual.([]int) {
			if !reflect.DeepEqual(v, expected[i]) {
				return
			}
		}
	case []string:
		for i, v := range actual.([]string) {
			if !reflect.DeepEqual(v, expected[i]) {
				return
			}
		}
	}

	return ""
}

func ShouldBeSuperSetOf(super interface{}, sub ...interface{}) (msg string) {
	superSlice, subSlice := super.([]string), sub[0].([]string)

	msg = fmt.Sprintf("SubSet:\t%q\nSuperSet:\t%q\n(Should contain every element of the given slice)", subSlice, superSlice)

	if len(superSlice) < len(subSlice) {
		return
	}

	superSet := make(map[string]bool)
	for _, e := range superSlice {
		superSet[e] = true
	}

	for _, a := range subSlice {
		if _, ok := superSet[a]; ok != true {
			return
		}
	}

	return ""
}
