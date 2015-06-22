package sys_test

import (
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/tmrts/flamingo/pkg/context"
	"github.com/tmrts/flamingo/pkg/sys"
)

func TestScriptValidation(t *testing.T) {
	Convey("Using a temporary file", t, func() {
		scriptPermissions := os.FileMode(0600)

		isScript := make(chan bool, 1)
		performScriptValidation := func(script *os.File) error {
			script.Close()

			hasShabang, err := sys.FileHasShabang(script.Name())

			isScript <- hasShabang

			return err
		}

		Convey("Given a file without Shabang", func() {
			tmpFile := &context.TempFile{
				Permissions: scriptPermissions,
				Content:     "Some stuff\n",
			}

			Convey("It should return false", func() {
				err := <-context.Using(tmpFile, performScriptValidation)

				So(err, ShouldBeNil)
				So(<-isScript, ShouldBeFalse)
			})
		})

		Convey("Given a file with malformed Shabang", func() {
			tmpFile := &context.TempFile{
				Permissions: scriptPermissions,
				Content:     "# ! /bin/bash\nSome command\n",
			}

			Convey("It should return false", func() {
				err := <-context.Using(tmpFile, performScriptValidation)

				So(err, ShouldBeNil)
				So(<-isScript, ShouldBeFalse)
			})
		})

		Convey("Given an empty file with Shabang", func() {
			tmpFile := &context.TempFile{
				Permissions: scriptPermissions,
				Content:     "#! /bin/bash\n\n",
			}

			Convey("It should return true", func() {
				err := <-context.Using(tmpFile, performScriptValidation)

				So(err, ShouldBeNil)
				So(<-isScript, ShouldBeTrue)
			})
		})

		Convey("Given a script file with Shabang", func() {
			tmpFile := &context.TempFile{
				Permissions: scriptPermissions,
				Content:     "#! /bin/bash\nprintf ${PWD}\n",
			}

			Convey("It should return true", func() {
				err := <-context.Using(tmpFile, performScriptValidation)

				So(err, ShouldBeNil)
				So(<-isScript, ShouldBeTrue)
			})
		})
	})
}

func TestScriptExecution(t *testing.T) {
	Convey("Using a temporary file", t, func() {
		scriptPermissions := os.FileMode(0700)

		outChan := make(chan string, 1)
		executeScript := func(script *os.File) error {
			script.Close()

			out, err := sys.DefaultExecutor.Execute(script.Name())

			outChan <- out

			return err
		}

		Convey("Given a script file", func() {
			scriptContent := "#! /usr/bin/env bash\n"
			scriptContent += "echo -n Yeehaw\n"
			scriptContent += "\n"

			tmpFile := &context.TempFile{
				Permissions: scriptPermissions,
				Content:     scriptContent,
			}

			Convey("It should execute the script and return its output", func() {
				err := <-context.Using(tmpFile, executeScript)

				So(err, ShouldBeNil)
				So(<-outChan, ShouldEqual, "Yeehaw")
			})
		})

	})
}
