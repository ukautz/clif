package cli

import (
	"fmt"
	"io"
	"os"
)

type Output interface {
	// Printf writes output
	Printf(msg string, args ...interface{})
	Die(msg string, args ...interface{})
	SetFormatter(f Formatter) Output
	SetDieHandler(func(msg string, args ...interface{})) Output
}

type DefaultOutput struct {
	fmt Formatter
	die func(msg string, args ...interface{})
}

func newOutput() *DefaultOutput {
	return &DefaultOutput{
		fmt: NewDefaultFormatter(),
		die: func(msg string, args ...interface{}) {
			fmt.Printf(msg+ "\n", args...)
			os.Exit(1)
		},
	}
}

func (this *DefaultOutput) SetFormatter(f Formatter) Output {
	this.fmt = f
	return this
}

func (this *DefaultOutput) SetDieHandler(die func(msg string, args ...interface{})) Output {
	this.die = die
	return this
}

func (this *DefaultOutput) Die(msg string, args ...interface{}) {
	this.die(msg, args...)
}

func (this *DefaultOutput) Printf(msg string, args ...interface{}) {
	fmt.Printf(msg, args...)
}

type IoOutput struct {
	DefaultOutput
	io io.Writer
}

func NewIoOutput(io io.Writer) *IoOutput {
	this := &IoOutput{
		DefaultOutput: *newOutput(),
		io:     io,
	}
	this.SetDieHandler(func(msg string, args ...interface{}) {
		this.Printf(msg+ "\n", args...)
		os.Exit(1)
	})
	return this
}

func (this *IoOutput) Printf(msg string, args ...interface{}) {
	this.io.Write([]byte(this.fmt.Sprintf(msg, args...)))
}
