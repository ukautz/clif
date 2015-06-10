// +build never

package main

import (
	"github.com/ukautz/go-cli"
	"reflect"
	"os"
	"fmt"
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

func callHello(out cli.Output) {
	out.Printf("Hello World\n")
}

func callFoo(c *cli.Command, out cli.Output, custom1 exampleInterface, custom2 *exampleStruct) {
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

func main() {
	// initialize the app with custom registered objects in the injection container
	c := cli.New("My App", "1.0.0", "An example application").
		Register(&exampleStruct{"bar1"}).
		RegisterAs(reflect.TypeOf((*exampleInterface)(nil)).Elem().String(), &exampleStruct{"bar2"}).
		New("hello", "The obligatory hello world", callHello)

	// extend output styles
	cli.DefaultStyles["mine"] = "\033[32;1m"

	// customize error handler
	cli.Die = func(msg string, args ...interface{}) {
		c.Output().Printf("<error>Everyting went wrong: %s<reset>\n\n", fmt.Sprintf(msg, args...))
		os.Exit(1)
	}

	// build & add a complex command
	cmd := cli.NewCommand("foo", "It does foo", callFoo).
		NewArgument("name", "Name for greeting", "", true, false).
		NewArgument("more-names", "And more names for greeting", "", false, true).
		NewOption("whatever", "w", "Some required option", "", true, false)
	cnt := cli.NewOption("counter", "c", "Show how high you can count", "", false, false)
	cnt.SetValidator(cli.IsInt)
	cmd.AddOption(cnt)
	c.Add(cmd)

	cb := func(c *cli.Command, out cli.Output) {
		out.Printf("Called %s\n", c.Name)
	}
	c.New("bar:baz", "A grouped command", cb).
		New("bar:zoing", "Another grouped command", cb).
		New("hmm:huh", "Yet another grouped command", cb).
		New("hmm:uhm", "And yet another grouped command", cb)

	// execute the main loop
	c.Run()
}
