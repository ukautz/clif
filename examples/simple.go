// +build never

package main

import "github.com/ukautz/go-cli"

func main() {
	c := cli.New("My App", "1.0.0", "An example application").
		New("hello", "The obligatory hello world", func(out cli.Output) {
		out.Printf("Hello World\n")
	})
	c.Run()
}
