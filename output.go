package clif

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// Output is interface for
type Output interface {

	// Escape escapes a string, so that no formatter tokens will be interpolated (eg `<foo>` -> `\<foo>`)
	Escape(s string) string

	// Printf applies format (renders styles) and writes to output
	Printf(msg string, args ...interface{})

	// ProgressBars returns the pool of progress bars
	ProgressBars() ProgressBarPool

	// Sprintf applies format (renders styles) and returns as string
	Sprintf(msg string, args ...interface{}) string

	// SetFormatter is builder method and replaces current formatter
	SetFormatter(f Formatter) Output

	// Table creates a table object
	Table(header []string, style ...*TableStyle) *Table

	// Writer returns the `io.Writer` used by this output
	Writer() io.Writer
}

// DefaultOutput is the default used output type
type DefaultOutput struct {
	fmt    Formatter
	io     io.Writer
	pbPool ProgressBarPool
}

var (
	DefaultOutputTableHeaderRenderer = func(out Output) func(string) string {
		return func(content string) string {
			return strings.Join(SplitFormattedString(out.Sprintf("<headline>%s<reset>", content)), "\n")
		}
	}
	DefaultOutputTableContentRenderer = func(out Output) func(string) string {
		return func(content string) string {
			from := strings.Replace(content, "\n", "<BR>", -1)
			from = rxControlCharacters.ReplaceAllString(from, "<CTRL>")
			//fmt.Printf("\n>> CONTENT FROM \"%s\"\n", from)
			return strings.Join(SplitFormattedString(out.Sprintf(content)), "\n")
		}
	}
)

// NewOutput generates a new (default) output with provided io writer (if nil
// then `os.Stdout` is used) and a formatter
func NewOutput(io io.Writer, f Formatter) *DefaultOutput {
	if io == nil {
		io = os.Stdout
	}
	return &DefaultOutput{
		fmt:    f,
		io:     io,
		pbPool: NewProgressBarPool(),
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
func NewColorOutput(io io.Writer) *DefaultOutput {
	return NewOutput(io, NewDefaultFormatter(DefaultStyles))
}

// NewDebugOutput is used for debugging the color formatter
func NewDebugOutput(io io.Writer) *DefaultOutput {
	return NewOutput(io, NewDefaultFormatter(DebugStyles))
}

func (this *DefaultOutput) SetFormatter(f Formatter) Output {
	this.fmt = f
	return this
}

func (this *DefaultOutput) Escape(msg string) string {
	return this.fmt.Escape(msg)
}

func (this *DefaultOutput) Printf(msg string, args ...interface{}) {
	this.io.Write([]byte(this.Sprintf(msg, args...)))
}

func (this *DefaultOutput) ProgressBars() ProgressBarPool {
	return this.pbPool
}

func (this *DefaultOutput) Sprintf(msg string, args ...interface{}) string {
	return this.fmt.Format(fmt.Sprintf(msg, args...))
}

func (this *DefaultOutput) Table(headers []string, style ...*TableStyle) *Table {
	if len(style) == 0 {
		style = []*TableStyle{NewDefaultTableStyle()}
	}
	style[0].HeaderRenderer = DefaultOutputTableHeaderRenderer(this)
	style[0].ContentRenderer = DefaultOutputTableContentRenderer(this)
	table := NewTable(headers, style...)
	return table
}

func (this *DefaultOutput) Writer() io.Writer {
	return this.io
}
