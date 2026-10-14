//go:build darwin

package server

import "syscall"

// osVersion reports the kernel release (e.g. "27.0.0" on macOS), which is
// what the ZCode app sends as x-os-version.
func osVersion() string {
	if version, err := syscall.Sysctl("kern.osrelease"); err == nil && version != "" {
		return version
	}
	return "0.0"
}
