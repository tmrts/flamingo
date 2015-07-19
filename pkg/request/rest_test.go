package request_test

import (
	"io/ioutil"
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/tmrts/flamingo/pkg/request"
	"github.com/tmrts/flamingo/pkg/util/testutil"
)

func TestRESTfulClient(t *testing.T) {
	Convey("Given a URL and a REST client", t, func() {
		server := testutil.NewMockServer(func(w http.ResponseWriter, r *http.Request) {
			buf, _ := ioutil.ReadFile("test_data/some.json")
			w.Write(buf)
		})

		Convey("When the client requests the contents", func() {
			response, err := request.Get(server.URL)
			So(err, ShouldBeNil)

			Convey("Then the response should be the raw contents of the response", func() {
				var data map[string]interface{}

				err := response.JSON(&data)
				So(err, ShouldBeNil)

				So(data["isJSON"], ShouldBeTrue)
			})
		})
	})
}
