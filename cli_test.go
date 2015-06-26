package clif

import (
	"bytes"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"io"
	"testing"
)

type testCliAlias interface {
	Hello() int
}

type testCliInject struct {
	Foo int
}

func (this *testCliInject) Hello() int {
	return this.Foo
}

func TestCliRun(t *testing.T) {
	Convey("Run cli command", t, func() {
		called := 0
		var handledErr error
		Die = func(msg string, args ...interface{}) {
			panic(fmt.Sprintf(msg, args...))
		}

		c := New("foo", "1.0.0", "").
			New("bar", "", func(c *Cli, o *Command) error {
			called = 1
			return nil
		}).
			New("zoing", "", func(x *testCliInject) error {
			called = x.Foo
			return nil
		}).
			New("zoing2", "", func(x testCliAlias) error {
			called = x.Hello()
			return nil
		}).
			New("oops", "", func(x io.Writer) error {
			panic("Should never be called")
			return nil
		}).
			New("errme", "", func() error {
			return fmt.Errorf("I error!")
		}).
			Register(&testCliInject{
			Foo: 100,
		}).
			RegisterAs("clif.testCliAlias", &testCliInject{
			Foo: 200,
		})

		cmdInvalid := NewCommand("bla", "Dont use me", func() {})
		argInvalid := NewArgument("something", "..", "", false, false)
		argInvalid.SetValidator(func(name, value string) error {
			return fmt.Errorf("Never works!")
		})
		cmdInvalid.AddArgument(argInvalid)
		c.Add(cmdInvalid)

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
		Convey("Run existing method with interface injection", func() {
			c.RunWith([]string{"zoing2"})
			So(handledErr, ShouldBeNil)
			So(called, ShouldEqual, 200)
		})
		Convey("Run not existing method", func() {
			So(func() {
				c.RunWith([]string{"baz"})
			}, ShouldPanicWith, "Command \"baz\" unknown")
		})
		Convey("Run without args describes and exits", func() {
			buf := bytes.NewBuffer(nil)
			out := NewOutput(buf, NewDefaultFormatter(map[string]string{}))
			c.SetOutput(out)
			c.RunWith([]string{})
			So(buf.String(), ShouldEqual, DescribeCli(c))
		})
		Convey("Run method with not registered arg fails", func() {
			So(func() {
				c.RunWith([]string{"oops"})
			}, ShouldPanicWith, "Missing callback parameter io.Writer")
		})
		Convey("Run method with invalid arg fails", func() {
			So(func() {
				c.RunWith([]string{"bla", "bla"})
			}, ShouldPanicWith, "Parse error: Parameter \"something\" invalid: Never works!")
		})
		Convey("Run method with resulting error returns it", func() {
			So(func() {
				c.RunWith([]string{"errme"})
			}, ShouldPanicWith, "Failure in execution: I error!")
		})
	})
}

func TestCliConstruction(t *testing.T) {
	Convey("Create new Cli with commands", t, func() {
		app := New("My App", "1.0.0", "Testing app")
		cb := func() {}

		Convey("Use command constructor", func() {
			app.New("foo", "For fooing", cb)
			So(len(app.Commands), ShouldEqual, 2)
			So(app.Commands["foo"], ShouldNotBeNil)
		})

		Convey("Use variadic adding", func() {
			app.New("foo", "For fooing", cb)
			cmds := []*Command{
				NewCommand("foo", "For fooing", cb),
				NewCommand("bar", "For baring", cb),
			}
			app.Add(cmds...)
			So(len(app.Commands), ShouldEqual, 3)
			So(app.Commands["foo"], ShouldNotBeNil)
			So(app.Commands["bar"], ShouldNotBeNil)
		})
	})
}
