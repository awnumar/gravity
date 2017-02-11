// +build windows

package memlock

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

// Lock is a wrapper for windows.VirtualLock()
func Lock(b []byte) error {
	p := unsafe.Pointer(&b[0])
	return windows.VirtualLock(uintptr(p), uintptr(len(b)))
}

// Unlock is a wrapper for windows.VirtualUnlock()
func Unlock(b []byte) error {
	p := unsafe.Pointer(&b[0])
	return windows.VirtualUnlock(uintptr(p), uintptr(len(b)))
}
