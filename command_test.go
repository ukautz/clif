package clif

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"regexp"
	"strings"
	"testing"
	"time"
	"os"
)

func _testInitCommand() *Command {
	return &Command{
		Arguments: []*Argument{
			{
				parameter: parameter{
					Name:     "foo",
					Required: true,
					Regex:    regexp.MustCompile(`^a`),
				},
			},
			{
				parameter: parameter{
					Name:     "bar",
					Multiple: true,
					Parse: func(name, value string) (string, error) {
						if strings.Index(value, "B") == -1 {
							return "", fmt.Errorf("Missing B")
						} else {
							return value, nil
						}
					},
				},
			},
		},
		Options: []*Option{
			{
				parameter: parameter{
					Name:     "baz",
					Required: true,
				},
			},
			{
				parameter: parameter{
					Name:     "bang",
					Default:  "the default",
					Multiple: true,
				},
			},
			{
				parameter: parameter{
					Name:     "zoing",
					Multiple: true,
				},
				Flag:  true,
				Alias: "z",
			},
		},
	}
}

func TestCommandCreation(t *testing.T) {
	Convey("Creating command with constructor", t, func() {
		Convey("Adding command with empty signature and empty return", func() {
			c := NewCommand("name", "Usage", func() {
				//
			})
			So(c, ShouldNotBeNil)
		})
		Convey("Adding command with non-empty signature and empty return", func() {
			c := NewCommand("name", "Usage", func(foo string, bar int, baz time.Time) {
				//
			})
			So(c, ShouldNotBeNil)
		})
		Convey("Adding command with empty signature and non-empty return", func() {
			c := NewCommand("name", "Usage", func() (string, int, time.Time, error) {
				return "", 0, time.Now(), nil
			})
			So(c, ShouldNotBeNil)
		})
		Convey("Adding command with non-empty signature and non-empty return", func() {
			c := NewCommand("name", "Usage", func(foo string, bar int, baz time.Time) (string, int, time.Time, error) {
				return "", 0, time.Now(), nil
			})
			So(c, ShouldNotBeNil)
		})
		Convey("Adding command which is not a func", func() {
			So(func() {
				NewCommand("name", "Usage", "foo")
			}, ShouldPanicWith, "Call must be method, but is string")
		})
		Convey("Adding command which is also not a func", func() {
			So(func() {
				NewCommand("name", "Usage", time.Now())
			}, ShouldPanicWith, "Call must be method, but is struct")
		})
	})
}

