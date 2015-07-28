package file_test

import (
	"os"
	"strconv"
	"syscall"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/tmrts/flamingo/pkg/file"
	"github.com/tmrts/flamingo/pkg/sys/nss"
)

func TestEnsureDirectoryExists(t *testing.T) {
	Convey("Given a unique dir name, a user name and file permissions", t, func() {
		currentUser, err := nss.GetCurrentUser()
		So(err, ShouldBeNil)

		dirname, err := file.UniqueName("/tmp", "filetests")
		So(err, ShouldBeNil)

		perms := os.FileMode(0600)

		userID, _ := strconv.Atoi(currentUser.Uid)
		groupID, _ := strconv.Atoi(currentUser.Gid)

		Convey("It should create a new directory with the given permissions if it doesn't exist", func() {
			err := file.EnsureDirectoryExists(dirname, perms, userID, groupID)
			So(err, ShouldBeNil)
			defer os.Remove(dirname)

			fi, err := os.Lstat(dirname)
			So(err, ShouldBeNil)

			So(fi.IsDir(), ShouldBeTrue)

			uid := strconv.Itoa(int(fi.Sys().(*syscall.Stat_t).Uid))
			So(currentUser.Uid, ShouldEqual, uid)

			gid := strconv.Itoa(int(fi.Sys().(*syscall.Stat_t).Gid))
			So(currentUser.Gid, ShouldEqual, gid)
		})

		Convey("It should change the ownership and the file permissions if it exists", func() {
			t.SkipNow()
			err := file.EnsureDirectoryExists(dirname, perms, userID, groupID)
			So(err, ShouldBeNil)
			defer os.Remove(dirname)

			fi, err := os.Lstat(dirname)
			So(err, ShouldBeNil)

			So(fi.IsDir(), ShouldBeTrue)

			uid := strconv.Itoa(int(fi.Sys().(*syscall.Stat_t).Uid))
			So(currentUser.Uid, ShouldEqual, uid)

			gid := strconv.Itoa(int(fi.Sys().(*syscall.Stat_t).Gid))
			So(currentUser.Gid, ShouldEqual, gid)
		})
	})
}
