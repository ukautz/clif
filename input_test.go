package clif

import (
	"bytes"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"regexp"
	"sync"
	"testing"
)

func TestDefaultInput(t *testing.T) {
	Convey("Input data", t, func() {
		bufIn := bytes.NewBuffer(nil)
		bufOut := bytes.NewBuffer(nil)
		out := NewMonochromeOutput(bufOut)
		in := NewDefaultInput(bufIn, out)

		Convey("Ask returns on any non-empty input", func() {
			var wg sync.WaitGroup
			wg.Add(2)
			go func() {
				defer wg.Done()
				bufIn.WriteString("Foo\n")
			}()
			res := ""
			go func() {
				defer wg.Done()
				res = in.Ask("Foo? ", nil)
			}()
			wg.Wait()
			So(res, ShouldEqual, "Foo")
			So(bufOut.String(), ShouldEqual, "Foo? ")
		})

		Convey("Ask with check tries until check ok", func() {
			var wg sync.WaitGroup
			var wgFirst sync.WaitGroup
			wgFirst.Add(1)
			wg.Add(3)
			go func() {
				defer wg.Done()
				defer wgFirst.Done()
				bufIn.WriteString("Foo\nBaz\n")
			}()
			go func() {
				defer wg.Done()
				wgFirst.Wait()
				bufIn.WriteString("Bar\n")
			}()
			res := ""
			go func() {
				defer wg.Done()
				res = in.Ask("Foo? ", func(c string) error {
					if c == "Bar" {
						return nil
					} else {
						return fmt.Errorf("Not Bar!")
					}
				})
			}()
			wg.Wait()
			So(res, ShouldEqual, "Bar")
			So(bufOut.String(), ShouldEqual, `Foo? Not Bar!

Foo? Not Bar!

Foo? `)
		})

		Convey("Ask with regular expression tries until matches", func() {
			var wg sync.WaitGroup
			var wgFirst sync.WaitGroup
			wgFirst.Add(1)
			wg.Add(3)
			rx := regexp.MustCompile(`Bar`)
			go func() {
				defer wg.Done()
				defer wgFirst.Done()
				bufIn.WriteString("Foo\nBaz\n")
			}()
			go func() {
				defer wg.Done()
				wgFirst.Wait()
				bufIn.WriteString("Bar\n")
			}()
			res := ""
			go func() {
				defer wg.Done()
				res = in.AskRegex("Foo? ", rx)
			}()
			wg.Wait()
			So(res, ShouldEqual, "Bar")
			So(bufOut.String(), ShouldEqual, `Foo? Input does not match criteria

Foo? Input does not match criteria

Foo? `)
		})

		Convey("Choose presents options and returns on valid choice", func() {
			var wg sync.WaitGroup
			var wgFirst sync.WaitGroup
			wgFirst.Add(1)
			wg.Add(3)
			go func() {
				defer wg.Done()
				defer wgFirst.Done()
				bufIn.WriteString("Foo\nBaz\n")
			}()
			go func() {
				defer wg.Done()
				wgFirst.Wait()
				bufIn.WriteString("the bar\n")
			}()
			res := ""
			go func() {
				defer wg.Done()
				res = in.Choose("Choose or loose!", map[string]string{
					"foo":     "Foo!!!",
					"the bar": "One bar please",
					"42":      "Take that",
				})
			}()
			wg.Wait()
			So(res, ShouldEqual, "the bar")
			So(bufOut.String(), ShouldEqual, `Choose or loose!
  42)      Take that
  foo)     Foo!!!
  the bar) One bar please
Choose: Choose one of: 42, foo, the bar

Choose or loose!
  42)      Take that
  foo)     Foo!!!
  the bar) One bar please
Choose: Choose one of: 42, foo, the bar

Choose or loose!
  42)      Take that
  foo)     Foo!!!
  the bar) One bar please
Choose: `)
		})
	})
}
