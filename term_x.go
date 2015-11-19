// +build linux darwin freebsd netbsd openbsd solaris

package clif

import (
	"runtime"
	"syscall"
	"unsafe"
)

func init() {
	TermWidthCall = func() (int, error) {
		w := new(TermWindow)
		tio := syscall.TIOCGWINSZ
		if runtime.GOOS == "darwin" {
			tio = TERM_TIOCGWINSZ_OSX
		}
		res, _, err := syscall.Syscall(sys_ioctl,
			tty.Fd(),
			uintptr(tio),
			uintptr(unsafe.Pointer(w)),
		)
		if int(res) == -1 {
			return 0, err
		}
		TermWidthCurrent = int(w.Col)
		return TermWidthCurrent, nil
	}

	TermWidthCall()
}
