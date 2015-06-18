package context_test

import (
	"errors"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/tmrts/flamingo/pkg/context"
)

type Signal int

const (
	Enter Signal = 1 + iota
	Exit
	Use
)

type RandomNumberContext struct {
	EventChannel chan Signal

	seed  int
	value int
}

// Enter prepares the resources to be used
func (rnc *RandomNumberContext) Enter() error {
	rnc.value = rnc.seed

	rnc.EventChannel <- Enter

	return nil
}

// Exit cleans up after the context is used
func (rnc *RandomNumberContext) Exit() error {
	rnc.value = 0

	rnc.EventChannel <- Exit
	return nil
}

// Use casts the given function to its own context function and calls it
func (rnc *RandomNumberContext) Use(fn interface{}) error {
	closure := fn.(func(int) error)

	rnc.EventChannel <- Use
	return closure(rnc.value)
}

func TestUsingContexts(t *testing.T) {
	Convey("Given a context and a context use function", t, func() {
		NotEvenErr := errors.New("not a good number")

		checkIfEven := func(no int) error {
			if no%2 == 1 {
				return NotEvenErr
			}
			return nil
		}

		Convey("It should use the context in the order, Enter -> Use -> Exit", func() {
			ch := make(chan Signal)
			randomNumber := &RandomNumberContext{
				seed:         2,
				EventChannel: ch,
			}

			errch := context.Using(randomNumber, checkIfEven)

			event1, event2, event3 := <-ch, <-ch, <-ch

			So(event1, ShouldEqual, Enter)
			So(event2, ShouldEqual, Use)
			So(event3, ShouldEqual, Exit)

			err := <-errch
			So(err, ShouldBeNil)
		})

		Convey("It should propagate the errors encountered during use", func() {
			ch := make(chan Signal, 3)
			randomNumber := &RandomNumberContext{
				seed:         3,
				EventChannel: ch,
			}
			errch := context.Using(randomNumber, checkIfEven)

			err := <-errch
			So(err, ShouldEqual, NotEvenErr)
		})
	})
}
