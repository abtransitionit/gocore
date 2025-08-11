package main

import (
	"runtime"
)

// name: GetOsType
// Return:
// string: the operating system type. Possible return values include:
//   - "darwin"   → macOS
//   - "linux"    → Linux
//   - "windows"  → Windows
func GetOsType() string {
	return runtime.GOOS
}

// name: GetOsCpuType
// Return :
// string: the operating system processor architecture. Possible return values include:
//   - "amd64"    → 64-bit x86
//   - "arm64"    → 64-bit ARM
//   - "386"      → 32-bit x86func getOSType() string {
//   - and others like "arm", "ppc64", etc.
func GetOsCpuType() string {
	return runtime.GOARCH
}
