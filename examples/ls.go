// +build ignore

package main

import (
	"gopkg.in/ukautz/clif.v1"
	"os/exec"
	"fmt"
)

func main() {
	clif.New("My App", "1.0.0", "An example application").
		New("ls", "", func() { out, _ := exec.Command("ls", "-lha").Output(); fmt.Println(string(out)) }).
		New("ps", "", func() { out, _ := exec.Command("ps", "-auxf").Output(); fmt.Println(string(out)) }).
		Run()
}
