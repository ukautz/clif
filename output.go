package clif

import (
	"fmt"
	"io"
	"os"
)

// Output is interface for outputter classes
type Output interface {
	// Printf writes to output
	Printf(msg string, args ...interface{})

	// SetFormatter is builder method and replaces current formatter
	SetFormatter(f Formatter) Output
}

type DefaultOutput struct {
	fmt Formatter
	io  io.Writer
}

func NewOutput(io io.Writer, f Formatter) *DefaultOutput {
	if io == nil {
		io = os.Stdout
	}
	return &DefaultOutput{
		fmt: f,
		io:  io,
	}
}

func NewPlainOutput(io io.Writer) *DefaultOutput {
	return NewOutput(io, NewDefaultFormatter(nil))
}

func NewFancyOutput(io io.Writer) *DefaultOutput {
	return NewOutput(io, NewDefaultFormatter(DefaultStyles))
}

func (this *DefaultOutput) SetFormatter(f Formatter) Output {
	this.fmt = f
	return this
}

func (this *DefaultOutput) Printf(msg string, args ...interface{}) {
	this.io.Write([]byte(this.fmt.Format(fmt.Sprintf(msg, args...))))
}