var testsCommandParse = []struct {
	in   []string
	vals map[string][]string
	err  error
}{
	{
		in:  []string{},
		err: fmt.Errorf("Argument \"foo\" is required but missing"),
	},
	{
		in:  []string{"first"},
		err: fmt.Errorf("Parameter \"foo\" invalid: Does not match criteria"),
	},
	{
		in:  []string{""},
		err: fmt.Errorf("Parameter \"foo\" invalid: Does not match criteria"),
	},
	{
		in:  []string{"afirst"},
		err: fmt.Errorf("Option \"baz\" is required but missing"),
	},
	{
		in:  []string{"afirst", "second"},
		err: fmt.Errorf("Parameter \"bar\" invalid: Missing B"),
	},
	{
		in:  []string{"afirst", "second with B"},
		err: fmt.Errorf("Option \"baz\" is required but missing"),
	},
	{
		in: []string{"afirst", "--baz=x"},
		vals: map[string][]string{
			"foo":  []string{"afirst"},
			"baz":  []string{"x"},
			"bang": []string{"the default"},
		},
	},
	{
		in: []string{"afirst", "--baz", "x"},
		vals: map[string][]string{
			"foo":  []string{"afirst"},
			"baz":  []string{"x"},
			"bang": []string{"the default"},
		},
	},
	{
		in:  []string{"afirst", "--baz=x", "--baz=y"},
		err: fmt.Errorf("Parameter \"baz\" does not support multiple values"),
	},
	{
		in:  []string{"afirst", "--baz", "--baz", "x"},
		err: fmt.Errorf("Missing value for option \"--baz\""),
	},
	{
		in: []string{"afirst", "second with B", "--baz=x"},
		vals: map[string][]string{
			"foo":  []string{"afirst"},
			"bar":  []string{"second with B"},
			"baz":  []string{"x"},
			"bang": []string{"the default"},
		},
	},
	{
		in:  []string{"afirst", "second with B", "another bar param", "--baz=x"},
		err: fmt.Errorf("Parameter \"bar\" invalid: Missing B"),
	},
	{
		in: []string{"afirst", "second with B", "another Bar param", "--baz=x"},
		vals: map[string][]string{
			"foo":  []string{"afirst"},
			"bar":  []string{"second with B", "another Bar param"},
			"baz":  []string{"x"},
			"bang": []string{"the default"},
		},
	},
	{
		in: []string{"afirst", "--baz=x", "second with B", "another Bar param"},
		vals: map[string][]string{
			"foo":  []string{"afirst"},
			"bar":  []string{"second with B", "another Bar param"},
			"baz":  []string{"x"},
			"bang": []string{"the default"},
		},
	},
	{
		in: []string{"afirst", "--baz=x", "second with B", "another Bar param", "--bang", "Bang!"},
		vals: map[string][]string{
			"foo":  []string{"afirst"},
			"bar":  []string{"second with B", "another Bar param"},
			"baz":  []string{"x"},
			"bang": []string{"Bang!"},
		},
	},
	{
		in: []string{"afirst", "--baz=x", "second with B", "another Bar param", "--bang", "Bang!", "--bang=Bang!!"},
		vals: map[string][]string{
			"foo":  []string{"afirst"},
			"bar":  []string{"second with B", "another Bar param"},
			"baz":  []string{"x"},
			"bang": []string{"Bang!", "Bang!!"},
		},
	},
	{
		in:  []string{"afirst", "--baz=x", "second with B", "another Bar param", "--bang", "Bang!", "--zoing=z"},
		err: fmt.Errorf("Flag \"--zoing=z\" cannot have value"),
	},
	{
		in: []string{"afirst", "--baz=x", "second with B", "another Bar param", "--bang", "Bang!", "--zoing"},
		vals: map[string][]string{
			"foo":   []string{"afirst"},
			"bar":   []string{"second with B", "another Bar param"},
			"baz":   []string{"x"},
			"bang":  []string{"Bang!"},
			"zoing": []string{"true"},
		},
	},
	{
		in: []string{"afirst", "--baz=x", "second with B", "another Bar param", "--bang", "Bang!", "--zoing", "-z", "-z"},
		vals: map[string][]string{
			"foo":   []string{"afirst"},
			"bar":   []string{"second with B", "another Bar param"},
			"baz":   []string{"x"},
			"bang":  []string{"Bang!"},
			"zoing": []string{"true", "true", "true"},
		},
	},
	{
		in:  []string{"-="},
		err: fmt.Errorf("Malformed option \"-=\""),
	},
	{
		in:  []string{"--not-there"},
		err: fmt.Errorf("Unrecognized option \"--not-there\""),
	},
}

func TestCommandParse(t *testing.T) {
	Convey("Parse commands", t, func() {
		for i, test := range testsCommandParse {
			Convey(fmt.Sprintf("%2d) \"%s\"", i, strings.Join(test.in, "\", \"")), func() {
				c := _testInitCommand()
				err := c.Parse(test.in)
				So(err, ShouldResemble, test.err)
				if err == nil {
					So(c.Input(), ShouldResemble, test.vals)
				}
			})
		}
	})
}

func TestCommandParseFallsBackToDefaults(t *testing.T) {
	Convey("Parse commands with default values for options and/or arguments", t, func() {
		c := NewCommand("command", "Usage", func() {}).
			NewArgument("foo", "For fooing", "FOO", true, false).
			NewOption("bar", "b", "For baring", "BAR", true, false)
		Convey("Missing required argument defaults", func() {
			err := c.Parse([]string{})
			So(err, ShouldBeNil)
			So(c.Argument("foo").String(), ShouldEqual, "FOO")
		})
		Convey("Missing required option defaults", func() {
			err := c.Parse([]string{})
			So(err, ShouldBeNil)
			So(c.Option("bar").String(), ShouldEqual, "BAR")
		})
	})
}

func TestCommandParseFallsBackToEnvBeforeDefault(t *testing.T) {
	Convey("Parse commands with default values for options and/or arguments", t, func() {
		os.Setenv("the_foo", "FOO_ENV")
		os.Setenv("the_bar", "BAR_ENV")
		c := NewCommand("command", "Usage", func() {}).
			AddArgument(NewArgument("foo", "For fooing", "FOO", true, false).SetEnv("the_foo")).
			AddOption(NewOption("bar", "b", "For baring", "BAR", true, false).SetEnv("the_bar"))
		Convey("Missing required argument defaults", func() {
			err := c.Parse([]string{})
			So(err, ShouldBeNil)
			So(c.Argument("foo").String(), ShouldEqual, "FOO_ENV")
		})
		Convey("Missing required option defaults", func() {
			err := c.Parse([]string{})
			So(err, ShouldBeNil)
			So(c.Option("bar").String(), ShouldEqual, "BAR_ENV")
		})
	})
}

