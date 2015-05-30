package cli

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

type testCliInject struct {
	Foo int
}

func TestCliRun(t *testing.T) {
	Convey("Run cli command", t, func() {
		called := 0
		var handledErr error
		c := New("foo", "1.0.0", "").
			New("bar", "", func(c *Cli, o *Command) error {
			called = 1
			return nil
		}).
			New("zoing", "", func(x *testCliInject) error {
			called = x.Foo
			return nil
		}).
			SetErrorHandler(func(err error) {
			handledErr = err
		}).
			Register(&testCliInject{
			Foo: 100,
		})

		Convey("Run existing method", func() {
			c.RunWith([]string{"bar"})
			So(handledErr, ShouldBeNil)
			So(called, ShouldEqual, 1)
		})

		Convey("Run existing method with injection", func() {
			c.RunWith([]string{"zoing"})
			So(handledErr, ShouldBeNil)
			So(called, ShouldEqual, 100)
		})

		Convey("Run not existing method", func() {
			c.RunWith([]string{"baz"})
			So(handledErr, ShouldResemble, fmt.Errorf("Command \"baz\" unknown"))
		})
	})
}
