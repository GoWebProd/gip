package pool

import (
	"sync/atomic"
	"unsafe"
)

type poolChainElt[T nilable[V], V any] struct {
	poolDequeue[T, V]

	// next and prev link to the adjacent poolChainElts in this
	// poolChain.
	//
	// next is written atomically by the producer and read
	// atomically by the consumer. It only transitions from nil to
	// non-nil.
	//
	// prev is written atomically by the consumer and read
	// atomically by the producer. It only transitions from
	// non-nil to nil.
	next, prev *poolChainElt[T, V]
}

func storePoolChainElt[T nilable[V], V any](pp **poolChainElt[T, V], v *poolChainElt[T, V]) {
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(pp)), unsafe.Pointer(v))
}

func loadPoolChainElt[T nilable[V], V any](pp **poolChainElt[T, V]) *poolChainElt[T, V] {
	return (*poolChainElt[T, V])(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(pp))))
}
