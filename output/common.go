package output

import (
	"unicode/utf8"
	"regexp"
)

const (
	LINE_BREAK = "\n"
)

var (
	rxControlCharacters = regexp.MustCompile(`\033\[[\d;]+m`)
)

func StringLength(str string) int {
	str = rxControlCharacters.ReplaceAllString(str, "")
	return utf8.RuneCountInString(str)
}
