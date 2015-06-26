package clif

import (
	"fmt"
	"os"
)

// Die is the default function executed on die. It can be used as a shorthand
// via `clif.Die("foo %s", "bar")` and can be overwritten to change the failure
// exit handling CLI-wide.
var Die = func(msg string, args ...interface{}) {
	fmt.Printf(msg+"\n", args...)
	os.Exit(1)
}
