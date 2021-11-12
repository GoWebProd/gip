package rtime

import "unsafe"

// noescape hides a pointer from escape analysis.  noescape is
// the identity function but escape analysis doesn't think the
// output depends on the input.  noescape is inlined and currently
// compiles down to zero instructions.
// USE CAREFULLY!
//go:nosplit
func Noescape[T any](p *T) unsafe.Pointer {
	x := uintptr(unsafe.Pointer(p))
	
	return unsafe.Pointer(x ^ 0)
}
