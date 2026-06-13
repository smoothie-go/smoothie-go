//go:build windows

package recipe

import "syscall"

func isOpenCLAvailable() bool {
	handle, err := syscall.LoadLibrary("OpenCL.dll")
	if err == nil {
		syscall.FreeLibrary(handle)
		return true
	}
	return false
}
