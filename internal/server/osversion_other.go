//go:build !darwin

package server

// osVersion mirrors the attribution header on non-macOS platforms, where the
// app does not read a sysctl kernel release.
func osVersion() string {
	return "0.0"
}
