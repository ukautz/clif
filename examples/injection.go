// +build ignore

package main

import (
	"fmt"
	"gopkg.in/ukautz/clif.v1"
	"reflect"
	"regexp"
	"strconv"
)

type (
	MyFoo interface {
		Bar() string
		SetBar(i int)
	}
	MyBar struct {
		bar int
	}
	MyBaz struct {
		baz string
	}
)

func (this *MyBar) Bar() string {
	return fmt.Sprintf("Bar is %d", this.bar)
}

func (this *MyBar) SetBar(i int) {
	this.bar = i
}

func (this *MyBaz) String() string {
	return fmt.Sprintf("~~ %s ~~", this.baz)
}

func callMe(out clif.Output, in clif.Input, c *clif.Command, foo MyFoo, baz *MyBaz) {
	barIn := in.Ask("Gimme a bar integer: ", func(v string) error {
		_, err := strconv.Atoi(v)
		return err
	})
	barInt, _ := strconv.Atoi(barIn)
	foo.SetBar(barInt)

	bazIn := in.AskRegex("Now please a baz: ", regexp.MustCompile(`^B`))
	baz.baz = bazIn

	out.Printf("Bar: <info>%s<reset>\nBaz: <headline>%s<reset>\n", foo.Bar(), baz)
}

func main() {
	cli := clif.New("my-app", "My kewl App", "0.8.5")
	cmd := clif.NewCommand("call", "Call me", callMe)
	cli.Add(cmd)
	cli.Register(new(MyBaz)).
		RegisterAs(reflect.TypeOf((*MyFoo)(nil)).Elem().String(), new(MyBar))
	cli.Run()
}
