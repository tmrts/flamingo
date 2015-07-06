package identity_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/TamerTas/flamingo/pkg/sys"
	"github.com/TamerTas/flamingo/pkg/sys/identity"
	"github.com/TamerTas/flamingo/pkg/util/testutil"
)

func TestUserCreation(t *testing.T) {
	Convey("Given a user object and an executor", t, func() {
		user := identity.User{
			Name:         "newUser",
			PasswordHash: "PHASH",
			UID:          "1001",
			GID:          "1001",
		}

		exec := sys.NewMockExecutor()

		idmgr := &identity.ManagerImplementation{
			Executor: exec,
		}

		Convey("It should create a new user in the system", func() {
			exec.OutStr <- ""
			exec.OutErr <- nil

			err := idmgr.CreateUser(user)

			So(<-exec.Exec, ShouldEqual, "useradd")
			So(<-exec.Args, testutil.ShouldSetEqual, []string{"newUser", "--password=PHASH", "--uid=1001", "--gid=1001"})

			So(err, ShouldEqual, nil)
		})
	})
}

func TestUserSetPassword(t *testing.T) {
	Convey("Given a user name and a password hash and an executor", t, func() {
		uname, phash := "existentUser", "PASSWORD_HASH"

		exec := sys.NewMockExecutor()

		idmgr := &identity.ManagerImplementation{
			Executor: exec,
		}

		Convey("It should change the password of the user", func() {
			exec.OutStr <- ""
			exec.OutErr <- nil

			err := idmgr.SetGroupPassword(uname, phash)

			So(<-exec.Exec, ShouldEqual, "chpasswd")
			So(<-exec.Args, testutil.ShouldSetEqual, []string{"-e", "existentUser:PASSWORD_HASH"})

			So(err, ShouldEqual, nil)
		})
	})
}

func TestGroupCreation(t *testing.T) {
	Convey("Given a group object and an executor", t, func() {
		user := identity.Group{
			Name:         "newGroup",
			PasswordHash: "PHASH",
			GID:          "1002",
		}

		exec := sys.NewMockExecutor()

		idmgr := &identity.ManagerImplementation{
			Executor: exec,
		}

		Convey("It should create a new group in the system", func() {
			exec.OutStr <- ""
			exec.OutErr <- nil

			err := idmgr.CreateGroup(user)

			So(<-exec.Exec, ShouldEqual, "groupadd")
			So(<-exec.Args, testutil.ShouldSetEqual, []string{"newGroup", "--password=PHASH", "--gid=1002"})

			So(err, ShouldEqual, nil)
		})
	})
}

func TestGroupSetPassword(t *testing.T) {
	Convey("Given a group name and a password hash and an executor", t, func() {
		gname, phash := "existentGroup", "PASSWORD_HASH"

		exec := sys.NewMockExecutor()

		idmgr := &identity.ManagerImplementation{
			Executor: exec,
		}

		Convey("It should change the password of the group", func() {
			exec.OutStr <- ""
			exec.OutErr <- nil

			err := idmgr.SetGroupPassword(gname, phash)

			So(<-exec.Exec, ShouldEqual, "groupmod")
			So(<-exec.Args, testutil.ShouldSetEqual, []string{"existentGroup", "--password=PASSWORD_HASH"})

			So(err, ShouldEqual, nil)
		})
	})
}
