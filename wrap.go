package clif

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"
	"unicode"
)

type (
	// WrapTrimMode defines whether and how wrapped lines are trimmed. See WRAP_TRIM_*
	WrapTrimMode int

	// WrapWhitespaceMode defines whether in-sentence whitespaces are reduced contracted or not. See WRAP_WHITESPACE_*
	WrapWhitespaceMode int

	// Wrapper transforms multi- or single-line strings into multi-line strings of maximum line length.
	Wrapper struct {

		// BreakWords controls whether over-length words can be broken or not.
		// Example: Wrapping "foo barrr" with limit of 3 and NO wrapping will
		// return "foo\nbarrr" and with wrapping will return "foo\nbar\nrr"
		BreakWords bool

		// KeepEmptyLines controls whether empty lines are removed or not. Defaults
		// to false. If set to true then "foo\n\nbar" stays "foo\n\n"bar". Otherwise
		// it becomes "foo\nbar"
		KeepEmptyLines bool

		// Limit is the max length per wrapped line
		Limit uint

		// TrimMode defines whether/how wrapped lines are trimmed or not
		TrimMode WrapTrimMode

		// WhitespaceMode defines whether in-sentence whitespaces are reduced
		// or not. Eg with disabled contracting
		WhitespaceMode WrapWhitespaceMode
	}
)

const (
	// WRAP_TRIM_NONE keeps "  foo bar baz  " as is
	WRAP_TRIM_NONE WrapTrimMode = iota

	// WRAP_TRIM_RIGHT transforms "  foo bar baz  " to "foo bar baz  "
	WRAP_TRIM_RIGHT

	// WRAP_TRIM_LEFT transforms "  foo bar baz  " to "  foo bar baz"
	WRAP_TRIM_LEFT

	// WRAP_TRIM_BOTH transforms "  foo bar baz  " to "foo bar baz"
	WRAP_TRIM_BOTH
)

const (
	// WRAP_WHITESPACE_CONTRACT contracts "foo   bar   baz" to "foo bar baz"
	WRAP_WHITESPACE_CONTRACT WrapWhitespaceMode = iota

	// WRAP_WHITESPACE_KEEP keeps "foo   bar   baz" as is
	WRAP_WHITESPACE_KEEP
)

func NewWrapper(limit uint) *Wrapper {
	return &Wrapper{
		Limit:          limit,
		TrimMode:       WRAP_TRIM_RIGHT,
		WhitespaceMode: WRAP_WHITESPACE_CONTRACT,
		KeepEmptyLines: false,
	}
}

func Wrap(s string, limit uint) string {
	return NewWrapper(limit).Wrap(s)
}

