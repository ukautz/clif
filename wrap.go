package clif

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"
)

// WrapString wraps the given string within lim width in characters.
//
// PROUDLY STOLEN from https://raw.githubusercontent.com/mitchellh/go-wordwrap/master/wordwrap.go
//
// Wrapping is currently naive and only happens at white-space. A future
// version of the library will implement smarter wrapping. This means that
// pathological cases can dramatically reach past the limit, such as a very
// long word.
func WrapString(s string, lim uint) string {
	return wrapString(s, lim, false)
}

// WrapStringExtreme does also break words, if they are longer than the given limit
func WrapStringExtreme(s string, lim uint) string {
	return wrapString(s, lim, true)
}

func wrapStringX(s string, lineLimit uint, splitWords bool) string {
	lines := []string{""}
	wordBuf := ""
	wordBufLen := 0
	curLineNum := 0
	curLineLen := 0
	controlCharSeq := 0
	lastControlChars := []rune{}
	controlCharBuf := bytes.NewBuffer(nil)

	finishLine := func(add string, hasControlChars bool) {
		if hasControlChars {
			lines[curLineNum] += "\033[0m"
		}
		lines[curLineNum] += add
		lines = append(lines, "")
		curLineNum++
		if hasControlChars {
			lines[curLineNum] += controlCharBuf.String()
		}
		curLineLen = 0
	}

	var lastChar rune
	for _, char := range s {
		if IsControlCharStart(byte(char)) {
			controlCharBuf.WriteRune(char)
			controlCharSeq = 1
		} else if controlCharSeq == 1 { // expect "["
			if char != 91 { // abort .. not "["
				lines[curLineNum] += controlCharBuf.String()
				controlCharBuf.Reset()
				controlCharSeq = 0
			} else {
				controlCharBuf.WriteRune(char)
				controlCharSeq = 2
			}
		} else if controlCharSeq == 2 {
			if char <= '0' && char >= '9' {
				controlCharBuf.WriteRune(char)
				lastControlChars = []rune{char}
				controlCharSeq = 3
			} else { // abort, "not valid char
				lines[curLineNum] += controlCharBuf.String()
				controlCharBuf.Reset()
				controlCharSeq = 0
			}
		} else if controlCharSeq == 3 {
			if char <= '0' && char >= '9' {
				lastControlChars = append(lastControlChars, char)
				controlCharBuf.WriteRune(char)
			} else if char == ';' {
				lastControlChars = []rune{}
				controlCharBuf.WriteRune(char)
				controlCharSeq = 2
			} else if char == 'm' { // end
				controlCharBuf.WriteRune(char)
				controlCharSeq = 0
				lines[curLineNum] += controlCharBuf.String()
				if len(lastControlChars) == 1 && lastControlChars[0] == '0' {
					controlCharBuf.Reset()
				}
				lastControlChars = []rune{}
				controlCharSeq = 0
			} else { // abort, "not valid char
				lines[curLineNum] += controlCharBuf.String()
				controlCharBuf.Reset()
				controlCharSeq = 0
			}
		} else {
			hasControlChars := controlCharBuf.Len() > 0
			if char == '\n' {
				finishLine(wordBuf, hasControlChars)
				wordBuf = ""
				wordBufLen = 0
			} else if unicode.IsSpace(char) {
				if wordBufLen > 0 {
					lines[curLineNum] += wordBuf
					curLineLen += wordBufLen
					lines[curLineNum] += string(char)
					curLineLen ++
					wordBuf = ""
					wordBufLen = 0
				} else if curLineLen > 0 && unicode.IsSpace(lastChar) {
					fmt.Printf("\n>> HERE IN CONTINUE\n\n")
				} else {
					fmt.Printf("\n>> HERE IN SPACE\n\n")
					lines[curLineNum] += string(char)
					curLineLen++
				}
			} else {
				totalLineLen := curLineLen + wordBufLen
				if totalLineLen+1 > int(lineLimit) {
					if curLineLen > 0 { // has prefix
						finishLine(string(char), hasControlChars)
						wordBuf = ""
						wordBufLen = 0

					} else { // the word itself is longer than line
						if splitWords {
							finishLine(string(char), hasControlChars)
							wordBuf = ""
							wordBufLen = 0
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
	return strings.Join(lines, "\n")
}

func wrapString(s string, lim uint, splitWords bool) string {
	// Initialize a buffer with a slightly larger size to account for breaks
	init := make([]byte, 0, len(s))
	buf := bytes.NewBuffer(init)

	var current uint
	var wordBuf, spaceBuf bytes.Buffer

	for _, char := range s {
		if char == '\n' {
			if wordBuf.Len() == 0 {
				if current+uint(spaceBuf.Len()) > lim {
					current = 0
				} else {
					current += uint(spaceBuf.Len())
					spaceBuf.WriteTo(buf)
				}
				spaceBuf.Reset()
			} else {
				current += uint(spaceBuf.Len() + wordBuf.Len())
				spaceBuf.WriteTo(buf)
				spaceBuf.Reset()
				wordBuf.WriteTo(buf)
				wordBuf.Reset()
			}
			buf.WriteRune(char)
			current = 0
		} else if unicode.IsSpace(char) {
			if spaceBuf.Len() == 0 || wordBuf.Len() > 0 {
				current += uint(spaceBuf.Len() + wordBuf.Len())
				spaceBuf.WriteTo(buf)
				spaceBuf.Reset()
				wordBuf.WriteTo(buf)
				wordBuf.Reset()
			}

			spaceBuf.WriteRune(char)
		} else {
			if splitWords && uint(wordBuf.Len()) == lim {
				wordBuf.WriteRune('\n')
				spaceBuf.WriteTo(buf)
				spaceBuf.Reset()
				wordBuf.WriteTo(buf)
				wordBuf.Reset()
			}
			wordBuf.WriteRune(char)

			if current+uint(spaceBuf.Len()+wordBuf.Len()) > lim && uint(wordBuf.Len()) < lim {
				buf.WriteRune('\n')
				current = 0
				spaceBuf.Reset()
			}
		}
	}

	if wordBuf.Len() == 0 {
		if current+uint(spaceBuf.Len()) <= lim {
			spaceBuf.WriteTo(buf)
		}
	} else {
		spaceBuf.WriteTo(buf)
		wordBuf.WriteTo(buf)
	}

	return buf.String()
}
