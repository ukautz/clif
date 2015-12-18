package clif

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"unicode/utf8"
	"io"
	"path/filepath"
	"time"
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
		if seq == 0 && IsControlCharStart(c) { // \033
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

// IsControlCharStart returns bool whether character is "\033" aka "\e"
func IsControlCharStart(c byte) bool {
	return c == 27
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

var dbgFh io.Writer
var Dbg = func(msg string, args ...interface{}) {
	if v := os.Getenv("DEBUG_CLIF"); v != "1" {
		return
	}
	if dbgFh == nil {
		dbgFh, _ = os.OpenFile(filepath.Join(os.TempDir(), "debug.clif"), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	}
	dbgFh.Write([]byte(fmt.Sprintf("[%s] %s\n", time.Now(), fmt.Sprintf(msg, args...))))
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

func _stringTableCompare(s1, s2 string) {
	out := NewTable([]string{"IS", "SHOULD"})
	l1 := len(s1)
	l2 := len(s2)
	max := l1
	if max < l2 {
		max = l2
	}
	for i := 0; i < max; i++ {
		row := []string{"", ""}
		if i < l1 {
			row[0] = s1[i : i+1]
		}
		if i < l2 {
			row[1] = s2[i : i+1]
		}
		for j, c := range row {
			//fmt.Printf(" .. %d\n", j)
			if c == "\n" {
				row[j] = "<BR>"
			} else if c == "" {
				row[j] = "-"
			} else {
				row[j] = fmt.Sprintf("%d (%c)", c[0], c[0])
			}
		}
		out.AddRow(row)
		//fmt.Printf("> R %d: %v\n", i, row)
	}

	style := NewDefaultTableStyle()
	fmt.Println(style.Render(out, 30))
}

func _stringRenderDump(s string) string {
	s = strings.Replace(s, "\n", "\\n", -1)
	s = strings.Replace(s, "\t", "→", -1)
	s = strings.Map(func(r rune) rune {
		if r == ' ' || r == '\t' {
			return '┈'
		} else if r >= 32 && r < 127 {
			return r
		}
		return '჻'
	}, s)
	return s
}

func _stringCompareDump(s1, s2 string) {
	fmt.Printf("\n-------------------------------------\n")
	fmt.Printf("IS:     \"%s\"\n", _stringRenderDump(s1))
	fmt.Printf("--------\n")
	fmt.Printf("SHOULD: \"%s\"\n", _stringRenderDump(s2))
	fmt.Printf("-------------------------------------\n")
}
