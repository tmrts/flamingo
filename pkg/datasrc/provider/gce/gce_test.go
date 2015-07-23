package gce_test

import (
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/tmrts/flamingo/pkg/sys/ssh"
	. "github.com/tmrts/flamingo/pkg/util/testutil"

	"github.com/tmrts/flamingo/pkg/datasrc/provider"
	"github.com/tmrts/flamingo/pkg/datasrc/provider/gce"
)

const (
	testMetadataDir = "test_metadata"
)

func TestGoogleComputeMetadataRetrieval(t *testing.T) {
	Convey("Given a GCE meta-data service", t, func() {
		// mock GCE metadata server
		server := NewMockServer(func(w http.ResponseWriter, r *http.Request) {
			var json_path string
			if r.Header.Get("Metadata-Flavor") != "Google" {
				http.Error(w, "metadata header is not found", http.StatusBadRequest)
				return
			}

			if strings.Contains(r.URL.String(), "project") {
				json_path = filepath.Join(testMetadataDir, "GCEv1_project.json")
			} else if strings.Contains(r.URL.String(), "instance") {
				json_path = filepath.Join(testMetadataDir, "GCEv1_instance.json")
			} else {
				http.Error(w, "requested resource is not found", http.StatusNotFound)
				return
			}

			buf, err := ioutil.ReadFile(json_path)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			w.Write(buf)
		})

		service := gce.MetadataService{
			URL: provider.FormatURL(server.URL + "/%v/%v"),
		}

		Convey("It should retrieve meta-data from GCE meta-data service", func() {
			digest, err := service.FetchMetadata()
			So(err, ShouldBeNil)

			So(digest.Hostname, ShouldEqual, "centos.internal")

			ifc := digest.PrimaryNetworkInterface()
			So(ifc.PrivateIP.String(), ShouldEqual, "10.240.45.128")

			So(ifc.PublicIPs[0].String(), ShouldEqual, "104.155.21.159")
			So(ifc.PublicIPs[1].String(), ShouldEqual, "104.155.21.160")

			sshKeys := digest.SSHKeys

			So(sshKeys["user1"], ShouldConsistOf,
				ssh.Key("ssh-rsa RSA_PUBLIC_KEY_FOR_USER_1 user1@machine"),
				ssh.Key("ssh-dsa DSA_PUBLIC_KEY_FOR_USER_1 user1@machine"))

			So(sshKeys["user2"], ShouldConsistOf,
				ssh.Key("ssh-rsa RSA_PUBLIC_KEY_FOR_USER_2 user2@machine"))
		})

		Convey("It should retrieve user-data from GCE meta-data service", func() {
			userdata, err := service.FetchUserdata()
			So(err, ShouldBeNil)

			So(userdata["user-data"], ShouldEqual, "#cloud-config\n")
		})

	})
}
