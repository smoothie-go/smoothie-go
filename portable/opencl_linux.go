//go:build linux

package portable

import "os"

func IsOpenCLAvailable() bool {
	if files, err := os.ReadDir("/etc/OpenCL/vendors"); err == nil && len(files) > 0 {
		return true
	}
	commonPaths := []string{
		"/usr/lib/libOpenCL.so",
		"/usr/lib/libOpenCL.so.1",
		"/usr/lib64/libOpenCL.so",
		"/usr/lib64/libOpenCL.so.1",
		"/usr/lib/x86_64-linux-gnu/libOpenCL.so",
		"/usr/lib/x86_64-linux-gnu/libOpenCL.so.1",
	}
	for _, path := range commonPaths {
		if _, err := os.Stat(path); err == nil {
			return true
		}
	}
	return false
}
