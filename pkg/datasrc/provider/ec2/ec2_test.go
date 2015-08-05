package ec2_test

import (
	"net/http"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	. "github.com/tmrts/flamingo/pkg/util/testutil"

	"github.com/tmrts/flamingo/pkg/datasrc/provider"
	"github.com/tmrts/flamingo/pkg/datasrc/provider/ec2"
	"github.com/tmrts/flamingo/pkg/sys/ssh"
)

const (
	testMetadataDir = "test_metadata"
)

func TestRetrievesDataFromEC2(t *testing.T) {
	Convey("Given an EC2 meta-data service", t, func() {
		// mock EC2 metadata server
		server := NewMockServer(func(w http.ResponseWriter, r *http.Request) {
			attributes := map[string]string{
				"hostname":                  "centos.ec2",
				"local-ipv4":                "10.240.51.29",
				"public-ipv4":               "104.155.21.99",
				"public-keys/0/openssh-key": "ssh-rsa OPENSSH_KEY",
			}

			if strings.Contains(r.URL.String(), "/2009-04-04/meta-data/") {
				for attr, value := range attributes {
					if strings.HasSuffix(r.URL.String(), attr) {
						w.Write([]byte(value))
					}
				}
			} else if strings.HasSuffix(r.URL.String(), "/2009-04-04/user-data") {
				w.Write([]byte("#cloud-config\n"))
			} else {
				http.Error(w, "requested resource is not found", http.StatusNotFound)
			}
		})

		service := ec2.MetadataService{
			URL: provider.FormatURL(server.URL + "/%v/%v/%v"),
		}

		Convey("It should retrieve meta-data from EC2 meta-data service", func() {
			digest, err := service.FetchMetadata()
			So(err, ShouldBeNil)

			So(digest.Hostname, ShouldEqual, "centos.ec2")

			ifc := digest.PrimaryNetworkInterface()
			So(ifc.PublicIPs[0].String(), ShouldEqual, "104.155.21.99")

			sshKeys := digest.SSHKeys

			So(sshKeys["root"], ShouldConsistOf, ssh.Key("ssh-rsa OPENSSH_KEY"))
		})

		Convey("It should retrieve user-data from EC2 meta-data service", func() {
			userdata, err := service.FetchUserdata()
			So(err, ShouldBeNil)

			So(userdata["user-data"], ShouldEqual, "#cloud-config\n")
		})

	})
}
