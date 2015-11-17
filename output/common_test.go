package output

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
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
}

func TestStringLength(t *testing.T) {
	Convey("Length of strings", t, func() {
		for idx, test := range testsStringLenth {
			Convey(fmt.Sprintf("%d)", idx), func() {
				l := StringLength(test.str)
				So(l, ShouldEqual, test.len)
			})
		}
	})
}
