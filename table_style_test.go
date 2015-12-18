package clif

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"regexp"
	"testing"
)

var _testTableRenderHeader = func(str string) string {
	//fmt.Printf("\n>>> USING HEADER RENDERER ON \"%s\"\n", strings.Replace(str, "\n", " ** ", -1))
	rxRight := regexp.MustCompile(`^(\s*)\**(\S.+?)\**(\s*)$`)
	if rxRight.MatchString(str) {
		match := rxRight.FindStringSubmatch(str)
		prefix := match[1]
		word := match[2]
		suffix := match[3]
		return fmt.Sprintf("%s*%s*%s", prefix, word, suffix)
	}
	return str
}

var testsTableStyle = []struct {
	data           [][]string
	expectedWidths []int
	renderedTable  string
}{
	{
		data: [][]string{
			{"foo", "bar", "baz"},
		},
		expectedWidths: []int{17, 17, 19},
		renderedTable:  `fixtures/table_0`,
	},
	{
		data: [][]string{
			{"foofoofoofoofoo", "bar", "baz"},
		},
		expectedWidths: []int{26, 12, 15},
		renderedTable:  `fixtures/table_1`,
	},
	{
		data: [][]string{
			{"foo\nfoo\nfoo\nfoo\nfoo", "bar", "baz"},
		},
		expectedWidths: []int{17, 17, 19},
		renderedTable:  `fixtures/table_2`,
	},
	{
		data: [][]string{
			{"foofoofoofoofoo", "bar", "baz"},
			{"foo", "barbarbarbarbar", "baz"},
		},
		expectedWidths: []int{21, 21, 11},
		renderedTable:  `fixtures/table_3`,
	},
	{
		data: [][]string{
			{"foofoofoofoofoo", "bar", "baz"},
			{"foo", "barbarbarbarbar", "baz"},
			{"foo", "bar", "bazbazbazbazbaz"},
		},
		expectedWidths: []int{17, 17, 19},
		renderedTable:  `fixtures/table_4`,
	},
	{
		data: [][]string{
			{"foofoofoofoofoofoofoofoofoofoofoofoofoofoofoo", "bar", "baz"},
			{"foo", "barbarbarbarbarbarbarbarbarbarbarbarbarbarbar", "baz"},
			{"foo", "bar", "bazbazbazbazbazbazbazbazbazbazbazbazbazbazbaz"},
		},
		expectedWidths: []int{17, 17, 19},
		renderedTable:  `fixtures/table_5`,
	},
	{
		data: [][]string{
			{
				"foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo",
				"bar bar bar bar bar bar bar bar bar bar bar bar bar bar bar bar bar bar bar",
				"baz baz baz baz baz baz baz baz baz baz baz baz baz baz baz baz baz baz baz",
			},
		},
		expectedWidths: []int{17, 17, 19},
		renderedTable:  `fixtures/table_6`,
	},
}

func TestTableStyle(t *testing.T) {
	headers := []string{"H1", "H2", "H3"}
	Convey("Create new style from data", t, func() {
		for idx, test := range testsTableStyle {
			Convey(fmt.Sprintf("%d)", idx), func() {
				table := NewTable(headers)
				style := NewDefaultTableStyle()
				table.style = style
				for _, row := range test.data {
					table.AddRow(row)
				}

				Convey("Total width with waste + col sizes must be equal to total render size", func() {
					widths := style.CalculateColWidths(table, 60)
					waste := style.Waste(table.colAmount)
					total := waste
					for _, w := range widths {
						total += w
					}
					So(waste, ShouldEqual, 7)
					So(total, ShouldEqual, 60)
					//fmt.Printf("EXP WIDTH: %v -- IS WIDTH: %v\n", test.expectedWidths, widths)
					So(test.expectedWidths, ShouldResemble, widths)
				})

				Convey("Render table", func() {
					style.HeaderRenderer = _testTableRenderHeader
					out := style.Render(table, 60)

					expect, _ := ioutil.ReadFile(test.renderedTable)
					//fmt.Printf("\n--IS: \n%s\n---\n-- EXPECT: \n%s\n--\n", out, expect)
					So(out, ShouldResemble, string(expect))
				})

				if test.renderedTable == "fixtures/table_5" {
					Convey("Render table with open style", func() {
						style = CopyTableStyle(OpenTableStyle)
						style.HeaderRenderer = _testTableRenderHeader
						out := style.Render(table, 60)

						expect, _ := ioutil.ReadFile(test.renderedTable + "_open")
						//fmt.Printf("\n--IS: \n%s\n---\n-- EXPECT: \n%s\n--\n", out, expect)
						So(out, ShouldResemble, string(expect))
					})
				}
			})
		}
	})
}
