package clif

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestDefaultFormatterFormat(t *testing.T) {
	Convey("Formatting output", t, func() {
		Convey("Formatting with stripped directives", func() {
			f := NewDefaultFormatter(nil)
			s := f.Format("Foo <headline>bar<reset> baz")
			So(s, ShouldEqual, "Foo bar baz")
		})
		Convey("Formatting with formatted directives", func() {
			f := NewDefaultFormatter(map[string]string{
				"headline": "H!",
				"reset":    "R!",
			})
			s := f.Format("Foo <headline>bar<reset> baz")
			So(s, ShouldEqual, "Foo H!barR! baz")
		})
		Convey("Formatting does not replace not registered default tokens", func() {
			f := NewDefaultFormatter(nil)
			s := f.Format("Foo <headline>bar<reset> <baz> boing")
			So(s, ShouldEqual, "Foo bar <baz> boing")
		})
		Convey("Formatting does not replace not registered custom tokens", func() {
			f := NewDefaultFormatter(map[string]string{
				"baz": "BAZ",
			})
			s := f.Format("Foo <headline>bar<reset> <baz> boing")
			So(s, ShouldEqual, "Foo <headline>bar<reset> BAZ boing")
		})
		Convey("Formatting works over multi-line string", func() {
			f := NewDefaultFormatter(map[string]string{
				"headline": "H!",
				"reset":    "R!",
			})
			s := f.Format("Foo <headline>bar baz\ndings<reset> baz")
			So(s, ShouldEqual, "Foo H!bar baz\ndingsR! baz")
		})
	})
}

func TestDefaultFormatterEscape(t *testing.T) {
	Convey("Escapinng output", t, func() {
		Convey("Escaping tokens", func() {
			f := NewDefaultFormatter(nil)
			s := f.Escape("Foo <headline>bar<reset> baz")
			So(s, ShouldEqual, "Foo \\<headline>bar\\<reset> baz")
		})
		Convey("Escaping already escaped tokens does nothing", func() {
			f := NewDefaultFormatter(nil)
			s := f.Escape("Foo \\<headline>bar\\<reset> baz")
			So(s, ShouldEqual, "Foo \\<headline>bar\\<reset> baz")
		})
	})
}