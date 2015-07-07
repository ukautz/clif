package clif

import (
	. "github.com/smartystreets/goconvey/convey"
	"regexp"
	"testing"
)

func TestArgumentBuilder(t *testing.T) {
	Convey("Argument builder methods", t, func() {
		arg := NewArgument("foo", "Do foo", "", false, false).
			SetUsage("Do foo foo").
			SetDescription("Do foo foo foo").
			SetDefault("bar").
			SetEnv("lala").
			SetParse(func(n, v string) (string, error) { return n + "1" + v, nil }).
			SetRegex(regexp.MustCompile(`^b`))
		So(arg.Name, ShouldEqual, "foo")
		So(arg.Usage, ShouldEqual, "Do foo foo")
		So(arg.Description, ShouldEqual, "Do foo foo foo")
		So(arg.Default, ShouldEqual, "bar")
		So(arg.Env, ShouldEqual, "lala")
		So(arg.Parse, ShouldNotBeNil)
		res, err := arg.Parse("x", "y")
		So(err, ShouldBeNil)
		So(res, ShouldEqual, "x1y")
		So(arg.Regex, ShouldResemble, regexp.MustCompile(`^b`))
	})
}
