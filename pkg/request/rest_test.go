package rest_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/tmrts/flamingo/pkg/rest"
)

func TestRESTfulClient(t *testing.T) {
	Convey("Given a URL and a REST client", t, func() {
		url := "http://httpbin.org/get"

		//client := rest.DefaultClient
		Convey("When the client requests the contents", func() {
			response, err := rest.Get(url)
			So(err, ShouldBeNil)

			Convey("Then the response should be the raw contents of the response", func() {
				json, err := response.JSON()
				So(err, ShouldBeNil)

				So(json, ShouldEqual, "http://httpbin.org/get")
			})
		})
	})
}
