package gce_test

import (
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	. "github.com/tmrts/flamingo/pkg/util/testutil"

	"github.com/tmrts/flamingo/pkg/datasrc/metadata"
	"github.com/tmrts/flamingo/pkg/datasrc/metadata/gce"
)

const (
	testMetadataDir = "test_metadata"
)

func TestGoogleComputeMetadataRetrieval(t *testing.T) {
	Convey("Given a REST client", t, func() {
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
			URL: metadata.FormatURL(server.URL + "/%v/%v"),
		}

		b, err := ioutil.ReadFile(filepath.Join(testMetadataDir, "GCEv1_instance.json"))
		So(err, ShouldBeNil)
		So(b, ShouldNotBeEmpty)

		Convey("It should retrieve v1 metadata from Google Compute Engine metadata service", func() {
			metadata, err := service.Metadata()
			So(err, ShouldBeNil)

			digest := metadata.Digest()
			So(digest.Hostname, ShouldEqual, "centos.internal")

			ifc := digest.PrimaryNetworkInterface()
			So(ifc.PrivateIP.String(), ShouldEqual, "10.240.45.128")

			So(ifc.PublicIPs[0].String(), ShouldEqual, "104.155.21.159")
			So(ifc.PublicIPs[1].String(), ShouldEqual, "104.155.21.160")
		})
	})
}
