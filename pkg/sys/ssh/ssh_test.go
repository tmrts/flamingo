package ssh_test

import (
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/tmrts/flamingo/pkg/context"
	"github.com/tmrts/flamingo/pkg/sys"
	"github.com/tmrts/flamingo/pkg/sys/nss"
	"github.com/tmrts/flamingo/pkg/sys/ssh"
)

func TestSSHKeyVerification(t *testing.T) {
	Convey("Given an ssh key", t, func() {
		Convey("It should validate the given key", func() {
			Convey("For an empty SSH key", func() {
				key := []byte("")

				Convey("It should return false", func() {
					err := ssh.Verify(key)
					So(err, ShouldNotBeNil)
				})
			})

			Convey("For a malformed SSH key", func() {
				key := []byte("GIBBERISH")

				Convey("It should return false", func() {
					err := ssh.Verify(key)
					So(err, ShouldNotBeNil)
				})
			})

			Convey("For a valid rsa SSH public key", func() {
				key, err := ioutil.ReadFile("test_keys/test_rsa.pub")
				So(err, ShouldBeNil)

				byteKey := []byte(key)

				Convey("It should return true", func() {
					err := ssh.Verify(byteKey)
					So(err, ShouldBeNil)
				})
			})

			Convey("For a valid rsa1 SSH public key", func() {
				key, err := ioutil.ReadFile("test_keys/test_rsa.pub")
				So(err, ShouldBeNil)

				byteKey := []byte(key)

				Convey("It should return true", func() {
					err := ssh.Verify(byteKey)
					So(err, ShouldBeNil)
				})
			})

			Convey("For a valid dsa SSH public key", func() {
				key, err := ioutil.ReadFile("test_keys/test_rsa.pub")
				So(err, ShouldBeNil)

				byteKey := []byte(key)

				Convey("It should return true", func() {
					err := ssh.Verify(byteKey)
					So(err, ShouldBeNil)
				})
			})

			Convey("For a valid ecdsa SSH public key", func() {
				key, err := ioutil.ReadFile("test_keys/test_rsa.pub")
				So(err, ShouldBeNil)

				byteKey := []byte(key)

				Convey("It should return true", func() {
					err := ssh.Verify(byteKey)
					So(err, ShouldBeNil)
				})
			})
		})
	})
}

func TestSSHDirectoryStructureInitialization(t *testing.T) {
	Convey("Given a user", t, func() {
		fakeHomeDir, err := ioutil.TempDir("", "flamingo-fakeuser")
		So(err, ShouldBeNil)

		defer os.RemoveAll(fakeHomeDir)

		curUser, err := nss.CurrentUser()
		So(err, ShouldBeNil)

		fakeUser := &user.User{
			Uid:      curUser.Uid,
			Gid:      curUser.Gid,
			Username: curUser.Username,
			Name:     "FakeUser",
			HomeDir:  fakeHomeDir,
		}

		Convey("It should initialize the ssh directory for that user", func() {
			err := ssh.InitializeFor(fakeUser)
			So(err, ShouldBeNil)

			_, err = os.Open(filepath.Join(fakeHomeDir, ".ssh"))
			So(err, ShouldBeNil)

			_, err = os.Open(filepath.Join(fakeHomeDir, ".ssh/authorized_keys"))
			So(err, ShouldBeNil)
		})
	})
}

func TestSSHKeyAuthorization(t *testing.T) {
	Convey("Given an SSH authorized keys file and public keys", t, func() {
		rsa_key, err := ioutil.ReadFile("test_keys/test_rsa.pub")
		So(err, ShouldBeNil)

		dsa_key, err := ioutil.ReadFile("test_keys/test_dsa.pub")
		So(err, ShouldBeNil)

		tmpAuthKeysFile := &context.TempFile{
			Content:     "",
			Permissions: os.FileMode(0600),
		}

		Convey("It should append the SSH keys to the authorized keys file", func() {
			authorizeKey := func(f *os.File) error {
				err := ssh.AuthorizeKeys(f, []byte(rsa_key), []byte(dsa_key))
				if err != nil {
					return err
				}
				f.Close()

				_, err = sys.DefaultExecutor.Execute("ssh-keygen", "-l", "-f", f.Name())

				return err
			}

			err := <-context.Using(tmpAuthKeysFile, authorizeKey)
			So(err, ShouldBeNil)
		})
	})
}
