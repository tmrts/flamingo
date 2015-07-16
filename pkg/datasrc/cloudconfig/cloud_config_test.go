package cloudconfig_test

import (
	"path/filepath"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	. "github.com/tmrts/flamingo/pkg/util/testutil"

	"github.com/tmrts/flamingo/pkg/datasrc/cloudconfig"
)

const (
	testConfigDirectory = "test_configs"
)

func TestRunCmdParsing(t *testing.T) {
	Convey("Given a cloud-config file with runcmd directive", t, func() {
		configFile := filepath.Join(testConfigDirectory, "runcmd.yaml")

		Convey("It should parse the execute statements", func() {
			conf, err := cloudconfig.Parse(configFile)
			So(err, ShouldBeNil)

			So(conf.Commands, ShouldConsistOf,
				"ls -l /",
				"sh -xc \"echo $(date) ': hello world!'\"",
				"sh -c echo \"=========hello world'=========\"",
				"ls -l /root",
				"wget http://slashdot.org -O /tmp/index.html",
			)
		})
	})
}

func TestSSHKeyParsing(t *testing.T) {
	Convey("Given a cloud-config file containing SSH-key directives", t, func() {
		configFile := filepath.Join(testConfigDirectory, "sshkeys.yaml")

		Convey("It should parse the ssh keys", func() {
			c, err := cloudconfig.Parse(configFile)
			So(err, ShouldBeNil)

			So(c.AuthorizedKeys[0], ShouldConsistOf, []byte("ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAGEA3FSyQwBI6Z+nCSjUUk8EEAnnkhXlukKoUPND/RRClWz2s5TCzIkd3Ou5+Cyz71X0XmazM3l5WgeErvtIwQMyT1KjNoMhoJMrJnWqQPOt5Q8zWd9qG7PBl9+eiH5qV7NZ mykey@host"))

			//So(c.SSHKeyPairs[1].Public, ShouldEqual, []byte())
		})
	})
}

func TestUserGroupParsing(t *testing.T) {
	Convey("Given a cloud-config file containing users and/or groups", t, func() {
		configFile := filepath.Join(testConfigDirectory, "usergroups.yaml")

		Convey("It should parse the users and groups", func() {
			c, err := cloudconfig.Parse(configFile)

			So(err, ShouldBeNil)

			So(c.Groups, ShouldEqual, "")
		})
	})
}
