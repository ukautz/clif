// +build linux darwin freebsd netbsd openbsd

package clif

import "syscall"

const sys_ioctl = syscall.SYS_IOCTL
