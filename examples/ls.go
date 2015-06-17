// +build ignore

package main

import (
	"gopkg.in/ukautz/clif.v0"
	"os/exec"
)

func main() {
	clif.New("My App", "1.0.0", "An example application").
		New("ls", "", func() { exec.Command("ls", "-lha", ".").Output() }).
		New("ps", "", func() { exec.Command("ps", "auxf") }).
		Run()
}
