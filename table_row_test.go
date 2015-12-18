package clif
import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"strings"
	"fmt"
)

var testsTableRow = []struct{
	content []string
	expectCalcContent []string
	expectLineCount int
	expectCalcLineCount int
}{
	{
		content: []string{"foo", "bar", "baz"},
		expectCalcContent: []string{
			"foo" + strings.Repeat(" ", 7),
			"bar" + strings.Repeat(" ", 7),
			"baz" + strings.Repeat(" ", 7),
		},
		expectLineCount: 1,
		expectCalcLineCount: 1,
	},
	{
		content: []string{"foo\nfoo", "bar", "baz"},
		expectCalcContent: []string{
			"foo" + strings.Repeat(" ", 7) + "\nfoo" + strings.Repeat(" ", 7),
			"bar" + strings.Repeat(" ", 7),
			"baz" + strings.Repeat(" ", 7),
		},
		expectLineCount: 2,
		expectCalcLineCount: 2,
	},
	{
		content: []string{"foofoo", "bar", "baz"},
		expectCalcContent: []string{
			"foofoo" + strings.Repeat(" ", 9),
			"bar" + strings.Repeat(" ", 4),
			"baz" + strings.Repeat(" ", 5),
		},
		expectLineCount: 1,
		expectCalcLineCount: 1,
	},
	{
		content: []string{"foofoofoofoo", "bar", "baz"},
		expectCalcContent: []string{
			"foofoofoofoo" + strings.Repeat(" ", 8),
			"bar" + strings.Repeat(" ", 2),
			"baz" + strings.Repeat(" ", 2),
		},
		expectLineCount: 1,
		expectCalcLineCount: 1,
	},
	{
		content: []string{"foo foo foo foo", "bar bar bar bar", "baz baz baz baz"},
		expectCalcContent: []string{
			"foo foo   \nfoo foo   ",
			"bar bar   \nbar bar   ",
			"baz baz   \nbaz baz   ",
		},
		expectLineCount: 1,
		expectCalcLineCount: 2,
	},
	{
		content: []string{"foofoofoofoo", "barbarbarbar", "bazbazbazbaz"},
		expectCalcContent: []string{
			"foofoofoof\noo        ",
			"barbarbarb\nar        ",
			"bazbazbazb\naz        ",
		},
		expectLineCount: 1,
		expectCalcLineCount: 2,
	},
	{
		content: []string{"foofoofoofoo foofoofoofoo", "barbarbarbar", "bazbazbazbaz"},
		expectCalcContent: []string{
			"foofoofoofoo   \nfoofoofoofoo   ",
			"barbarb\narbar  ",
			"bazbazba\nzbaz    ",
		},
		expectLineCount: 1,
		expectCalcLineCount: 2,
	},
}

func TestTableRow(t *testing.T) {
	Convey("Create new row from data", t, func() {
		for idx, test := range testsTableRow {
			out := strings.Replace(strings.Join(test.content, " | "), "\n", "\\n", -1)
			if l := len(out); l > 30 {
				out = fmt.Sprintf("%s ... (%d)", out[0:30], l)
			}
			Convey(fmt.Sprintf("%d) %s", idx, out), func() {
				row := NewTableRow(test.content)

				Convey("Rendered without width restrictions", func() {
					rendered, lineCount := row.Render(0)
					So(rendered, ShouldResemble, test.content)
					So(lineCount, ShouldResemble, test.expectLineCount)
				})

				Convey("Rendered with width restrictions", func() {
					rendered, lineCount := row.Render(30)
					So(rendered, ShouldResemble, test.expectCalcContent)
					So(lineCount, ShouldResemble, test.expectCalcLineCount)
				})
			})
		}
	})
}