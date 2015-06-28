package injection_test

import (
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/tmrts/flamingo/pkg/context"
	"github.com/tmrts/flamingo/pkg/injection"
)

type toLowerCase func(string) string

func (d *toLowerCase) Inject(fn interface{}) injection.Restore {
	maskedFn := *d
	*d = fn.(func(string) string)

	return func() {
		*d = maskedFn
	}
}

var ToLower toLowerCase = func(s string) string {
	return strings.ToLower(s)
}

func TestDependencyInjection(t *testing.T) {
	Convey("Given a dependency and a replacement for that dependency with the same type", t, func() {
		So(ToLower("HELLO"), ShouldEqual, "hello")

		fakeToLower := func(string) string {
			return "This function has been hijacked."
		}

		Convey("It should inject the replacement to the dependency", func() {
			restoreDependency := ToLower.Inject(fakeToLower)

			So(ToLower("HELLO"), ShouldEqual, "This function has been hijacked.")
			So(ToLower("WORLD"), ShouldEqual, "This function has been hijacked.")

			restoreDependency()

			So(ToLower("HELLO"), ShouldEqual, "hello")
		})

		Convey("Within an injection context", func() {
			mockToLower := &injection.ContextManager{
				Dependency:  &ToLower,
				Replacement: fakeToLower,
			}

			Convey("It should use the injected function within the context", func() {
				outch := make(chan string)
				context.Using(mockToLower, func() error {
					outch <- ToLower("HELLO")
					return nil
				})

				So(<-outch, ShouldEqual, "This function has been hijacked.")
			})

			Convey("It should restore the injected function after the context is exited", func() {
				So(ToLower("HELLO"), ShouldEqual, "hello")
			})
		})
	})
}
