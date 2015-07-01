package nss_test

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/tmrts/flamingo/pkg/sys/nss"
)

type fakeExecutor struct {
	expectedCmd    string
	expectedResult string
}

func (fe *fakeExecutor) Execute(e string, args ...string) (string, error) {
	if e != fe.expectedCmd {
		return "", fmt.Errorf("Execute(%q) -> expected: %q, got: %q", e, e, fe.expectedCmd)
	}

	return fe.expectedResult, nil
}

func TestShadowDatabaseEntryRetrieval(t *testing.T) {
	Convey("Given an nss server and a user name", t, func() {
		uname := "root"

		server := &nss.Server{
			Exec: &fakeExecutor{
				expectedCmd:    "getent",
				expectedResult: "root:$4$SALT$PASSWORD_HASH:16590:0:99999:7::::",
			},
		}

		Convey("It should query the NSS shadow database and get its entry", func() {
			root, err := nss.GetShadowEntry(server, uname)
			So(err, ShouldBeNil)

			So(root.UserName, ShouldEqual, "root")
			So(root.PasswordHash, ShouldEqual, "$4$SALT$PASSWORD_HASH")
		})
	})
}
