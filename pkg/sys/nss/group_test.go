package nss_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/tmrts/flamingo/pkg/sys/nss"
	"github.com/tmrts/flamingo/pkg/util/testutil"
)

func TestGroupDatabaseEntryRetrieval(t *testing.T) {
	Convey("Given an nss server and a group name", t, func() {
		gname := "root"

		server := &nss.Server{
			Exec: &fakeExecutor{
				expectedCmd:    "getent",
				expectedResult: "root:x:0:user1,user2,user3",
			},
		}

		Convey("It should query the NSS group database and get its entry", func() {
			root, err := nss.GetGroupEntry(server, gname)
			So(err, ShouldBeNil)

			So(root.GroupName, ShouldEqual, "root")

			So(root.GID, ShouldEqual, 0)

			So(root.IsSystemGroup, ShouldBeTrue)

			So(root.Members, testutil.ShouldSetEqual, []string{"user1", "user2", "user3"})
		})
	})
}
