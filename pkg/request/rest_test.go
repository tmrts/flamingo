package request_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/tmrts/flamingo/pkg/request"
)

type mockHandler struct{}

func (m mockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	buf, _ := ioutil.ReadFile("test_data/some.json")
	w.Write(buf)
}

func TestRESTfulClient(t *testing.T) {
	Convey("Given a URL and a REST client", t, func() {
		server := httptest.NewServer(mockHandler{})

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
