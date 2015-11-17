// +build ignore

package main

import (
	"fmt"
	"gopkg.in/ukautz/clif.v1"
	"os"
	"reflect"
)

type exampleInterface interface {
	Foo() string
}

type exampleStruct struct {
	foo string
}

func (this *exampleStruct) Foo() string {
	return this.foo
}

func callHello(out clif.Output) {
	out.Printf("Hello <mine>World<reset>\n")
}

func callStyles(out clif.Output) {
	for token, _ := range clif.DefaultStyles {
		if token == "mine" {
			continue
		}
		out.Printf("Token \\<%s>: <%s>%s<reset>\n", token, token, token)
	}
}

func callFoo(c *clif.Command, out clif.Output, custom1 exampleInterface, custom2 *exampleStruct) {
	out.Printf("Hello %s, how is the %s?\n", c.Argument("name").String(), c.Option("whatever").String())
	if m := c.Argument("more-names").Strings(); m != nil && len(m) > 0 {
		for _, n := range m {
			out.Printf("  Say hello to <info>%s<reset>\n", n)
		}
	}
	if c.Option("counter").Int() > 5 {
		out.Printf("  You can count real high!\n")
	}
	out.Printf("  <headline>Custom 1: %s<reset>\n", custom1.Foo())
	out.Printf("  <subline>Custom 2: %s<reset>\n", custom2.foo)
}

/*
go run extended.go foo peter -w bla everybody -c=12 else

	Hello peter, how is the bla?
	  Say hello to everybody
	  Say hello to else
	  You can count real high!
	  Custom 1: bar2
	  Custom 2: bar1
*/

func setStyle(style string, c *clif.Cli) {
	switch style {
	case "sunburn":
		if c != nil {
			c.Output().SetFormatter(clif.NewDefaultFormatter(clif.SunburnStyles))
		}
		clif.DefaultStyles = clif.SunburnStyles
	case "winter":
		if c != nil {
			c.Output().SetFormatter(clif.NewDefaultFormatter(clif.WinterStyles))
		}
		clif.DefaultStyles = clif.WinterStyles
	}
}

func main() {
	setStyle(os.Getenv("CLI_STYLE"), nil)

	// extend output styles
	clif.DefaultStyles["mine"] = "\033[32;1m"

	// initialize the app with custom registered objects in the injection container
	c := clif.New("My App", "1.0.0", "An example application").
		Register(&exampleStruct{"bar1"}).
		RegisterAs(reflect.TypeOf((*exampleInterface)(nil)).Elem().String(), &exampleStruct{"bar2"}).
		New("hello", "The obligatory hello world", callHello)

	styleArg := clif.NewArgument("style", "Name of a style. Available: default, sunburn, winter", "default", true, false).
		SetParse(func(name, value string) (string, error) { setStyle(value, c); return value, nil })
	c.Add(clif.NewCommand("styles", "Print all color style tokens", callStyles).AddArgument(styleArg))

	// customize error handler
	clif.Die = func(msg string, args ...interface{}) {
		c.Output().Printf("<error>Everyting went wrong: %s<reset>\n\n", fmt.Sprintf(msg, args...))
		clif.Exit(1)
	}

	// build & add a complex command
	cmd := clif.NewCommand("foo", "It does foo", callFoo).
		NewArgument("name", "Name for greeting", "", true, false).
		NewArgument("more-names", "And more names for greeting", "", false, true).
		NewOption("whatever", "w", "Some required option", "", true, false)
	cnt := clif.NewOption("counter", "c", "Show how high you can count", "", false, false)
	cnt.SetParse(clif.IsInt)
	cmd.AddOption(cnt)
	c.Add(cmd)

	cb := func(c *clif.Command, out clif.Output) {
		out.Printf("Called %s\n", c.Name)
	}
	c.New("bar:baz", "A grouped command", cb).
		New("bar:zoing", "Another grouped command", cb).
		New("hmm:huh", "Yet another grouped command", cb).
		New("hmm:uhm", "And yet another grouped command", cb)

	// execute the main loop
	c.Run()
}
