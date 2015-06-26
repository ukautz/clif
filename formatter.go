package clif

import (
	"regexp"
	"strings"
)

// Formatter is used by Output for rendering. It supports style directives in the form <info> or <end> or suchlike
type Formatter interface {

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
	"info":      "\033[36m",
	"success":   "\033[32m",
	"debug":     "\033[30;1m",
	"headline":  "\033[4;1m",
	"subline":   "\033[4m",
	"important": "\033[47;30;1m",
	"query":     "\033[34m",
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

// NewDefaultFormatter constructs a new constructor with the given styles
func NewDefaultFormatter(styles map[string]string) *DefaultFormatter {
	return &DefaultFormatter{styles}
}

func (this *DefaultFormatter) Format(msg string) string {
	// Since Go regexp does not implement lock-behinds: using a multi-pass approach
	// to first replace all escaped style tokens (\<token>) with string, which is
	// rather unlikely to occur, then replacing style tokens with color control
	// characters and then re-replacing the placeholder strings back
	msg = strings.Replace(msg, `\<`, "~~~~#~~~~", -1)
	rxToken := regexp.MustCompile(`(<[^>]+>)`)
	msg = rxToken.ReplaceAllStringFunc(msg, func(token string) string {
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
	msg = strings.Replace(msg, "~~~~#~~~~", "<", -1)

	return msg
}
