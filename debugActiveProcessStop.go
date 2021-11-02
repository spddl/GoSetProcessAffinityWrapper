package main

import "syscall"

// https://docs.microsoft.com/en-us/windows/win32/api/debugapi/nf-debugapi-debugactiveprocess

var (
	modkernel32 = syscall.NewLazyDLL("kernel32.dll")

	procDebugActiveProcessStop = modkernel32.NewProc("DebugActiveProcessStop")
)

func DebugActiveProcessStop(processid uint32) (err error) {
	r1, _, e1 := syscall.Syscall(procDebugActiveProcessStop.Addr(), 1, uintptr(processid), 0, 0)
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}