// WrapString wraps the given string within lim width in characters.
// Code is partially stolen from https://raw.githubusercontent.com/mitchellh/go-wordwrap/master/wordwrap.go
func (this *Wrapper) Wrap(s string) string {
	lines := []string{""}
	wordBuf := ""
	wordBufLen := uint(0)
	curLineNum := 0
	curLineLen := uint(0)
	controlCharSeq := 0
	lastControlChars := []rune{}
	controlCharBuf := bytes.NewBuffer(nil)
	controlCharNoneEndBuf := bytes.NewBuffer(nil)
	var lastChar rune

	rxReplaceLeftBeforeCtrlChars := regexp.MustCompile(`^\s+\033`)
	rxReplaceRightBeforeCtrlChars := regexp.MustCompile(`\s+(\033\[[0-9]+(?:;[0-9]+)*m)$`)
	trimCurrent := func() {
		switch this.TrimMode {
		case WRAP_TRIM_RIGHT:
			lines[curLineNum] = strings.TrimRight(lines[curLineNum], " \t")
			lines[curLineNum] = rxReplaceRightBeforeCtrlChars.ReplaceAllStringFunc(lines[curLineNum], func(in string) string {
				return strings.TrimLeft(in, " \t")
			})
		case WRAP_TRIM_LEFT:
			lines[curLineNum] = rxReplaceLeftBeforeCtrlChars.ReplaceAllString(lines[curLineNum], "\033")
			lines[curLineNum] = strings.TrimLeft(lines[curLineNum], " \t")
		case WRAP_TRIM_BOTH:
			lines[curLineNum] = strings.TrimSpace(lines[curLineNum])
			lines[curLineNum] = rxReplaceLeftBeforeCtrlChars.ReplaceAllString(lines[curLineNum], "\033")
			lines[curLineNum] = rxReplaceRightBeforeCtrlChars.ReplaceAllStringFunc(lines[curLineNum], func(in string) string {
				return strings.TrimLeft(in, " \t")
			})
		}
	}

	finishLine := func(add string) {
		hasControlChars := controlCharNoneEndBuf.Len() > 0
		_wrapDebug("++ FINISH LINE (%v)\n", hasControlChars)
		lines[curLineNum] += add
		if hasControlChars {
			lines[curLineNum] += "\033[0m"
		}
		trimCurrent()
		lines = append(lines, "")
		curLineNum++
		if hasControlChars {
			lines[curLineNum] += controlCharNoneEndBuf.String()
		}
		curLineLen = 0
	}

	ctrlCharsEnded := func() bool {
		l := len(lastControlChars)
		if l == 1 && lastControlChars[0] == '0' {
			return true
		} else {
			for i := 0; i < l-1; i++ {
				if lastControlChars[i] == ';' && lastControlChars[i+1] == '0' && (i+2 <= l || lastControlChars[i+2] == ';') {
					return true
				}
			}
		}
		return false
	}

	for _, char := range s {
		if IsControlCharStart(byte(char)) {
			_wrapDebug(">> INIT CTRL CHARS\n")
			if wordBufLen > 0 {
				_wrapDebug("  >> PREPEND WORD BUF \"%s\"\n", _stringRenderDump(wordBuf))
				curLineLen += wordBufLen
				lines[curLineNum] += wordBuf
				wordBufLen = 0
				wordBuf = ""
			} else {
				_wrapDebug("  >> NO WORD BUF\n")
			}
			controlCharBuf.WriteRune(char)
			controlCharSeq = 1
		} else if controlCharSeq == 1 { // expect "["
			if char != 91 { // abort .. not "["
				_wrapDebug(">> ABORT CTRL CHARS\n")
				lines[curLineNum] += controlCharBuf.String()
				controlCharBuf.Reset()
				controlCharSeq = 0
			} else {
				_wrapDebug(">> START CTRL CHARS\n")
				controlCharBuf.WriteRune(char)
				controlCharSeq = 2
			}
		} else if controlCharSeq == 2 {
			if char >= '0' && char <= '9' {
				_wrapDebug(">> CONTINUE CTRL CHARS\n")
				controlCharBuf.WriteRune(char)
				lastControlChars = []rune{char}
				controlCharSeq = 3
			} else { // abort, "not valid char
				_wrapDebug(">> ABORT CTRL CHARS 2(%c)\n", char)
				lines[curLineNum] += controlCharBuf.String()
				controlCharBuf.Reset()
				controlCharSeq = 0
			}
		} else if controlCharSeq == 3 {
			if char >= '0' && char <= '9' {
				_wrapDebug(">> CONTINUE CTRL CHARS 2\n")
				lastControlChars = append(lastControlChars, char)
				controlCharBuf.WriteRune(char)
			} else if char == ';' {
				_wrapDebug(">> SEP CTRL CHARS\n")
				lastControlChars = []rune{}
				controlCharBuf.WriteRune(char)
				controlCharSeq = 2
			} else if char == 'm' { // end
				_wrapDebug(">> END CTRL CHARS\n")
				controlCharBuf.WriteRune(char)
				controlCharSeq = 0
				_wrapDebug(">> ADD CTRL CHARS: %s\n", _stringRenderDump(controlCharBuf.String()))
				lines[curLineNum] += controlCharBuf.String()
				if ctrlCharsEnded() {
					controlCharNoneEndBuf.Reset()
				} else {
					controlCharNoneEndBuf.WriteString(controlCharBuf.String())
				}
				controlCharBuf.Reset()
				/*if len(lastControlChars) == 1 && lastControlChars[len(lastControlChars)-1] == '0' {
					controlCharBuf.Reset()
				} else {
					_wrapDebug(" >> NO RESET CTRL CHARS CAUSE: %c\n", lastControlChars[len(lastControlChars)-1])
				}*/
				lastControlChars = []rune{}
				controlCharSeq = 0
			} else { // abort, "not valid char
				_wrapDebug(">> ABORT CTRL CHARS 3\n")
				lines[curLineNum] += controlCharBuf.String()
				controlCharBuf.Reset()
				controlCharSeq = 0
			}
		} else {
			if char == '\n' {
				_wrapDebug(">> ADD BREAK\n")
				finishLine(wordBuf)
				wordBuf = ""
				wordBufLen = 0
			} else if unicode.IsSpace(char) {
				_wrapDebug(">> ADD SPACE\n")
				if wordBufLen > 0 || this.WhitespaceMode == WRAP_WHITESPACE_KEEP {
					_wrapDebug(" >> SPACE WITH WORD OR KEEP\n")
					lines[curLineNum] += wordBuf
					curLineLen += wordBufLen
					lines[curLineNum] += string(char)
					curLineLen++
					wordBuf = ""
					wordBufLen = 0
					if curLineLen == this.Limit {
						_wrapDebug("  >> ADD SPACE FINISH LINE\n")
						finishLine("")
					}
				} else if curLineLen > 0 && unicode.IsSpace(lastChar) {
					_wrapDebug("\n>> HERE IN CONTINUE\n\n")
				} else {
					_wrapDebug("\n>> HERE IN SPACE ADD\n\n")
					lines[curLineNum] += string(char)
					curLineLen++
				}
			} else {
				_wrapDebug(">> ADD CHAR '%c'\n", char)
				totalLineLen := curLineLen + wordBufLen
				if totalLineLen+1 > this.Limit {
					if curLineLen > 0 { // has prefix before current word
						_wrapDebug("\n>> FINISH LINE WITH WORDBUF \"%s\"\n", wordBuf)
						finishLine("")
						wordBuf += string(char)
						wordBufLen++

					} else { // the word itself is longer than line
						_wrapDebug("\n>> WORD IS BIGGER \"%s\" (%v)\n", wordBuf, this.BreakWords)
						if this.BreakWords {
							finishLine(wordBuf)
							wordBuf = string(char)
							wordBufLen = 1
						} else {
							wordBuf += string(char)
							wordBufLen++
						}
					}
				} else {
					wordBuf += string(char)
					wordBufLen++
				}
			}
		}
		lastChar = char
	}
	if wordBufLen > 0 {
		lines[curLineNum] += string(wordBuf)
	}
	if controlCharNoneEndBuf.Len() > 0 {
		lines[curLineNum] += "\033[0m"
	}
	trimCurrent()

	rendered := strings.TrimRight(strings.Join(lines, "\n"), "\n")
	if !this.KeepEmptyLines {
		rxIsEmpty := regexp.MustCompile(`^(?:\s|\033\[[0-9;]+m)*$`)
		lines := strings.Split(rendered, "\n")
		filled := []string{}
		for _, line := range lines {
			if !rxIsEmpty.MatchString(line) {
				filled = append(filled, line)
			}
		}
		rendered = strings.Join(filled, "\n")
	}

	return rendered
}

func _wrapDebug(str string, args ...interface{}) {
	if dbg := os.Getenv("WRAP_DEBUG"); dbg == "yes" || dbg == "1" {
		fmt.Printf(str, args...)
	}
}
