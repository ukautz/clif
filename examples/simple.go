// +build never

package main

import "github.com/ukautz/clif"

func main() {
	c := clif.New("My App", "1.0.0", "An example application").
		New("hello", "The obligatory hello world", func(out clif.Output) {
		out.Printf("Hello World\n")
	})
	c.Run()
}
