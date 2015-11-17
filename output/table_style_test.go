package output

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"io/ioutil"
	"regexp"
)

var testsTableStyle = []struct {
	data           [][]string
	expectedWidths []uint
	renderedTable string
}{
	{
		data: [][]string{
			{"foo", "bar", "baz"},
		},
		expectedWidths: []uint{17, 17, 19},
		renderedTable: `fixtures/table_0`,
	},
	{
		data: [][]string{
			{"foofoofoofoofoo", "bar", "baz"},
		},
		expectedWidths: []uint{37, 7, 9},
		renderedTable: `fixtures/table_1`,
	},
	{
		data: [][]string{
			{"foo\nfoo\nfoo\nfoo\nfoo", "bar", "baz"},
		},
		expectedWidths: []uint{17, 17, 19},
		renderedTable: `fixtures/table_2`,
	},
	{
		data: [][]string{
			{"foofoofoofoofoo", "bar", "baz"},
			{"foo", "barbarbarbarbar", "baz"},
		},
		expectedWidths: []uint{23, 23, 7},
		renderedTable: `fixtures/table_3`,
	},
	{
		data: [][]string{
			{"foofoofoofoofoo", "bar", "baz"},
			{"foo", "barbarbarbarbar", "baz"},
			{"foo", "bar", "bazbazbazbazbaz"},
		},
		expectedWidths: []uint{17, 17, 19},
		renderedTable: `fixtures/table_4`,
	},
	{
		data: [][]string{
			{"foofoofoofoofoofoofoofoofoofoofoofoofoofoofoo", "bar", "baz"},
			{"foo", "barbarbarbarbarbarbarbarbarbarbarbarbarbarbar", "baz"},
			{"foo", "bar", "bazbazbazbazbazbazbazbazbazbazbazbazbazbazbaz"},
		},
		expectedWidths: []uint{17, 17, 19},
		renderedTable: `fixtures/table_5`,
	},
}

func TestTableStyle(t *testing.T) {
	headers := []string{"H1", "H2", "H3"}
	Convey("Create new style from data", t, func() {
		for idx, test := range testsTableStyle {
			Convey(fmt.Sprintf("%d)", idx), func() {
				table := NewTable(headers)
				style := NewDefaultTableStyle()
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
					So(test.expectedWidths, ShouldResemble, widths)
				})

				Convey("Render table", func() {
					style.HeaderRenderer = func(str string) string {
						rxRight := regexp.MustCompile(`^(.+\S+)(\s+)$`)
						if rxRight.MatchString(str) {
							match := rxRight.FindStringSubmatch(str)
							word := match[1]
							spaces := match[2]
							if len(spaces) > 2 {
								return fmt.Sprintf("*%s*%s", word, spaces[0:len(spaces)-2])
							}
						}
						return str
					}
					out := style.Render(table, 60)
					//fmt.Printf("\nALL:\n--\n%s\n--\n", out)

					expect, _ := ioutil.ReadFile(test.renderedTable)
					So(out, ShouldResemble, string(expect))
				})

				if idx == 0 {
					Convey("Render table without top & bottom", func() {
						style = NewDefaultTableStyle()
						style.Top = ""
						style.Bottom = ""
						style.HeaderRenderer = func(str string) string {
							rxRight := regexp.MustCompile(`^(.+\S+)(\s+)$`)
							if rxRight.MatchString(str) {
								match := rxRight.FindStringSubmatch(str)
								word := match[1]
								spaces := match[2]
								if len(spaces) > 2 {
									return fmt.Sprintf("*%s*%s", word, spaces[0:len(spaces)-2])
								}
							}
							return str
						}
						out := style.Render(table, 60)
						//fmt.Printf("\nALL:\n--\n%s\n--\n", out)

						expect, _ := ioutil.ReadFile(test.renderedTable+ "_light")
						So(out, ShouldResemble, string(expect))
					})
				}
			})
		}
	})
}
