package clif
import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"fmt"
	"strings"
	"os"
)

var testsWrapText = []struct{
	from string
	expect string
}{
	/**/
	{
		from: "",
		expect: "",
	},
	{
		from: "foo",
		expect: "foo",
	},
	{
		from: "foo bar baz",
		expect: "foo bar baz",
	},

	{
		from: "foo   bar   baz",
		expect: "foo bar baz",
	},
	{
		from: "foobarbaz",
		expect: "foobarbaz",
	},
	{
		from: "foobarbaz foobarbaz",
		expect: "foobarbaz\nfoobarbaz",
	},
	{
		from: "foo bar baz foo bar baz",
		expect: "foo bar baz\nfoo bar baz",
	},
	{
		from: "foo bar baz\nfoo bar baz",
		expect: "foo bar baz\nfoo bar baz",
	},
	{
		from: "foo bar baz foo bar baz foo bar baz",
		expect: "foo bar baz\nfoo bar baz\nfoo bar baz",
	},
	{
		from: "foo bar baz\nfoo bar baz\nfoo bar baz",
		expect: "foo bar baz\nfoo bar baz\nfoo bar baz",
	},
	{
		from: "foobarbazfoobarbazfoobarbaz",
		expect: "foobarbazfoo\nbarbazfoobar\nbaz",
	},
	{
		from: "foobarbazfoobarbazfoobarbazfoobarbaz",
		expect: "foobarbazfoo\nbarbazfoobar\nbazfoobarbaz",
	},
	{
		from: "xx foobarbazfoobarbazfoobarbaz yy",
		expect: "xx\nfoobarbazfoo\nbarbazfoobar\nbaz yy",
	},
	{
		from: "foo bar baz foo bar baz",
		expect: "foo bar baz\nfoo bar baz",
	},
	{
		from: "\033[34mfoo\033[0m",
		expect: "\033[34mfoo\033[0m",
	},
	{
		from: "\033[34mfoo\033[0m bar baz foo bar baz",
		expect: "\033[34mfoo\033[0m bar baz\nfoo bar baz",
	},
	{
		from: "\033[34mfoo bar baz foo bar\033[0m baz",
		expect: strings.Join([]string{
			"\033[34mfoo bar baz\033[0m",
			"\033[34mfoo bar\033[0m baz",
		}, "\n"),
	},
	{
		from: "\033[34mfoo \033[1mbar",
		expect: "\033[34mfoo \033[1mbar\033[0m",
	},
	{
		from: "\033[34mfoo \033[1mbar baz\nfoo bar \033[0mbaz",
		expect: "\033[34mfoo \033[1mbar baz\033[0m\n\033[34m\033[1mfoo bar \033[0mbaz",
	},
	{
		from: "\033[34mfoobarbazfoobarbazfoobarbazfoobarbaz",
		expect: "\033[34mfoobarbazfoo\033[0m\n\033[34mbarbazfoobar\033[0m\n\033[34mbazfoobarbaz\033[0m",
	},
}

func TestWrapText(t *testing.T) {
	Convey("Wrapping text", t, func() {
		wrapper := NewWrapper(12)
		wrapper.WhitespaceMode = WRAP_WHITESPACE_CONTRACT
		wrapper.TrimMode = WRAP_TRIM_RIGHT
		wrapper.BreakWords = true
		for idx, test := range testsWrapText {
			out := _stringRenderDump(test.from)
			Convey(fmt.Sprintf("%d) %s\033[0m", idx, out), func() {
				to := wrapper.Wrap(test.from)
				if dbg := os.Getenv("WRAP_DEBUG"); dbg == "yes" || dbg == "1" {
					_stringCompareDump(to, test.expect)
				}
				So(to, ShouldEqual, test.expect)
			})
		}
	})
}

func TestWrapTextTrim(t *testing.T) {
	Convey("Wrap vs Trimming", t, func() {
		wrapper := NewWrapper(8)
		wrapper.WhitespaceMode = WRAP_WHITESPACE_KEEP
		wrapper.BreakWords = true
		Convey("Right trim none", func() {
			wrapper.TrimMode = WRAP_TRIM_NONE
			So(wrapper.Wrap("  Line 1  \n  Line 2  "), ShouldEqual, "  Line 1  \n  Line 2  ")
		})
		Convey("Right trim", func() {
			wrapper.TrimMode = WRAP_TRIM_RIGHT
			So(wrapper.Wrap("  Line 1  \n  Line 2  "), ShouldEqual, "  Line 1\n  Line 2")
		})
		Convey("Left trim", func() {
			wrapper.TrimMode = WRAP_TRIM_LEFT
			So(wrapper.Wrap("  Line 1  \n  Line 2  "), ShouldEqual, "Line 1  \nLine 2  ")
		})
		Convey("Both trim", func() {
			wrapper.TrimMode = WRAP_TRIM_BOTH
			So(wrapper.Wrap("  Line 1  \n  Line 2  "), ShouldEqual, "Line 1\nLine 2")
		})
		Convey("Trim without whitespace mode", func() {
			wrapper.WhitespaceMode = WRAP_WHITESPACE_CONTRACT
			wrapper.TrimMode = WRAP_TRIM_BOTH
			wrapped := wrapper.Wrap("  Line 1  \n  Line 2  ")
			expect := "Line 1\nLine 2"
			if dbg := os.Getenv("WRAP_DEBUG"); dbg == "yes" || dbg == "1" {
				_stringCompareDump(wrapped, expect)
			}
			So(wrapped, ShouldEqual, expect)
		})
	})
}