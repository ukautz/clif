package clif
import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"strings"
	"fmt"
)

var testsWrapText = []struct{
	from string
	expect string
}{
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
		from: "xx foobarbazfoobarbazfoobarbaz yy",
		expect: "xx\nfoobarbazfoo\nbarbazfoobar\nbaz yy",
	},
	{
		from: "foo bar baz foo bar baz",
		expect: "foo bar baz\nfoo bar baz",
	},
	{
		from: "\033[34mfoo\033[0m bar baz foo bar baz",
		expect: "\033[34mfoo\033[0m bar baz\nfoo bar baz",
	},
}

func TestWrapText(t *testing.T) {
	Convey("Wrapping text", t, func() {
		for idx, test := range testsWrapText {
			out := strings.Replace(test.from, "\n", "\\n", -1)
			Convey(fmt.Sprintf("%d) %s", idx, out), func() {
				to := WrapStringExtreme(test.from, 12)
				So(to, ShouldEqual, test.expect)
			})
		}
	})
}