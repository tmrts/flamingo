package cloudconfig_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/tmrts/flamingo/pkg/sys/identity"
	"github.com/tmrts/flamingo/pkg/sys/ssh"
	. "github.com/tmrts/flamingo/pkg/util/testutil"

	"github.com/tmrts/flamingo/pkg/datasrc/userdata/cloudconfig"
)

const (
	testConfigDirectory = "test_configs"
)

func TestRunCmdParsing(t *testing.T) {
	Convey("Given a cloud-config file with runcmd directive", t, func() {
		configFile, err := os.Open(filepath.Join(testConfigDirectory, "runcmd.yaml"))
		So(err, ShouldBeNil)

		Convey("It should parse the execute statements", func() {
			conf, err := cloudconfig.Parse(configFile)
			So(err, ShouldBeNil)

			So(conf.Commands[0], ShouldConsistOf, "ls", "-l", "/")
			So(conf.Commands[1], ShouldConsistOf, "sh", "-xc", "echo $(date) ': hello world!'")
			So(conf.Commands[2], ShouldConsistOf, "sh", "-c", "echo \"=========hello world'=========\"")
			So(conf.Commands[3], ShouldConsistOf, "sh", "-c", "ls -l /root")
			So(conf.Commands[4], ShouldConsistOf, "wget", "http://slashdot.org", "-O", "/tmp/index.html")
		})
	})
}

func TestSSHKeyParsing(t *testing.T) {
	Convey("Given a cloud-config file containing SSH-key directives", t, func() {
		configFile, err := os.Open(filepath.Join(testConfigDirectory, "sshkeys.yaml"))
		So(err, ShouldBeNil)

		Convey("It should parse the ssh keys", func() {
			conf, err := cloudconfig.Parse(configFile)
			So(err, ShouldBeNil)

			So(conf.AuthorizedKeys["root"], ShouldConsistOf,
				ssh.Key("ssh-rsa RSA_PUBLIC_KEY_1 mykey@host"),
				ssh.Key("ssh-rsa RSA_PUBLIC_KEY_2 mykey@host"),
			)

			So(conf.SSHKeyPairs, ShouldConsistOf,
				ssh.KeyPair{
					Public:  ssh.Key("ssh-rsa RSA_PUBLIC_KEY smoser@localhost"),
					Private: ssh.Key("-----BEGIN RSA PRIVATE KEY-----\nRSA_PRIVATE_KEY\n-----END RSA PRIVATE KEY-----\n"),
				},
			)
		})
	})
}

func TestUserGroupParsing(t *testing.T) {
	Convey("Given a cloud-config file containing users and/or groups", t, func() {
		configFile, err := os.Open(filepath.Join(testConfigDirectory, "usergroups.yaml"))
		So(err, ShouldBeNil)

		Convey("It should parse the users and groups", func() {
			c, err := cloudconfig.Parse(configFile)

			So(err, ShouldBeNil)

			So(c.Groups["cloud-users"], ShouldBeEmpty)

			So(c.Users["foobar"], ShouldStructEqual, identity.User{
				Name:            "foobar",
				GECOS:           "Foo B. Bar",
				GID:             "foobar",
				SecondaryGroups: []string{"users"},
				PasswordHash:    "$6$SHA256$PASSWORD_HASH",
				ExpireDate:      "2012-09-01",
				SELinuxUser:     "staff_u",
			})

			So(c.Users["barfoo"], ShouldStructEqual, identity.User{
				Name:            "barfoo",
				GECOS:           "Bar B. Foo",
				SecondaryGroups: []string{"users", "admin"},
			})

			So(c.Users["daemon"], ShouldStructEqual, identity.User{
				Name:            "daemon",
				GECOS:           "Magic Cloud App Daemon User",
				IsInactive:      true,
				IsSystemAccount: true,
			})
		})
	})
}

func TestWriteFileDirective(t *testing.T) {
	Convey("Given a configuration file containing files to be written", t, func() {
		configFile, err := os.Open(filepath.Join(testConfigDirectory, "write_files.yaml"))
		So(err, ShouldBeNil)

		Convey("It should parse the file names and the file contents", func() {
			conf, err := cloudconfig.Parse(configFile)

			So(err, ShouldBeNil)

			So(conf.Files[0], ShouldStructEqual, cloudconfig.WriteFile{
				Path:        "/etc/sysconfig/selinux",
				Permissions: "0644",
				Owner:       "some_user:some_group",
				Encoding:    "b64",
				Content:     "STRING_FILE_CONTENT",
			})

			So(conf.Files[1], ShouldStructEqual, cloudconfig.WriteFile{
				Path:    "/etc/sysconfig/samba",
				Content: "# My new /etc/sysconfig/samba file\n\nSMBDOPTIONS=\"-D\"\n",
			})

			/* TODO(tmrts): Test binary file contents
			 *
			 *So(conf.Files[2], ShouldStructEqual, cloudconfig.WriteFile{
			 *    Path:        "/bin/arch",
			 *    Permissions: "0555",
			 *    Content:     base64.StdEncoding.EncodeToString([]byte("f0VMRgIBAQAAAAAAAAAAAAIAPgABAAAAwARAAAAAAABAAAAAAAAAAJAVAAAAAAAAAAAAAEAAOAAI")),
			 *})
			 */
		})
	})
}

func TestCloudConfigCheck(t *testing.T) {
	Convey("Given a cloud-config file", t, func() {
		content := strings.NewReader("#cloud-config\nkey:\nvalue\n")

		mockContent := ioutil.NopCloser(content)

		Convey("It should return no errors", func() {
			err := cloudconfig.IsValid(mockContent)
			So(err, ShouldBeNil)
		})
	})

	Convey("Given an invalid cloud-config", t, func() {
		content := strings.NewReader("#! /usr/bin/env bash\nls .\n")

		mockContent := ioutil.NopCloser(content)

		Convey("It should return invalid cloud-config error", func() {
			err := cloudconfig.IsValid(mockContent)

			So(err, ShouldEqual, cloudconfig.ErrNotACloudConfigFile)
		})
	})
}
