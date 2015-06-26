package clif

import (
	"fmt"
	"io"
	"os"
)

// Output is interface for
type Output interface {

	// Printf applies format (renders styles) and writes to output
	Printf(msg string, args ...interface{})

	// Sprintf applies format (renders styles) and returns as string
	Sprintf(msg string, args ...interface{}) string

	// SetFormatter is builder method and replaces current formatter
	SetFormatter(f Formatter) Output
}

// DefaultOutput is the default used output type
type DefaultOutput struct {
	fmt Formatter
	io  io.Writer
}

// NewOutput generates a new (default) output with provided io writer (if nil
// then `os.Stdout` is used) and a formatter
func NewOutput(io io.Writer, f Formatter) *DefaultOutput {
	if io == nil {
		io = os.Stdout
	}
	return &DefaultOutput{
		fmt: f,
		io:  io,
	}
}

// NewMonochromeOutput returns default output (on `os.Stdout`, if io is nil) using
// a formatter which strips all style tokens (hence rendering non-colored, plain
// strings)
func NewMonochromeOutput(io io.Writer) *DefaultOutput {
	return NewOutput(io, NewDefaultFormatter(nil))
}

// NewColoredOutput returns default output (on `os.Stdout`, if io is nil) using
// a formatter which applies the default color styles to style tokens on output.
func NewColoredOutput(io io.Writer) *DefaultOutput {
	return NewOutput(io, NewDefaultFormatter(DefaultStyles))
}

func (this *DefaultOutput) SetFormatter(f Formatter) Output {
	this.fmt = f
	return this
}

func (this *DefaultOutput) Printf(msg string, args ...interface{}) {
	this.io.Write([]byte(this.Sprintf(msg, args...)))
}

func (this *DefaultOutput) Sprintf(msg string, args ...interface{}) string {
	return this.fmt.Format(fmt.Sprintf(msg, args...))
}
