package rest_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRESTfulClient(t *testing.T) {
	Convey("Given a URL and a REST client", t, func() {
		//url := "http://localhost"

		/*
		 *client := &rest.Client{
		 *    http.Client{
		 *        RoundTripper: func(rq *http.Request) (*http.Response, error) {
		 *            mockResponse := io.ReadCloser{}
		 *            return mockResponse, nil
		 *        },
		 *    },
		 *}
		 */

		//client := rest.DefaultClient
		Convey("When the client requests the contents", func() {
			//response, err := client.Get(url)
			//So(err, ShouldBeNil)

			Convey("Then the response should be the raw contents of the response", func() {
				//json, err := rest.JSON(response)
				//So(err, ShouldBeNil)

				//So(json, ShouldEqual, "http://httpbin.org/get")
			})
		})
	})
}
