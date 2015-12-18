package clif

import (
	"bytes"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"strings"
	"testing"
)

var testsTableCol = []struct {
	content               string
	expectContent         string
	expectCalcContent     string
	expectRenderedContent string
	expectWidth           int
	expectCalcWidth       int
	expectLineCount       int
	expectCalcLineCount   int
}{
	{
		content:               "",
		expectContent:         "",
		expectCalcContent:     "",
		expectRenderedContent: strings.Repeat(" ", 10),
		expectWidth:           0,
		expectCalcWidth:       0,
		expectLineCount:       1,
		expectCalcLineCount:   1,
	},
	{
		content:               "foo",
		expectContent:         "foo",
		expectCalcContent:     "foo",
		expectRenderedContent: "foo" + strings.Repeat(" ", 7),
		expectWidth:           3,
		expectCalcWidth:       3,
		expectLineCount:       1,
		expectCalcLineCount:   1,
	},
	{
		content:               "foo   ",
		expectContent:         "foo",
		expectCalcContent:     "foo",
		expectRenderedContent: "foo" + strings.Repeat(" ", 7),
		expectWidth:           3,
		expectCalcWidth:       3,
		expectLineCount:       1,
		expectCalcLineCount:   1,
	},
	{
		content:           "foo bar baz",
		expectContent:     "foo bar baz",
		expectCalcContent: "foo bar\nbaz",
		expectRenderedContent: strings.Join([]string{
			"foo bar" + strings.Repeat(" ", 3),
			"baz" + strings.Repeat(" ", 7),
		}, "\n"),
		expectWidth:         11,
		expectCalcWidth:     7,
		expectLineCount:     1,
		expectCalcLineCount: 2,
	},
	{
		content:           "foo\nbarrr\nbaz",
		expectContent:     "foo\nbarrr\nbaz",
		expectCalcContent: "foo\nbarrr\nbaz",
		expectRenderedContent: strings.Join([]string{
			"foo" + strings.Repeat(" ", 7),
			"barrr" + strings.Repeat(" ", 5),
			"baz" + strings.Repeat(" ", 7),
		}, "\n"),
		expectWidth:         5,
		expectCalcWidth:     5,
		expectLineCount:     3,
		expectCalcLineCount: 3,
	},
	{
		content:               "0123456789 0123456789",
		expectContent:         "0123456789 0123456789",
		expectCalcContent:     "0123456789\n0123456789",
		expectRenderedContent: "0123456789\n0123456789",
		expectWidth:           21,
		expectCalcWidth:       10,
		expectLineCount:       1,
		expectCalcLineCount:   2,
	},
	{
		content:               "01234567890123456789",
		expectContent:         "01234567890123456789",
		expectCalcContent:     "0123456789\n0123456789",
		expectRenderedContent: "0123456789\n0123456789",
		expectWidth:           20,
		expectCalcWidth:       10,
		expectLineCount:       1,
		expectCalcLineCount:   2,
	},
	{
		content:               strings.Repeat("0123456789", 100),
		expectContent:         strings.Repeat("0123456789", 100),
		expectCalcContent:     "0123456789" + strings.Repeat("\n0123456789", 99),
		expectRenderedContent: "0123456789" + strings.Repeat("\n0123456789", 99),
		expectWidth:           1000,
		expectCalcWidth:       10,
		expectLineCount:       1,
		expectCalcLineCount:   100,
	},
	{
		content:               strings.Repeat("012345678 ", 10),
		expectContent:         "012345678" + strings.Repeat(" 012345678", 9),
		expectCalcContent:     "012345678" + strings.Repeat("\n012345678", 9),
		expectRenderedContent: "012345678 " + strings.Repeat("\n012345678 ", 9),
		expectWidth:           99,
		expectCalcWidth:       9,
		expectLineCount:       1,
		expectCalcLineCount:   10,
	},
}

func TestTableCol(t *testing.T) {
	Convey("Create new col from data", t, func() {
		for idx, test := range testsTableCol {
			out := strings.Replace(test.content, "\n", "\\n", -1)
			if l := len(out); l > 30 {
				out = fmt.Sprintf("%s ... (%d)", out[0:30], l)
			}
			Convey(fmt.Sprintf("%d) \"%s\"", idx, out), func() {
				col := NewTableCol(nil, test.content)

				Convey("Rendered without width restrictions", func() {
					So(col.Content(), ShouldEqual, test.expectContent)
					So(col.Width(), ShouldEqual, test.expectWidth)
					So(col.LineCount(), ShouldEqual, test.expectLineCount)
				})

				Convey("Rendered with width restrictions", func() {
					//fmt.Printf("\n  IS CNT=%s, OUGHT CNT=%s\n", col.Content(10), test.expectCalcContent)
					So(col.Content(10), ShouldEqual, test.expectCalcContent)
					//fmt.Printf("\n  IS WIDTH=%d, OUGHT WIDTH=%d\n", col.Width(10), test.expectCalcWidth)
					So(col.Width(10), ShouldEqual, test.expectCalcWidth)
					So(col.LineCount(10), ShouldEqual, test.expectCalcLineCount)

					rendered, _, _ := col.Render(10)
					//fmt.Printf("\n  IS=\"%s\", OUGHT=\"%s\"\n", strings.Replace(rendered, "\n", " ** ", -1), strings.Replace(test.expectRenderedContent, "\n", " ** ", -1))
					So(rendered, ShouldEqual, test.expectRenderedContent)
				})
			})
		}
	})
}

func TestTableColWithFormatting(t *testing.T) {
	Convey("Create new col with formatting", t, func() {
		tests := map[string][]string{
			/*"Unformatted": []string{
				"foo bar baz foo bar baz!",
				strings.Join([]string{
					"foo bar   ",
					"baz foo   ",
					"bar baz!  ",
				}, "\n"),
			},*/
			"Formatted": []string{
				"foo <info>bar baz foo bar<reset> baz!",
				strings.Join([]string{
					"foo \033[34mbar\033[0m   ",
					"\033[34mbaz foo\033[0m   ",
					"\033[34mbar\033[0m baz!  ",
				}, "\n"),
			},
		}
		for name, ref := range tests {
			Convey(fmt.Sprintf("\n>>>>>>>> %s\n", name), func() {

				buf := bytes.NewBuffer(nil)
				out := NewColorOutput(buf)
				col := NewTableCol(nil, ref[0])
				col.SetRenderer(DefaultOutputTableContentRenderer(out))

				res, _, _ := col.Render(10)
				//_stringCompareDump(res, ref[1])
				//_testDumpStrings(res, ref[1])
				//fmt.Printf("\033[0m-- IS:\n%s\n\033[0m-- SHOULD:\n%s\n\033[0m--", res, ref[1])
				So(res, ShouldEqual, ref[1])
			})
		}
	})
}
