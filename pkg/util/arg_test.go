package util_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/tmrts/flamingo/pkg/util"
	"github.com/tmrts/flamingo/pkg/util/testutil"
)

func TestStructToArgsConversion(t *testing.T) {
	Convey("Given a struct with `flag` tags", t, func() {
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

		Convey("It should serialize the struct into flag form with the given flag tags", func() {
			argSlice := util.GetArgumentFormOfStruct(fakeUserInfo)

			expectedArgs := []string{
				"--uid=990",
				"--system",
				"--comment=This is a Comment.",
				"--items=item1,item2,item3",
			}

			So(argSlice, testutil.ShouldSetEqual, expectedArgs)
		})
	})

	Convey("Given an uninitialized struct with `flag` tags", t, func() {
		fakeUserInfo := struct {
			IsAvailable bool     `flag:"available"`
			Item        string   `flag:"item"`
			Items       []string `flag:"items"`
			UserID      int      `flag:"id"`
			UserIDs     []int    `flag:"ids"`
		}{}

		Convey("It should return an empty slice", func() {
			argSlice := util.GetArgumentFormOfStruct(fakeUserInfo)

			So(argSlice, ShouldBeEmpty)
		})
	})
}
