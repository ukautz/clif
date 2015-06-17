package clif

import (
	. "github.com/smartystreets/goconvey/convey"
	"reflect"
	"testing"
)

type testFoo interface {
	Bar() string
}

type testBar struct{}

func (this *testBar) Bar() string {
	return "baz"
}

func TestRegistryFull(t *testing.T) {
	Convey("With new registry", t, func() {
		reg := NewRegistry()

		Convey("Not registered, not found", func() {
			So(reg.Has("foo"), ShouldEqual, false)
		})
		Convey("Is registered, is found", func() {
			v := new(testBar)
			reg.Register(v)
			So(reg.Has("*clif.testBar"), ShouldEqual, true)
			So(reg.Has(reflect.TypeOf(v).String()), ShouldEqual, true)
			So(reg.Get("*clif.testBar").Interface(), ShouldEqual, v)
		})
		Convey("Is aliased & registered, is found", func() {
			v := new(testBar)
			a := reflect.TypeOf((*testFoo)(nil)).Elem()
			reg.Alias(a.String(), v)
			So(reg.Has("clif.testFoo"), ShouldEqual, true)
			So(reg.Get("clif.testFoo").Interface(), ShouldEqual, v)
		})
	})
}
