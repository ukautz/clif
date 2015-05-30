package cli

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestDescriberCommand(t *testing.T) {
	Convey("Description of Command", t, func() {
		d := &DefaultDescriber{}
		c := NewCommand("foo", "It does foo", func(){}).
			SetDescription("It does really, really foo")

		c.NewArgument("bar", "The bar", "", true, false).
			NewArgument("baz", "The baz", "", false, true)

		c.NewOption("boing", "b", "The boing!", "", true, false).
			NewOption("zoing", "z", "The ZOING!", "", false, true)
		c.Option("zoing").IsFlag()

		s := d.Command(c)
		expect := `It does really, really foo

Usage:
	foo <bar> ([baz] ...) --boing|-b <val> ((--zoing|-z) ...)

Arguments:
	bar  The bar (req)
	baz  The baz (mult)

Options:
	--boing|-b <val>  The boing! (req)
	--zoing|-z        The ZOING! (mult)
`
		So(s, ShouldEqual, expect)
	})
}

func TestDescriberCli(t *testing.T) {
	Convey("Description of Cli", t, func() {
		d := &DefaultDescriber{}
		c := New("cli", "1.0.1", "My CLI").
			New("foo", "It does foo", func(){}).
			New("bar", "It does bar", func(){}).
			New("bazoing", "It does bazoing", func(){});

		s := d.Cli(c)
		expect := `cli (1.0.1)
My CLI

Usage:
	cli.test <command> [<arg> ..] [--opt <val> ..]

Available commands:
	bar      It does bar
	bazoing  It does bazoing
	foo      It does foo
	help     Show this help`
		So(s, ShouldEqual, expect)
	})
}
