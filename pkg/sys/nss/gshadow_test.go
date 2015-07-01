package nss_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/tmrts/flamingo/pkg/sys/nss"
	"github.com/tmrts/flamingo/pkg/util/testutil"
)

func TestGroupShadowDatabaseEntryRetrieval(t *testing.T) {
	Convey("Given an nss server and a group name", t, func() {
		gname := "root"

		server := &nss.Server{
			Exec: &fakeExecutor{
				expectedCmd:    "getent",
				expectedResult: "root:$4$SALT$PASSWORD_HASH:admin1,admin2:user1,user2",
			},
		}

		Convey("It should query the NSS gshadow database and get its entry", func() {
			root, err := nss.GetGroupShadowEntry(server, gname)
			So(err, ShouldBeNil)

			So(root.GroupName, ShouldEqual, "root")

			So(root.PasswordHash, ShouldEqual, "$4$SALT$PASSWORD_HASH")

			So(root.Admins, testutil.ShouldSetEqual, []string{"admin1", "admin2"})

			So(root.Members, testutil.ShouldSetEqual, []string{"user1", "user2"})
		})
	})
}
