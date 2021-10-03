package futex

import (
	_ "unsafe"
)

//go:linkname Sleep runtime.futexsleep
func Sleep(addr *uint32, value uint32, ns int64)

//go:linkname Wake runtime.futexwakeup
func Wake(addr *uint32, count uint32)
