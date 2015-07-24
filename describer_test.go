package clif

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestDescribeCommand(t *testing.T) {
	Convey("Description of Command", t, func() {
		c := NewCommand("foo", "It does foo", func() {}).
			SetDescription("It does really, really foo")

		c.NewArgument("bar", "The bar", "", true, false).
			NewArgument("baz", "The baz", "", false, true)

		c.NewOption("boing", "b", "The boing!", "", true, false).
			NewOption("zoing", "z", "The ZOING!", "", false, true)
		c.Option("zoing").IsFlag().SetEnv("THE_ZOING")

		s := DescribeCommand(c)
		expect := `Command: <headline>foo<reset>
<info>It does really, really foo<reset>

<subline>Usage:<reset>
  foo bar [baz ...] [--help|-h] --boing|-b val [--zoing|-z ...]

<subline>Arguments:<reset>
  <info>bar<reset>  The bar (<important>req<reset>)
  <info>baz<reset>  The baz (<debug>mult<reset>)

<subline>Options:<reset>
  <info>--help|-h     <reset>  Display this help message
  <info>--boing|-b val<reset>  The boing! (<important>req<reset>)
  <info>--zoing|-z    <reset>  The ZOING! (<debug>mult<reset>, env: <debug>THE_ZOING<reset>)

`
		So(s, ShouldEqual, expect)
	})
}

func TestDescriberCli(t *testing.T) {
	Convey("Description of Cli", t, func() {
		c := New("cli", "1.0.1", "My CLI").
			New("foo", "It does foo", func() {}).
			New("bar", "It does bar", func() {}).
			New("bazoing", "It does bazoing", func() {}).
			New("zzz:uno", "A sub of zzz", func() {}).
			New("zzz:due", "A sub of zzz", func() {}).
			New("bla:due", "A sub of bla", func() {}).
			New("bla:uno", "A sub of bla", func() {})

		s := DescribeCli(c)
		expect := `<headline>cli<reset> <debug>(1.0.1)<reset>
<info>My CLI<reset>

<subline>Usage:<reset>
  clif.test command [arg ..] [--opt val ..]

<subline>Available commands:<reset>
  <info>bar    <reset>  It does bar
  <info>bazoing<reset>  It does bazoing
  <info>foo    <reset>  It does foo
  <info>help   <reset>  Show this help
  <info>list   <reset>  List all available commands
 <subline>bla<reset>
  <info>bla:due<reset>  A sub of bla
  <info>bla:uno<reset>  A sub of bla
 <subline>zzz<reset>
  <info>zzz:due<reset>  A sub of zzz
  <info>zzz:uno<reset>  A sub of zzz
`
		So(s, ShouldEqual, expect)
	})
}
