package clif

import (
	"bytes"
	. "github.com/smartystreets/goconvey/convey"
	"strings"
	"testing"
)

type testCustomFormatter struct{}

func (this *testCustomFormatter) Format(msg string) string {
	return strings.ToUpper(msg)
}

func (this *testCustomFormatter) Escape(msg string) string {
	return msg
}

func TestOutput(t *testing.T) {
	Convey("Default output rendering", t, func() {
		b := bytes.NewBuffer(nil)
		o := NewMonochromeOutput(b)
		o.Printf("With <headline>formatted<reset> input")
		So(b.String(), ShouldEqual, "With formatted input")
	})
	Convey("Fancy output rendering", t, func() {
		b := bytes.NewBuffer(nil)
		o := NewColorOutput(b)
		o.Printf("With <headline>formatted<reset> input")
		So(b.String(), ShouldEqual, "With \033[4;1mformatted\033[0m input")
	})
	Convey("Custom output rendering", t, func() {
		b := bytes.NewBuffer(nil)
		o := NewOutput(b, &testCustomFormatter{})
		o.Printf("With <headline>formatted<reset> input")
		So(b.String(), ShouldEqual, "WITH <HEADLINE>FORMATTED<RESET> INPUT")
	})
	Convey("Switching formatter later on", t, func() {
		b := bytes.NewBuffer(nil)
		o := NewMonochromeOutput(b)
		o.SetFormatter(&testCustomFormatter{})
		o.Printf("With <headline>formatted<reset> input")
		So(b.String(), ShouldEqual, "WITH <HEADLINE>FORMATTED<RESET> INPUT")
	})
	Convey("Multi-line rendering", t, func() {
		b := bytes.NewBuffer(nil)
		o := NewColorOutput(b)
		o.Printf("With <headline>formatted\nfoo\nbar\nbaz<reset> input")
		s := o.Sprintf("With <headline>formatted\nfoo\nbar\nbaz<reset> input")
		So(b.String(), ShouldEqual, "With \033[4;1mformatted\nfoo\nbar\nbaz\033[0m input")
		So(b.String(), ShouldEqual, s)
	})
}
