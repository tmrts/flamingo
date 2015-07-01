package nss_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/tmrts/flamingo/pkg/sys"
	"github.com/tmrts/flamingo/pkg/sys/nss"
	"github.com/tmrts/flamingo/pkg/util/rand"
)

func TestNSSErrors(t *testing.T) {
	Convey("Given an NSS server", t, func() {
		server := nss.New(sys.DefaultExecutor)

		Convey("When queried with a wrong database", func() {
			var (
				db  nss.Database = "InvalidDatabase"
				key string       = "some_entry"
			)

			Convey("It should return an error message containing: Unknown database", func() {
				_, err := server.GetEntryFrom(db, key)
				So(err, ShouldNotBeNil)

				So(err.Error(), ShouldContainSubstring, "Unknown database")
			})
		})

		Convey("When queried with a non-existent user", func() {
			var (
				db  nss.Database = nss.UserDatabase
				key string       = rand.String(10)
			)

			Convey("It should return an error message containing: Key could not be found", func() {
				_, err := server.GetEntryFrom(db, key)

				So(err, ShouldNotBeNil)

				So(err.Error(), ShouldContainSubstring, "Key could not be found")
			})
		})
	})
}

func TestUserIDLookup(t *testing.T) {
	Convey("Given a user name", t, func() {
		uname := "root"

		Convey("It should get the UID and GID of the user", func() {
			uid, gid, err := nss.GetIDsForUser(uname)

			So(err, ShouldBeNil)

			So(uid, ShouldEqual, 0)
			So(gid, ShouldEqual, 0)
		})
	})
}
