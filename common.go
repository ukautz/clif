package clif

import (
	"os"
)

// Die is the default function executed on die. It can be used as a shorthand
// via `clif.Die("foo %s", "bar")` and can be overwritten to change the failure
// exit handling CLI-wide.
var Die = func(msg string, args ...interface{}) {
	NewColorOutput(os.Stderr).Printf("<error>"+ msg+"<reset>\n", args...)
	Exit(1)
}

// Exit is wrapper for os.Exit, so it can be overwritten for tests or edge use cases
var Exit = func(s int) {
	os.Exit(s)
}

// CommandSort implements the `sort.Sortable` interface for commands, based on
// the command `Name` attribute
type CommandsSort []*Command

func (this CommandsSort) Len() int {
	return len(this)
}

func (this CommandsSort) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

func (this CommandsSort) Less(i, j int) bool {
	return this[i].Name < this[j].Name
}