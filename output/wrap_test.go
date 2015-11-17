package output
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