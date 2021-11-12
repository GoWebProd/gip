package pool

import "unsafe"

// Local per-P Pool appendix.
type poolLocalInternal[T nilable[V], V any] struct {
	private T           // Can be used only by the respective P.
	shared  poolChain[T, V] // Local P can pushHead/popHead; any P can popTail.
}

type poolLocal[T nilable[V], V any] struct {
	poolLocalInternal[T, V]

	// Prevents false sharing on widespread platforms with
	// 128 mod (cache line size) = 0 .
	// pad [128 - unsafe.Sizeof(poolLocalInternal[T, V]{})%128]byte
}

func indexLocal[T nilable[V], V any](l unsafe.Pointer, i int) *poolLocal[T, V] {
	lp := unsafe.Pointer(uintptr(l) + uintptr(i)*unsafe.Sizeof(poolLocal[T, V]{}))
	return (*poolLocal[T, V])(lp)
}