func TestCommandAccess(t *testing.T) {
	Convey("Comand argument & option access", t, func() {
		c := _testInitCommand()
		Convey("Existing argument accessible", func() {
			So(c.Argument("foo"), ShouldNotBeNil)
		})
		Convey("Not existing argument not accessible", func() {
			So(c.Argument("fooz"), ShouldBeNil)
		})
		Convey("Existing option accessible", func() {
			So(c.Option("baz"), ShouldNotBeNil)
		})
		Convey("Not existing option not accessible", func() {
			So(c.Option("bazz"), ShouldBeNil)
		})
	})
}

func TestCommandAddingArgument(t *testing.T) {
	Convey("Add arguments to command", t, func() {
		c := NewCommand("foo", "Doing foo", func(c *Command) error {
			return nil
		})
		Convey("Adding single argument", func() {
			c.NewArgument("bar", "A bar", "123", true, false)
			So(len(c.Arguments), ShouldEqual, 1)
			So(c.Arguments[0], ShouldResemble, &Argument{
				parameter: parameter{
					Name:     "bar",
					Usage:    "A bar",
					Default:  "123",
					Required: true,
				},
			})

			Convey("Adding additional argument with different name", func() {
				c.NewArgument("baz", "A baz", "", false, false)
				So(len(c.Arguments), ShouldEqual, 2)

				Convey("Adding argument of same name does not fly", func() {
					So(func() {
						c.NewArgument("baz", "A baz", "", false, false)
					}, ShouldPanicWith, `Argument with name "baz" already existing`)
				})
			})

			Convey("Adding multiple style argument", func() {
				c.NewArgument("baz", "A baz", "", false, true)
				So(len(c.Arguments), ShouldEqual, 2)

				Convey("Cannot add another argument after that", func() {
					So(func() {
						c.NewArgument("other", "A baz", "", false, false)
					}, ShouldPanicWith, `Cannot add argument after multiple style argument`)
				})
			})

			Convey("Adding argument with existing option does not fly", func() {
				c.Options = []*Option{
					{
						parameter: parameter{
							Name: "aaa",
						},
					},
					{
						parameter: parameter{
							Name: "bbb",
						},
						Alias: "b",
					},
				}
				So(func() {
					c.NewArgument("aaa", "A baz", "", false, false)
				}, ShouldPanicWith, `Option with name or alias "aaa" already existing`)
				So(func() {
					c.NewArgument("b", "A baz", "", false, false)
				}, ShouldPanicWith, `Option with name or alias "b" already existing`)
			})
		})
	})
}

func TestCommandAddingOptions(t *testing.T) {
	Convey("Add options to command", t, func() {
		c := NewCommand("foo", "Doing foo", func(c *Command) error {
			return nil
		})
		Convey("Adding single option", func() {
			c.NewOption("bar", "b", "A bar", "123", true, false)
			So(len(c.Options), ShouldEqual, 2)
			So(c.Options, ShouldResemble, []*Option{DefaultHelpOption, {
				parameter: parameter{
					Name:     "bar",
					Usage:    "A bar",
					Default:  "123",
					Required: true,
				},
				Alias: "b",
			}})

			Convey("Adding additional option with different name", func() {
				c.NewOption("baz", "", "A baz", "", false, false)
				So(len(c.Options), ShouldEqual, 3)

				Convey("Adding option with existing name does not fly", func() {
					So(func() {
						c.NewOption("baz", "", "A baz", "", false, false)
					}, ShouldPanicWith, `Option with name or alias "baz" already existing`)
				})

				Convey("Adding option with existing alias does not fly", func() {
					So(func() {
						c.NewOption("bazz", "b", "A baz", "", false, false)
					}, ShouldPanicWith, `Option with name or alias "b" already existing`)
				})
			})

			Convey("Adding option with existing argument does not fly", func() {
				c.Arguments = []*Argument{
					{
						parameter: parameter{
							Name: "a",
						},
					},
				}
				So(func() {
					c.NewOption("a", "", "A baz", "", false, false)
				}, ShouldPanicWith, `Argument with name "a" already existing`)
				So(func() {
					c.NewOption("aaa", "a", "A baz", "", false, false)
				}, ShouldPanicWith, `Cannot use alias: Argument with name "a" already existing`)
			})
		})
	})
}
