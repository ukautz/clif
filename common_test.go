package clif

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"strings"
	"testing"
)

var testsStringLenth = []struct {
	str string
	len int
}{
	{
		str: "foo",
		len: 3,
	},
	{
		str: "♞♞♞",
		len: 3,
	},
	{
		str: "\033[1mfoo\033[0m",
		len: 3,
	},
	{
		str: "\033[34m\033[1mfoo\033[0m",
		len: 3,
	},
}

func TestStringLength(t *testing.T) {
	Convey("Length of strings", t, func() {
		for idx, test := range testsStringLenth {
			name := rxControlCharacters.ReplaceAllString(test.str, "<CTRL>")
			Convey(fmt.Sprintf("%d) %s", idx, name), func() {
				l := StringLength(test.str)
				So(l, ShouldEqual, test.len)
			})
		}
	})
}

var testsSplitFormattedString = []struct {
	str    string
	expect string
}{
	{
		str:    "fffff",
		expect: "fffff",
	},
	{
		str:    "\033[1mfffff\033[0m",
		expect: "\033[1mfffff\033[0m",
	},
	{
		str:    "\033[1mfffff\033[0m\n",
		expect: "\033[1mfffff\033[0m\n",
	},
	{
		str:    "\033[1mfff\nfff\033[0m",
		expect: "\033[1mfff\033[0m\n\033[1mfff\033[0m",
	},
	{
		str: strings.Join([]string{
			"\033[34;1mfff",
			"fff",
			"fff\033[0m",
		}, "\n"),
		expect: strings.Join([]string{
			"\033[34;1mfff\033[0m",
			"\033[34;1mfff\033[0m",
			"\033[34;1mfff\033[0m",
		}, "\n"),
	},
	{
		str: strings.Join([]string{
			"\033[34;1mfff",
			"fff",
			"fff",
			"fff",
			"fff\033[0m",
		}, "\n"),
		expect: strings.Join([]string{
			"\033[34;1mfff\033[0m",
			"\033[34;1mfff\033[0m",
			"\033[34;1mfff\033[0m",
			"\033[34;1mfff\033[0m",
			"\033[34;1mfff\033[0m",
		}, "\n"),
	},
	{
		str: strings.Join([]string{
			"\033[34;1mfff\033[35m",
			"fff",
			"fff\033[0m",
		}, "\n"),
		expect: strings.Join([]string{
			"\033[34;1mfff\033[35m\033[0m",
			"\033[34;1m\033[35mfff\033[0m",
			"\033[34;1m\033[35mfff\033[0m",
		}, "\n"),
	},
	{
		str: strings.Join([]string{
			"\033[34;1mfff",
			"\033[35mfff",
			"fff\033[0m",
		}, "\n"),
		expect: strings.Join([]string{
			"\033[34;1mfff\033[0m",
			"\033[34;1m\033[35mfff\033[0m",
			"\033[34;1m\033[35mfff\033[0m",
		}, "\n"),
	},
	{
		str:    "\033[34;1mfff",
		expect: "\033[34;1mfff\033[0m",
	},
	{
		str:    "fff\nfff",
		expect: "fff\nfff",
	},
}

func TestSplitFormattedString(t *testing.T) {
	Convey("Split lines with control characters", t, func() {
		for idx, test := range testsSplitFormattedString {
			from := strings.Replace(test.str, "\n", "<BR>", -1)
			from = rxControlCharacters.ReplaceAllString(from, "<CTRL>")
			Convey(fmt.Sprintf("%d) from \"%s\"", idx, from), func() {
				to := strings.Join(SplitFormattedString(test.str), "\n")
				//_testDumpStrings(to, test.expect)
				So(to, ShouldEqual, test.expect)
			})
		}
	})
}
