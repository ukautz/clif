package clif

import (
	"regexp"
	"strings"
	"fmt"
)

// Formatter is used by Output for rendering. It supports style directives in the form <info> or <end> or suchlike
type Formatter interface {

	// Escape escapes a string, so that no formatter tokens will be interpolated (eg `<foo>` -> `\<foo>`)
	Escape(msg string) string

	// Format renders message for output by applying <style> tokens
	Format(msg string) string
}

// DefaultFormatter strips all formatting from the output message
type DefaultFormatter struct {
	styles map[string]string
}

var DefaultStyles = map[string]string{
	"error":     "\033[31;1m",
	"warn":      "\033[33m",
	"info":      "\033[34m",
	"success":   "\033[32m",
	"debug":     "\033[30;1m",
	"headline":  "\033[4;1m",
	"subline":   "\033[4m",
	"important": "\033[47;30;1m",
	"query":     "\033[36m",
	"reset":     "\033[0m",
}

// http://misc.flogisoft.com/bash/tip_colors_and_formatting#colors1
var SunburnStyles = map[string]string{
	"error":     "\033[97;48;5;196;1m",
	"warn":      "\033[30;48;5;208;2m",
	"info":      "\033[38;5;142;2m",
	"success":   "\033[38;5;2;2m",
	"debug":     "\033[38;5;242;2m",
	"headline":  "\033[38;5;226;1m",
	"subline":   "\033[38;5;228;1m",
	"important": "\033[38;5;15;2;4m",
	"query":     "\033[38;5;77m",
	"reset":     "\033[0m",
}

// http://misc.flogisoft.com/bash/tip_colors_and_formatting#colors1
var WinterStyles = map[string]string{
	"error":     "\033[97;48;5;89;1m",
	"warn":      "\033[30;48;5;97;2m",
	"info":      "\033[38;5;69;2m",
	"success":   "\033[38;5;45;1m",
	"debug":     "\033[38;5;239;2m",
	"headline":  "\033[38;5;21;1m",
	"subline":   "\033[38;5;27;1m",
	"important": "\033[38;5;15;2;4m",
	"query":     "\033[38;5;111m",
	"reset":     "\033[0m",
}

var DebugStyles = map[string]string{
	"error":     "E:",
	"warn":      "W:",
	"info":      "I:",
	"success":   "S:",
	"debug":     "D:",
	"headline":  "H:",
	"subline":   "U:",
	"important": "P:",
	"query":     "Q:",
	"reset":     "R:",
}

// NewDefaultFormatter constructs a new constructor with the given styles
func NewDefaultFormatter(styles map[string]string) *DefaultFormatter {
	return &DefaultFormatter{styles}
}

// DefaultFormatterTokenRegex is a regex to find all tokens, used by the DefaultFormatter
var DefaultFormatterTokenRegex = regexp.MustCompile(`(<[^>]+>)`)

// DefaultFormatterPre is a pre-callback, before DefaultFormatterTokenRegex is used
var DefaultFormatterPre = func(msg string) string {
	return strings.Replace(msg, `\<`, "~~~~#~~~~", -1)
}

// DefaultFormatterPost is a post-callback, after DefaultFormatterTokenRegex is used
var DefaultFormatterPost = func(msg string) string {
	return strings.Replace(msg, "~~~~#~~~~", "<", -1)
}

// DefaultFormatterEscape is a callback to replace all tokens with escaped versions
var DefaultFormatterEscape = func(msg string) string {
	return DefaultFormatterTokenRegex.ReplaceAllStringFunc(msg, func(token string) string {
		return `\` + token
	})
}

func (this *DefaultFormatter) Escape(msg string) string {
	msg = DefaultFormatterPre(msg)
	msg = DefaultFormatterEscape(msg)
	msg = strings.Replace(msg, "~~~~#~~~~", "\\<", -1)
	return msg
}

func (this *DefaultFormatter) Format(msg string) string {
	// Since Go regexp does not implement lock-behinds: using a multi-pass approach
	// to first replace all escaped style tokens (\<token>) with string, which is
	// rather unlikely to occur, then replacing style tokens with color control
	// characters and then re-replacing the placeholder strings back
	msg = DefaultFormatterPre(msg)
	msg = DefaultFormatterTokenRegex.ReplaceAllStringFunc(msg, func(token string) string {
		style := token[1 : len(token)-1]
		if this.styles == nil {
			if _, ok := DefaultStyles[style]; ok {
				return ""
			} else {
				return token
			}
		} else if replace, ok := this.styles[style]; ok {
			return replace
		} else {
			return token
		}
	})
	msg = DefaultFormatterPost(msg)

	return msg
}

func init() {

	// for each token `foo` add a token `/foo`, which contains reset, so we can do "<error>bla</error>"
	// instead of "<error>bla<reset>".
	for _, m := range []*map[string]string{&DefaultStyles, &SunburnStyles, &WinterStyles} {
		for k, _ := range *m {
			(*m)[fmt.Sprintf("/%s", k)] = (*m)["reset"]
		}
	}
}