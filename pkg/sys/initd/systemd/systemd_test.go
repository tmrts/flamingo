package systemd_test

import (
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	. "github.com/tmrts/flamingo/pkg/util/testutil"

	"github.com/tmrts/flamingo/pkg/sys"
	"github.com/tmrts/flamingo/pkg/sys/initd/systemd"
)

func TestSystemdInitManager(t *testing.T) {
	Convey("Given a systemd implementation", t, func() {
		exec := sys.NewStubExecutor("", nil)
		sysd := &systemd.Implementation{
			UnitDir: os.TempDir(),
			Exec:    exec,
		}

		Convey("Given a systemd unit", func() {
			testUnit := systemd.NewUnit("testUnit", "here")

			Convey("It should enable the component", func() {
				sysd.Install(testUnit)

				So(<-exec.Exec, ShouldEqual, "systemctl")
				So(<-exec.Args, ShouldConsistOf, "enable", "--system", testUnit.Name())
			})

			Convey("It should disable the component", func() {
				sysd.Disable(testUnit)

				So(<-exec.Exec, ShouldEqual, "systemctl")
				So(<-exec.Args, ShouldConsistOf, "disable", testUnit.Name())
			})

			Convey("It should start the component", func() {
				sysd.Start(testUnit)

				So(<-exec.Exec, ShouldEqual, "systemctl")
				So(<-exec.Args, ShouldConsistOf, "start", testUnit.Name())
			})

			Convey("It should stop the component", func() {
				sysd.Stop(testUnit)

				So(<-exec.Exec, ShouldEqual, "systemctl")
				So(<-exec.Args, ShouldConsistOf, "stop", testUnit.Name())
			})
		})
	})
}
