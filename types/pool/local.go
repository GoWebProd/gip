package pool

import "unsafe"

// Local per-P Pool appendix.
type poolLocalInternal[T any] struct {
	private *T           // Can be used only by the respective P.
	shared  poolChain[T] // Local P can pushHead/popHead; any P can popTail.
}

type poolLocal[T any] struct {
	poolLocalInternal[T]

	// Prevents false sharing on widespread platforms with
	// 128 mod (cache line size) = 0 .
	pad [128 - unsafe.Sizeof(poolLocalInternal[T]{})%128]byte
}

func indexLocal[T any](l unsafe.Pointer, i int) *poolLocal[T] {
	lp := unsafe.Pointer(uintptr(l) + uintptr(i)*unsafe.Sizeof(poolLocal[T]{}))
	return (*poolLocal[T])(lp)
}
