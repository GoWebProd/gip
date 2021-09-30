package rtime

import (
	_ "unsafe"
)

//go:linkname ProcPin runtime.procPin
func ProcPin() int

//go:linkname ProcUnpin runtime.procUnpin
func ProcUnpin()
