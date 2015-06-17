package clif

import (
	"fmt"
	"os"
)

var Die = func(msg string, args ...interface{}) {
	fmt.Printf(msg+"\n", args...)
	os.Exit(1)
}
