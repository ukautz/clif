// +build ignore

package main

import (
	"fmt"
	"gopkg.in/ukautz/clif.v0"
)

func callIn(in clif.Input, out clif.Output) {
	name := in.Ask("Who are you?", func(v string) error {
		if len(v) > 0 {
			return nil
		} else {
			return fmt.Errorf("Didn't catch that")
		}
	})
	out.Printf("\n")
	father := ""
	for {
		father = in.Choose(fmt.Sprintf("Hello %s. Who is your father?", name), map[string]string{
			"yoda":  "The small, green guy",
			"darth": "The one with the dark cloark and a smokers voice",
			"obi":   "The old man with the light thingy",
		})
		if in.Confirm("You're sure about that? (y/n)") {
			break
		} else {
			out.Printf("\n")
		}
	}

	out.Printf("Well, <important>%s<reset>, ", name)
	if father != "darth" {
		out.Printf("<success>may the force be with you!<reset>\n")
	} else {
		out.Printf("<error>NOOOOOOOO!<reset>\n")
	}
}

func main() {
	clif.New("inputter", "0.1.0", "Input example").
		New("in", "Test input", callIn).
		Run()
}
