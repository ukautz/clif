package clif

import (
	. "github.com/smartystreets/goconvey/convey"
	"regexp"
	"testing"
)

func TestOptionBuilder(t *testing.T) {
	Convey("Option builder methods", t, func() {
		opt := NewOption("foo", "f", "Do foo", "", false, false).
			SetUsage("Do foo foo").
			SetDescription("Do foo foo foo").
			SetDefault("bar").
			SetEnv("lala").
			SetParse(func(n, v string) (string, error) { return n + "1" + v, nil }).
			SetRegex(regexp.MustCompile(`^b`))
		So(opt.Name, ShouldEqual, "foo")
		So(opt.Usage, ShouldEqual, "Do foo foo")
		So(opt.Description, ShouldEqual, "Do foo foo foo")
		So(opt.Default, ShouldEqual, "bar")
		So(opt.Env, ShouldEqual, "lala")
		So(opt.Parse, ShouldNotBeNil)
		res, err := opt.Parse("x", "y")
		So(err, ShouldBeNil)
		So(res, ShouldEqual, "x1y")
		So(opt.Regex, ShouldResemble, regexp.MustCompile(`^b`))
	})
	Convey("Option as flag", t, func() {
		opt := NewOption("foo", "f", "Do foo", "", false, false).IsFlag()
		So(opt.Flag, ShouldBeTrue)
	})
}
