package cli

import "fmt"

type Formatter interface {
	Sprintf(msg string, args ...interface{}) string
}

type DefaultFormatter struct{}

func (this *DefaultFormatter) Sprintf(msg string, args ...interface{}) string {
	return fmt.Sprintf(msg, args...)
}

func NewDefaultFormatter() *DefaultFormatter {
	return &DefaultFormatter{}
}