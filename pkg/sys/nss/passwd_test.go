package nss_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/tmrts/flamingo/pkg/sys/nss"
)

func TestPasswdDatabaseEntryRetrieval(t *testing.T) {
	Convey("Given an nss server and a user name", t, func() {
		uname := "root"

		server := &nss.Server{
			Exec: &fakeExecutor{
				expectedCmd:    "getent",
				expectedResult: "root:x:0:0:root:/root:/bin/bash",
			},
		}

		Convey("It should query the NSS passwd database and get its entry", func() {
			root, err := nss.GetPasswdEntry(server, uname)
			So(err, ShouldBeNil)

			So(root.UserName, ShouldEqual, "root")

			So(root.UID, ShouldEqual, 0)
			So(root.GID, ShouldEqual, 0)

			So(root.HomeDir, ShouldEqual, "/root")

			So(root.IsSystemAccount, ShouldBeTrue)
		})
	})
}
