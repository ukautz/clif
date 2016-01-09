package clif

import (
	"bytes"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"io"
	"reflect"
	"strings"
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

func TestCliCall(t *testing.T) {
	Convey("Call cli command", t, func() {
		val := []string{}
		cmd := NewCommand("foo", "Do foo", func() {
			val = append(val, "called")
		})
		cli := New("cli", "0.1.0", "The cli")
		vals, err := cli.Call(cmd)
		So(err, ShouldBeNil)
		So(len(vals), ShouldEqual, 0)
		So(strings.Join(val, " ** "), ShouldEqual, "called")

		Convey("With return values", func() {
			val = []string{}
			cmd.Call = reflect.ValueOf(func() (int, string, error) {
				val = append(val, "called")
				return 10, "foo", nil
			})
			vals, err = cli.Call(cmd)
			So(len(vals), ShouldEqual, 2)
			So(err, ShouldBeNil)
			So(strings.Join(val, " ** "), ShouldEqual, "called")
			So(vals[0].Interface(), ShouldResemble, 10)
			So(vals[1].Interface(), ShouldResemble, "foo")

			Convey("With error and values", func() {
				val = []string{}
				cmd.Call = reflect.ValueOf(func() (int, string, error) {
					val = append(val, "called")
					return 10, "foo", fmt.Errorf("The error")
				})
				vals, err = cli.Call(cmd)
				So(len(vals), ShouldEqual, 2)
				So(err, ShouldResemble, NewCallError(fmt.Errorf("The error")))
				So(strings.Join(val, " ** "), ShouldEqual, "called")
				So(vals[0].Interface(), ShouldResemble, 10)
				So(vals[1].Interface(), ShouldResemble, "foo")
			})
		})

		Convey("With pre call on command", func() {
			val = []string{}
			cmd.SetPreCall(func(c *Command) {
				val = append(val, fmt.Sprintf("pre (%s)", c.Name))
			})
			_, err = cli.Call(cmd)
			So(err, ShouldBeNil)
			So(strings.Join(val, " ** "), ShouldEqual, "pre (foo) ** called")

			Convey("With values, which are not delegated", func() {
				val = []string{}
				cmd.SetPreCall(func(c *Command) (int, string) {
					val = append(val, fmt.Sprintf("pre (%s)", c.Name))
					return 10, "foo"
				})
				vals, err := cli.Call(cmd)
				So(err, ShouldBeNil)
				So(len(vals), ShouldEqual, 0)
				So(strings.Join(val, " ** "), ShouldEqual, "pre (foo) ** called")
			})

			Convey("With error, which is delegated", func() {
				val = []string{}
				cmd.SetPreCall(func(c *Command) error {
					val = append(val, fmt.Sprintf("pre (%s)", c.Name))
					return fmt.Errorf("Pre Error")
				})
				_, err = cli.Call(cmd)
				So(err, ShouldResemble, NewCallError(fmt.Errorf("Pre Error")))
				So(strings.Join(val, " ** "), ShouldEqual, "pre (foo)")

				Convey("With error and values, which are not delegated", func() {
					val = []string{}
					cmd.SetPreCall(func(c *Command) (int, string, error) {
						val = append(val, fmt.Sprintf("pre (%s)", c.Name))
						return 10, "foo", fmt.Errorf("Pre Error")
					})
					vals, err := cli.Call(cmd)
					So(err, ShouldResemble, NewCallError(fmt.Errorf("Pre Error")))
					So(len(vals), ShouldEqual, 0)
					So(strings.Join(val, " ** "), ShouldEqual, "pre (foo)")
				})
			})

			Convey("With post call", func() {
				val = []string{}
				cmd.SetPostCall(func(c *Command) {
					val = append(val, fmt.Sprintf("post (%s)", c.Name))
				})
				_, err = cli.Call(cmd)
				So(err, ShouldBeNil)
				So(strings.Join(val, " ** "), ShouldEqual, "pre (foo) ** called ** post (foo)")
			})
		})

		Convey("With post call", func() {
			val = []string{}
			cmd.SetPostCall(func(c *Command) {
				val = append(val, fmt.Sprintf("post (%s)", c.Name))
			})
			_, err = cli.Call(cmd)
			So(err, ShouldBeNil)
			So(strings.Join(val, " ** "), ShouldEqual, "called ** post (foo)")

			Convey("With values, which are not delegated", func() {
				val = []string{}
				cmd.SetPostCall(func(c *Command) (int, string) {
					val = append(val, fmt.Sprintf("post (%s)", c.Name))
					return 10, "foo"
				})
				vals, err := cli.Call(cmd)
				So(err, ShouldBeNil)
				So(len(vals), ShouldEqual, 0)
				So(strings.Join(val, " ** "), ShouldEqual, "called ** post (foo)")
			})

			Convey("With error, which is delegated", func() {
				val = []string{}
				cmd.SetPostCall(func(c *Command) error {
					val = append(val, fmt.Sprintf("post (%s)", c.Name))
					return fmt.Errorf("Post Error")
				})
				_, err = cli.Call(cmd)
				So(err, ShouldResemble, NewCallError(fmt.Errorf("Post Error")))
				So(strings.Join(val, " ** "), ShouldEqual, "called ** post (foo)")

				Convey("With error and values, which are not delegated", func() {
					val = []string{}
					cmd.SetPostCall(func(c *Command) (int, string, error) {
						val = append(val, fmt.Sprintf("post (%s)", c.Name))
						return 10, "foo", fmt.Errorf("Post Error")
					})
					vals, err := cli.Call(cmd)
					So(err, ShouldResemble, NewCallError(fmt.Errorf("Post Error")))
					So(strings.Join(val, " ** "), ShouldEqual, "called ** post (foo)")
					So(len(vals), ShouldEqual, 0)
				})
			})
		})
	})
}

func TestCliRun(t *testing.T) {
	Convey("Run cli command", t, func() {
		called := 0
		var handledErr error
		Die = func(msg string, args ...interface{}) {
			panic(fmt.Sprintf(msg, args...))
		}
		Exit = func(s int) {
			panic(fmt.Sprintf("Exit %d", s))
		}
		namedActual := make(map[string]interface{})

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
			New("named", "", func(named NamedParameters) {
			namedActual = map[string]interface{}(named)
		}).
			New("named2", "", func(x testCliAlias, named NamedParameters, y *testCliInject) {
			namedActual = map[string]interface{}(named)
		}).
			Register(&testCliInject{
			Foo: 100,
		}).
			RegisterAs("clif.testCliAlias", &testCliInject{
			Foo: 200,
		})

		cmdInvalid := NewCommand("bla", "Dont use me", func() {})
		argInvalid := NewArgument("something", "..", "", false, false)
		argInvalid.SetParse(func(name, value string) (string, error) {
			return "", fmt.Errorf("Never works!")
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
		Convey("Run existing method with named parameters", func() {
			c.RegisterNamed("foo", "bar")
			c.RegisterNamed("baz", 213)
			c.RunWith([]string{"named"})
			So(namedActual, ShouldResemble, map[string]interface{}{"foo": "bar", "baz": 213})
		})
		Convey("Run existing method with named parameters on arbitrary position", func() {
			c.RegisterNamed("foo", "bar")
			c.RegisterNamed("baz", 213)
			c.RunWith([]string{"named2"})
			So(namedActual, ShouldResemble, map[string]interface{}{"foo": "bar", "baz": 213})
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
			}, ShouldPanicWith, `Callback parameter of type io.Writer for command "oops" was not found in registry`)
		})
		Convey("Run method with invalid arg fails", func() {
			So(func() {
				buf := bytes.NewBuffer(nil)
				out := NewOutput(buf, NewDefaultFormatter(map[string]string{}))
				c.SetOutput(out)
				c.RunWith([]string{"bla", "bla"})
			}, ShouldPanicWith, "Parse error: Parameter \"something\" invalid: Never works!")
		})
		Convey("Run method with resulting error returns it", func() {
			So(func() {
				c.RunWith([]string{"errme"})
			}, ShouldPanicWith, "Failure in execution: I error!")
		})
		Convey("Run with cli-wide pre call", func() {
			name := "NOT"
			c.SetPreCall(func(c *Command) error {
				name = c.Name
				return nil
			})
			c.RunWith([]string{"bar"})
			So(name, ShouldEqual, "bar")
			Convey("Run with error in cli-wide pre call", func() {
				c.SetPreCall(func(c *Command) error {
					return fmt.Errorf("Abort "+ c.Name)
				})
				So(func() {
					c.RunWith([]string{"bar"})
				}, ShouldPanicWith, `Abort bar`)
			})
		})
	})
}

func TestCliConstruction(t *testing.T) {
	Convey("Create new Cli with commands", t, func() {
		app := New("My App", "1.0.0", "Testing app")
		cb := func() {}

		Convey("Two default commands exist", func() {
			So(len(app.Commands), ShouldEqual, 2)
			Convey("One is \"help\"", func() {
				_, ok := app.Commands["help"]
				So(ok, ShouldBeTrue)
				Convey("Other is \"list\"", func() {
					_, ok := app.Commands["list"]
					So(ok, ShouldBeTrue)
				})
			})
		})

		Convey("Command constructur adds new command", func() {
			app.New("foo", "For fooing", cb)
			So(len(app.Commands), ShouldEqual, 3)
			So(app.Commands["foo"], ShouldNotBeNil)
		})

		Convey("Adding can be used variadic", func() {
			app.New("foo", "For fooing", cb)
			cmds := []*Command{
				NewCommand("foo", "For fooing", cb),
				NewCommand("bar", "For baring", cb),
			}
			app.Add(cmds...)
			So(len(app.Commands), ShouldEqual, 4)
			So(app.Commands["foo"], ShouldNotBeNil)
			So(app.Commands["bar"], ShouldNotBeNil)
		})
	})
}

func TestCliDefaultCommand(t *testing.T) {
	Convey("Change default command of cli", t, func() {
		x := 0
		app := New("My App", "1.0.0", "Testing app").
			SetDefaultCommand("other").
			New("other", "Something else", func() { x += 1 })
		So(app.DefaultCommand, ShouldEqual, "other")
		Convey("Calling default command", func() {
			app.RunWith(nil)
			So(x, ShouldEqual, 1)
		})
	})
}

func TestCliDefaultOptions(t *testing.T) {
	Convey("Adding default options to cli", t, func() {
		app := New("My App", "1.0.0", "Testing app")
		So(len(app.DefaultOptions), ShouldEqual, 0)

		Convey("Using default option creator adds option", func() {
			app.NewDefaultOption("foo", "f", "fooing", "", false, false)
			So(len(app.DefaultOptions), ShouldEqual, 1)
		})

		Convey("Adding default option .. adds them", func() {
			app.AddDefaultOptions(
				NewOption("foo", "f", "fooing", "", false, false),
				NewOption("bar", "b", "baring", "", false, false),
			)
			So(len(app.DefaultOptions), ShouldEqual, 2)
		})

		Convey("Cli default options are not added to command on command create", func() {
			app.NewDefaultOption("foo", "f", "fooing", "", false, false)
			cmd := NewCommand("bla", "bla", func() {})
			app.Add(cmd)
			So(len(cmd.Options), ShouldEqual, len(DefaultOptions))

			Convey("Default options are added in run", func() {
				app.RunWith([]string{"bla"})
				So(len(cmd.Options), ShouldEqual, len(DefaultOptions)+1)
			})
		})
	})
}

func TestCliHeralds(t *testing.T) {
	Convey("Command heralds are add late, in run", t, func() {
		app := New("My App", "1.0.0", "Testing app")
		So(len(app.Commands), ShouldEqual, 2)
		So(len(app.Heralds), ShouldEqual, 0)

		Convey("Heralding command does not add it to list", func() {
			x := 0
			app.Herald(func(c *Cli) *Command {
				return NewCommand("foo", "fooing", func() { x = 2 })
			})
			So(len(app.Commands), ShouldEqual, 2)
			So(len(app.Heralds), ShouldEqual, 1)

			Convey("Running adds heralded commands", func() {
				app.RunWith([]string{"foo"})
				So(x, ShouldEqual, 2)
				So(len(app.Commands), ShouldEqual, 3)
				So(len(app.Heralds), ShouldEqual, 0)
			})
		})

	})
}

func TestCliNamedRegistryParameter(t *testing.T) {
	Convey("Setting and accessing named parameter in registry", t, func() {
		app := New("My App", "1.0.0", "Testing app")
		app.RegisterNamed("foo", "foo")
		app.RegisterNamed("bar", 123)
		obj := &testCliInject{}
		app.RegisterNamed("baz", obj)

		Convey("Accessing named parameters", func() {
			So(app.Named("foo"), ShouldEqual, "foo")
			So(app.Named("bar"), ShouldEqual, 123)
			So(app.Named("baz"), ShouldEqual, obj)

			Convey("Accessing not existing named paraemter", func() {
				So(app.Named("zoing"), ShouldBeNil)
			})
		})
	})
}

var testCliSeparateArgs = []struct {
	args       []string
	expectName string
	expectArgs []string
}{
	{
		args:       []string{"foo", "--bar", "baz"},
		expectName: "foo",
		expectArgs: []string{"--bar", "baz"},
	},
	{
		args:       []string{"foo", "-bar", "baz"},
		expectName: "foo",
		expectArgs: []string{"-bar", "baz"},
	},
	{
		args:       []string{"foo", "baz", "--bar"},
		expectName: "foo",
		expectArgs: []string{"baz", "--bar"},
	},
	{
		args:       []string{"--bar", "foo", "baz"},
		expectName: "foo",
		expectArgs: []string{"--bar", "baz"},
	},
	{
		args:       []string{"--bar=boing", "foo", "baz"},
		expectName: "foo",
		expectArgs: []string{"--bar=boing", "baz"},
	},
	{
		args:       []string{"-bar=boing", "foo", "baz"},
		expectName: "foo",
		expectArgs: []string{"-bar=boing", "baz"},
	},
	{
		args:       []string{"--bar", "boing", "foo", "baz"},
		expectName: "boing",
		expectArgs: []string{"--bar", "foo", "baz"},
	},
}

func TestCliSeparateArgs(t *testing.T) {
	Convey("Separate command line args", t, func() {
		app := New("My App", "1.0.0", "Testing app")

		for i, test := range testCliSeparateArgs {
			Convey(fmt.Sprintf("%d) From %v", i, test.args), func() {
				cname, cargs := app.SeparateArgs(test.args)
				So(cname, ShouldEqual, test.expectName)
				So(cargs, ShouldResemble, test.expectArgs)
			})
		}
	})
}
