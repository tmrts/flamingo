package configdrive_test

import (
	"fmt"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/tmrts/flamingo/pkg/datasrc/provider/configdrive"
	"github.com/tmrts/flamingo/pkg/sys"
)

func TestLocatingConfigDriveMountTarget(t *testing.T) {
	Convey("Given a mounted device with the label 'config-2'", t, func() {
		executor := sys.NewFuncExecutor(func(cmd string, args ...string) (string, error) {
			cmd = strings.Join(append([]string{cmd}, args...), " ")

			switch cmd {
			case "blkid -t LABEL='config-2' -odevice":
				return "/dev/cfg", nil
			case "findmnt --raw --noheadings --output TARGET /dev/cfg":
				return "/mnt/config", nil
			}

			return "", fmt.Errorf("unrecognized command")
		})

		Convey("It should locate where the device is mounted", func() {
			target, err := configdrive.FindMountTarget(executor)
			So(err, ShouldBeNil)

			So(target, ShouldEqual, "/mnt/config")
		})
	})
}
