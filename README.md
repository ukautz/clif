Go CLI framework
================

Go framework to rapidly develop command line applications.

Example
-------

![Example]https://cloud.githubusercontent.com/assets/600604/7931252/36ed6694-090f-11e5-9e38-54302fe98efc.gif!

```go
package main

import "github.com/ukautz/go-cli"

func main() {
	c := cli.New("My App", "1.0.0", "An example application").
		New("hello", "The obligatory hello world", func(out cli.Output) {
		out.Printf("Hello World\n")
	})
	c.Run()
}
```

Until readme is written: please see examples folder.