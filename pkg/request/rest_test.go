package request_test

import (
	"fmt"
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/tmrts/flamingo/pkg/request"
)

type mockHandler struct{}

func (m mockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("{'isJSON': true}"))
}

func TestRESTfulClient(t *testing.T) {
	Convey("Given a URL and a REST client", t, func() {
		localServer := &http.Server{
			Handler: mockHandler{},
		}
		localServer.ListenAndServe()

		url := localServer.Addr

		fmt.Printf("url %\n", url)
		//client := rest.DefaultClient
		Convey("When the client requests the contents", func() {
			response, err := request.Get(url)
			So(err, ShouldBeNil)

			Convey("Then the response should be the raw contents of the response", func() {
				var json interface{}
				err := response.JSON(&json)
				So(err, ShouldBeNil)

				So(json.isJSON, ShouldBeTrue)
			})
		})
	})
}
