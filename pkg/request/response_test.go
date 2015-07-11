package request_test

import (
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	. "github.com/tmrts/flamingo/pkg/util/testutil"

	"github.com/tmrts/flamingo/pkg/request"
	"github.com/tmrts/flamingo/pkg/util/mock"
)

func TestTextToJSON(t *testing.T) {
	Convey("Given a RESTful response", t, func() {
		mockBody := mock.NewReadCloser("{\"response\": [1,2,3]}")

		r := &request.Response{
			&http.Response{Body: mockBody},
		}

		Convey("It should convert it into given JSON struct", func() {
			var x struct {
				Response []int
			}

			err := r.JSON(&x)
			So(err, ShouldBeNil)

			So(x.Response, ShouldConsistOf, 1, 2, 3)
			So(mockBody.IsClosed, ShouldBeTrue)
		})
	})
}
