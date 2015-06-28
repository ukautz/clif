package clif

import (
	. "github.com/smartystreets/goconvey/convey"
	"reflect"
	"testing"
	"fmt"
)

type testFoo interface {
	Bar() string
}

type testBar struct{}

func (this *testBar) Bar() string {
	return "baz"
}

func TestRegistrySettingAndGetting(t *testing.T) {
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
			So(reg.Names(), ShouldResemble, []string{"*clif.testBar"})
		})
		Convey("Is aliased & registered, is found", func() {
			v := new(testBar)
			a := reflect.TypeOf((*testFoo)(nil)).Elem()
			reg.Alias(a.String(), v)
			So(reg.Has("clif.testFoo"), ShouldEqual, true)
			So(reg.Get("clif.testFoo").Interface(), ShouldEqual, v)
			So(reg.Names(), ShouldResemble, []string{"clif.testFoo"})
		})
	})
}


func TestRegistryReduce(t *testing.T) {
	Convey("With new registry", t, func() {
		reg := NewRegistry()
		names := []string{}
		values := []interface{}{}
		for i := 0; i < 100; i++ {
			name := fmt.Sprintf("V:%03d", i)
			reg.Alias(name, i)
			names = append(names, name)
			values = append(values, i)
		}
		reg.Register(new(testBar))

		Convey("Reducing synced", func() {
			found := reg.Reduce(func(name string, value interface{}) bool {
				if name == "V:098" {
					return true
				} else if v, ok := value.(int); ok && v == 17 {
					return true
				} else {
					return false
				}
			})
			So(found, ShouldResemble, []interface{}{17, 98})

			Convey("Reducing synced works in order", func() {
				orderNames := []string{}
				orderValues := reg.Reduce(func(name string, value interface{}) bool {
					if _, ok := value.(int); ok {
						orderNames = append(orderNames, name)
						return true
					} else {
						return false
					}
				})
				So(orderNames, ShouldResemble, names)
				So(orderValues, ShouldResemble, values)
			})
		})

		Convey("Reducing asynchron", func() {
			res := reg.ReduceAsync(func(name string, value interface{}) bool {
				if name == "V:098" {
					return true
				} else if v, ok := value.(int); ok && v == 17 {
					return true
				} else {
					return false
				}
			})
			found := 0
			fail := 0
			sum := 0
			for v := range res {
				if i, ok := v.(int); ok && (i == 17 || i == 98) {
					found ++
					sum += i
				} else {
					fail ++
				}
			}
			So(found, ShouldEqual, 2)
			So(sum, ShouldEqual, 115)
			So(fail, ShouldEqual, 0)
		})

	})
}