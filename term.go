// Determine terminal width which can come in handy for rendering tables and somesuch.
//
//
// DISCLAIMER:
// The code contents of all term*.go files is PROUDLY STOLEN FROM https://github.com/cheggaaa/pb
// which sadly does not export this nicely written functions and to whom all credits should go.
// Only slight modifications.
package clif

import (
	"os"
)

const (
	TERM_TIOCGWINSZ     = 0x5413
	TERM_TIOCGWINSZ_OSX = 1074295912
	TERM_DEFAULT_WIDTH  = 78
)

var (
	tty *os.File

	// TermWidthCall is the callback returning the terminal width
	TermWidthCall func() (int, error)

	// TermWidthCurrent contains the current terminal width from the last
	// call of `TerminalWidth()` (which is called in `init()`)
	TermWidthCurrent = TERM_DEFAULT_WIDTH
)

type (
	termWindow struct {
		Row    uint16
		Col    uint16
		Xpixel uint16
		Ypixel uint16
	}
)

// TerminalWidth returns the terminal width in amount of characters
func TermWidth() (int, error) {
	return TermWidthCall()
}
