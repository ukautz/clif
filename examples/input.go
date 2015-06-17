package main

import (
	"fmt"
	"github.com/ukautz/clif"
)

func callIn(in clif.Input, out clif.Output) {
	name := in.Ask("Who are you? ", func(v string) error {
		if len(v) > 0 {
			return nil
		} else {
			return fmt.Errorf("Didn't catch that")
		}
	})
	father := in.Choose("Who is your father?", map[string]string{
		"yoda":  "The small, green guy",
		"darth": "NOOOOOOOO!",
		"obi":   "The old man with the light thingy",
	})

	out.Printf("Well, %s, ", name)
	if father != "darth" {
		out.Printf("<success>may the force be with you!<reset>\n")
	} else {
		out.Printf("<error>u bad!<reset>\n")
	}
}

func main() {
	clif.New("inputter", "0.1.0", "Input example").
		New("in", "Test input", callIn).
		Run()
}
