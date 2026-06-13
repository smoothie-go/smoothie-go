//go:build windows

package portable

import "syscall"

func IsOpenCLAvailable() bool {
	handle, err := syscall.LoadLibrary("OpenCL.dll")
	if err == nil {
		syscall.FreeLibrary(handle)
		return true
	}
	return false
}
