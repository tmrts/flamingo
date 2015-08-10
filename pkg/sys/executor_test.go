package sys_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/tmrts/flamingo/pkg/sys"
)

func TestExecutesACommand(t *testing.T) {
	Convey("Given a command and a list of arguments", t, func() {
		cmd, args := "echo", []string{"-n", "fi fye fo fum"}

		Convey("It should execute the command and return its output", func() {
			out, err := sys.DefaultExecutor.Execute(cmd, args...)
			So(err, ShouldBeNil)

			So(out, ShouldEqual, "fi fye fo fum")
		})
	})
}

func TestExecutesAnInvalidCommand(t *testing.T) {
	Convey("Given an erroneous command", t, func() {
		cmd := "tHiScOmMaNdDoEsNtExIsT"

		Convey("It should return the stderr as a new error", func() {
			_, err := sys.DefaultExecutor.Execute(cmd)
			So(err, ShouldNotBeNil)
		})
	})
}
