package clif

import (
	"os"
	"reflect"
)

// Die is the default function executed on die. It can be used as a shorthand
// via `clif.Die("foo %s", "bar")` and can be overwritten to change the failure
// exit handling CLI-wide.
var Die = func(msg string, args ...interface{}) {
	NewColorOutput(os.Stderr).Printf("<error>"+ msg+"<reset>\n", args...)
	Exit(1)
}

var Exit = func(s int) {
	os.Exit(s)
}

func clone(v interface{}) interface{} {
	return reflect.ValueOf(v).Elem().Addr().Interface()
}