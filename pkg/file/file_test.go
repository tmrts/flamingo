package file_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/tmrts/flamingo/pkg/file"
	"github.com/tmrts/flamingo/pkg/sys/nss"
)

func TestNewFileCreation(t *testing.T) {
	Convey("Given a file path", t, func() {
		fname := filepath.Join(os.TempDir(), "__test_filepath")

		Convey("It should create an empty file", func() {
			err := file.New(fname)
			So(err, ShouldBeNil)
			defer os.Remove(fname)

			_, err = os.Stat(fname)
			So(err, ShouldBeNil)
		})

		Convey("With optional argument Contents", func() {
			c := "This is a file with text in it.\n"
			Convey("It should create a new file with the given contents", func() {
				err := file.New(fname, file.Contents(c))
				So(err, ShouldBeNil)
				defer os.Remove(fname)

				out, err := ioutil.ReadFile(fname)

				So(err, ShouldBeNil)
				So(string(out), ShouldEqual, c)
			})
		})

		Convey("With optional argument Permissions", func() {
			p := os.FileMode(0700)

			Convey("It should create a new file with the same permissions", func() {
				err := file.New(fname, file.Permissions(p))
				So(err, ShouldBeNil)
				defer os.Remove(fname)
				fi, err := os.Lstat(fname)

				So(err, ShouldBeNil)
				So(fi.Mode(), ShouldEqual, p)
			})
		})

		Convey("With optional argument User ID and/or Group ID", func() {
			userID, groupID := os.Getuid(), os.Getgid()

			Convey("It should create a new file belonging to given IDs", func() {
				err := file.New(fname, file.UID(userID), file.GID(groupID))
				So(err, ShouldBeNil)
				defer os.Remove(fname)
				fi, err := os.Lstat(fname)
				So(err, ShouldBeNil)

				uid := fi.Sys().(*syscall.Stat_t).Uid
				So(uid, ShouldEqual, userID)

				gid := fi.Sys().(*syscall.Stat_t).Gid
				So(gid, ShouldEqual, groupID)
			})
		})
	})
}

func TestRandomFileNameGeneration(t *testing.T) {
	Convey("Given a dirname and a prefix string", t, func() {
		dirname, prefix := os.TempDir(), "fnamegeneration_test"

		Convey("It should return a unique file name in the dir with the given prefixes", func() {
			fname, err := file.UniqueName(dirname, prefix)
			So(err, ShouldBeNil)

			_, err = os.Lstat(fname)
			So(os.IsNotExist(err), ShouldBeTrue)

			So(strings.Contains(fname, dirname), ShouldBeTrue)
			So(strings.Contains(fname, prefix), ShouldBeTrue)
		})
	})
}

func TestEnsureFileExists(t *testing.T) {
	Convey("Given a unique file name, a user name and file permissions", t, func() {
		currentUser, err := nss.GetCurrentUser()
		So(err, ShouldBeNil)

		fname, err := file.UniqueName(os.TempDir(), "filetests")
		So(err, ShouldBeNil)

		perms := os.FileMode(0700)

		userID, _ := strconv.Atoi(currentUser.Uid)
		groupID, _ := strconv.Atoi(currentUser.Gid)

		// TODO(tmrts): Refactor the test
		Convey("It should create a new file if it doesn't exist", func() {
			err := file.EnsureExists(fname, perms, userID, groupID)
			So(err, ShouldBeNil)
			defer os.Remove(fname)

			fi, err := os.Lstat(fname)
			So(err, ShouldBeNil)

			uid := strconv.Itoa(int(fi.Sys().(*syscall.Stat_t).Uid))
			So(currentUser.Uid, ShouldEqual, uid)

			gid := strconv.Itoa(int(fi.Sys().(*syscall.Stat_t).Gid))
			So(currentUser.Gid, ShouldEqual, gid)
		})
	})
}

func TestCreateTempFileWithContent(t *testing.T) {
	Convey("Given a contents string and file permissions", t, func() {
		contents := "This is a text about nothing.\n"
		perms := os.FileMode(0600)

		Convey("It should create a temporary file with the given contents", func() {
			tmpFile, err := file.Temp(contents, perms)
			So(err, ShouldBeNil)

			tmpFile.Close()

			tmpFilename := tmpFile.Name()
			defer os.Remove(tmpFilename)

			byteContent, err := ioutil.ReadFile(tmpFilename)

			So(contents, ShouldEqual, string(byteContent))

			fi, err := os.Lstat(tmpFilename)

			So(err, ShouldBeNil)
			So(fi.Mode(), ShouldEqual, perms)
		})
	})
}
