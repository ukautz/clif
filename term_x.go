// +build linux darwin freebsd netbsd openbsd solaris

package clif

import (
	"os"
	"runtime"
	"syscall"
	"unsafe"
)

func init() {
	TermWidthCall = func() (int, error) {
		w := new(termWindow)
		tio := syscall.TIOCGWINSZ
		if runtime.GOOS == "darwin" {
			tio = TERM_TIOCGWINSZ_OSX
		}
		res, _, err := syscall.Syscall(sys_ioctl,
			uintptr(syscall.Stdin),
			uintptr(tio),
			uintptr(unsafe.Pointer(w)),
		)
		if err != 0 || int(res) == -1 {
			return TERM_DEFAULT_WIDTH, os.NewSyscallError("GetWinsize", err)
		}
		return int(w.Col) - 4, nil
	}

	TermWidthCurrent, _ = TermWidthCall()
}
