package clif

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"unicode/utf8"
)

const (
	LINE_BREAK = "\n"
)

var (
	rxControlCharacters = regexp.MustCompile(`\033\[[\d;]+m`)
)

// StringLength returns the length of an UTF8 string without any control characters
func StringLength(str string) int {
	str = rxControlCharacters.ReplaceAllString(str, "")
	return utf8.RuneCountInString(str)
}

// SplitFormattedString splits formatted string into multiple lines while making
// sure that control characters end at line end and possibly re-start at next line
func SplitFormattedString(str string) []string {
	chars := []byte(str)
	lastIdx := len(chars) - 1
	seq := 0
	has := false
	cache := []byte{}
	current := []byte{}
	result := []byte{}
	for idx, c := range chars {
		add := []byte{c}
		if seq == 0 && c == 27 { // \033
			seq++
			cache = append(cache, c)
			current = []byte{}
		} else if seq == 1 && c == 91 { // "["
			seq++
			cache = append(cache, 91)
		} else if seq == 2 && ((c >= '0' && c <= '9') || c == ';') {
			seq = 3
			cache = append(cache, c)
			current = append(current, c)
		} else if seq == 3 && ((c >= '0' && c <= '9') || c == ';' || c == 'm') { // "m"
			cache = append(cache, c)
			if c == 'm' {
				seq = 0
				l := len(current)
				if current[l-1] == '0' {
					cache = []byte{}
					current = []byte{}
					has = false
				} else {
					has = true
				}
			} else {
				current = append(current, c)
			}
		} else if c == '\n' && has && idx != lastIdx {
			if false {
				fmt.Printf(" -> PREPEND CACHE\n")
			}
			add = []byte("\033[0m")
			add = append(add, c)
			add = append(add, cache...)
		} else if has && idx == lastIdx {
			add = []byte{c}
			add = append(add, []byte("\033[0m")...)
		}
		result = append(result, add...)
		if false {
			fmt.Printf("%d) C = %d (%c)\n", idx, c, c)
		}
	}
	return strings.Split(string(result), "\n")
}

// Die is the default function executed on die. It can be used as a shorthand
// via `clif.Die("foo %s", "bar")` and can be overwritten to change the failure
// exit handling CLI-wide.
var Die = func(msg string, args ...interface{}) {
	NewColorOutput(os.Stderr).Printf("<error>"+msg+"<reset>\n", args...)
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
