package clif

import (
	"regexp"
)

// Formatter is used by Output for rendering. It supports style directives in the form <info> or <end> or suchlike
type Formatter interface {

	// Format renders message with args for output
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

func NewDefaultFormatter(styles map[string]string) *DefaultFormatter {
	return &DefaultFormatter{styles}
}

func (this *DefaultFormatter) Format(msg string) string {
	rx := regexp.MustCompile(`(<[^>]+>)`)
	return rx.ReplaceAllStringFunc(msg, func(token string) string {
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
}
