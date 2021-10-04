package rtime

import (
	"unsafe"
)

//go:linkname ProcPin runtime.procPin
func ProcPin() int

//go:linkname ProcUnpin runtime.procUnpin
func ProcUnpin()

//go:linkname GoPark runtime.gopark
func GoPark(unlockf func(unsafe.Pointer, unsafe.Pointer) bool, lock unsafe.Pointer, reason uint8, traceEv byte, traceskip int)

//go:linkname GoReady runtime.goready
func GoReady(gp unsafe.Pointer, traceskip int)
