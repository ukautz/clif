// +build windows

package clif

import (
	"github.com/olekukonko/ts"
)

func init() {
	TermWidthCall = func() (int, error) {
		size, err := ts.GetSize()
		TermWidthCurrent = size.Col()
		return TermWidthCurrent, err
	}

	TermWidthCall()
}