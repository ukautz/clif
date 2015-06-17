// +build ignore

package main

import (
	"github.com/ukautz/clif"
	"os/exec"
)

func main() {
	clif.New("My App", "1.0.0", "An example application").
		New("ls", "", func() { exec.Command("ls", "-lha", ".").Output() }).
		New("ps", "", func() { exec.Command("ps", "auxf") }).
		Run()
}
