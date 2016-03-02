package retry_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/nproc/errorgroup-go"
	"github.com/nproc/retry-go"
	. "github.com/smartystreets/goconvey/convey"
)

func TestRetry(t *testing.T) {
	Convey("Retry", t, func() {
		singleZeroInterval := []time.Duration{0}
		tryFnSuccess := func(int, int) error { return nil }
		makeTryFnFail := func(message string) func(int, int) error {
			return func(int, int) error {
				return errors.New(message)
			}
		}
		makeTryFnPanic := func(message string) func(int, int) error {
			return func(int, int) error {
				panic(errors.New(message))
				return nil
			}
		}

		Convey("WithBackoffArray", func() {
			Convey("it should return nil if TryFn returns nil", func() {
				err := retry.WithBackoffArray(singleZeroInterval, tryFnSuccess)
				So(err, ShouldBeNil)
			})

			Convey("it should return the errors returned by TryFn in a ErrorGroup", func() {
				err := retry.WithBackoffArray(nil, makeTryFnFail("Error A"))
				So(err, ShouldHaveSameTypeAs, errorgroup.New(nil))
				errGrp := err.(*errorgroup.ErrorGroup)
				So(errGrp.Errors, ShouldHaveLength, 1)
				So(errGrp.Errors[0].Error(), ShouldEqual, "Error A")
			})

			Convey("it should recover from panic", func() {
				err := retry.WithBackoffArray(nil, makeTryFnPanic("Panic A"))
				So(err, ShouldHaveSameTypeAs, errorgroup.New(nil))
				errGrp := err.(*errorgroup.ErrorGroup)
				So(errGrp.Errors, ShouldHaveLength, 1)
				So(errGrp.Errors[0].Error(), ShouldEqual, "Panic A")
			})

			Convey("it should return an ErrorGroup with all errors", func() {
				err := retry.WithBackoffArray(singleZeroInterval, func(a, l int) error {
					return fmt.Errorf("Error %d", a)
				})
				So(err, ShouldHaveSameTypeAs, errorgroup.New(nil))
				errGrp := err.(*errorgroup.ErrorGroup)
				So(errGrp.Errors, ShouldHaveLength, 2)
				So(errGrp.Errors[0].Error(), ShouldEqual, "Error 0")
				So(errGrp.Errors[1].Error(), ShouldEqual, "Error 1")
			})

			Convey("it should wait based on the intervals in the BackoffArray", func() {
				interval := 10 * time.Millisecond
				backoff := []time.Duration{interval, interval, interval, interval}
				start := time.Now()
				err := retry.WithBackoffArray(backoff, func(a, l int) error {
					if a != l {
						return fmt.Errorf("Error %d", a)
					}
					return nil
				})
				end := time.Now()
				So(err, ShouldBeNil)
				diff := end.Sub(start)
				So(diff, ShouldBeGreaterThanOrEqualTo, 40*time.Millisecond)
				So(diff, ShouldBeLessThan, 50*time.Millisecond)
			})

			Convey("it should return an error if TryFn is nil", func() {
				err := retry.WithBackoffArray(singleZeroInterval, nil)
				So(err, ShouldEqual, retry.ErrTryFnNil)
			})
		})

		Convey("WithFixedInterval", func() {
			Convey("it should return nil if TryFn returns nil", func() {
				err := retry.WithFixedInterval(0, 1, tryFnSuccess)
				So(err, ShouldBeNil)
			})

			Convey("it should return the errors returned by TryFn in a ErrorGroup", func() {
				err := retry.WithFixedInterval(0, 0, makeTryFnFail("Error A"))
				So(err, ShouldHaveSameTypeAs, errorgroup.New(nil))
				errGrp := err.(*errorgroup.ErrorGroup)
				So(errGrp.Errors, ShouldHaveLength, 1)
				So(errGrp.Errors[0].Error(), ShouldEqual, "Error A")
			})

			Convey("it should recover from panic", func() {
				err := retry.WithFixedInterval(0, 0, makeTryFnPanic("Panic A"))
				So(err, ShouldHaveSameTypeAs, errorgroup.New(nil))
				errGrp := err.(*errorgroup.ErrorGroup)
				So(errGrp.Errors, ShouldHaveLength, 1)
				So(errGrp.Errors[0].Error(), ShouldEqual, "Panic A")
			})

			Convey("it should return an ErrorGroup with all errors", func() {
				err := retry.WithFixedInterval(0, 1, func(a, l int) error {
					return fmt.Errorf("Error %d", a)
				})
				So(err, ShouldHaveSameTypeAs, errorgroup.New(nil))
				errGrp := err.(*errorgroup.ErrorGroup)
				So(errGrp.Errors, ShouldHaveLength, 2)
				So(errGrp.Errors[0].Error(), ShouldEqual, "Error 0")
				So(errGrp.Errors[1].Error(), ShouldEqual, "Error 1")
			})

			Convey("it should wait based on the intervals in the BackoffArray", func() {
				start := time.Now()
				err := retry.WithFixedInterval(10*time.Millisecond, 4, func(a, l int) error {
					if a != l {
						return fmt.Errorf("Error %d", a)
					}
					return nil
				})
				end := time.Now()
				So(err, ShouldBeNil)
				diff := end.Sub(start)
				So(diff, ShouldBeGreaterThanOrEqualTo, 40*time.Millisecond)
				So(diff, ShouldBeLessThan, 50*time.Millisecond)
			})

			Convey("it should return an error if TryFn is nil", func() {
				err := retry.WithFixedInterval(1, 1, nil)
				So(err, ShouldEqual, retry.ErrTryFnNil)
			})
		})
	})
}
