package datasrc_test

import (
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/tmrts/flamingo/pkg/datasrc"
	"github.com/tmrts/flamingo/pkg/datasrc/provider"
	"github.com/tmrts/flamingo/pkg/datasrc/provider/gce"
	"github.com/tmrts/flamingo/pkg/util/testutil"
)

func TestFetchMetadata(t *testing.T) {
	Convey("Given a list of datasources and a timeout duration", t, func() {
		mockGCEServer := testutil.NewMockServer(func(w http.ResponseWriter, r *http.Request) {
			var json_path string
			if r.Header.Get("Metadata-Flavor") != "Google" {
				http.Error(w, "metadata header is not found", http.StatusBadRequest)
				return
			}

			testMetadataDir := "provider/gce/test_metadata"
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

		mockGCEProvider := &gce.MetadataService{
			URL: provider.FormatURL(mockGCEServer.URL + "/%v/%v"),
		}

		mockDataSources := map[string]datasrc.Provider{
			"gce": mockGCEProvider,
		}

		timeout := time.Millisecond * 500

		Convey("It should find the available data source provider", func() {
			provider, err := datasrc.FindProvider(mockDataSources, timeout)
			So(err, ShouldBeNil)

			So(provider, ShouldEqual, mockGCEProvider)
		})

		Convey("When datasources are unavailable it should timeout", func() {
			_, err := datasrc.FindProvider(map[string]datasrc.Provider{}, 0*time.Second)
			So(err, ShouldEqual, datasrc.ErrDatasourceRetrievalTimeout)
		})
	})
}
