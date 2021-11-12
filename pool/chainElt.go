package pool

import (
	"sync/atomic"
	"unsafe"
)

type poolChainElt[T any] struct {
	poolDequeue[T]

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
	next, prev *poolChainElt[T]
}

func storePoolChainElt[T any](pp **poolChainElt[T], v *poolChainElt[T]) {
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(pp)), unsafe.Pointer(v))
}

func loadPoolChainElt[T any](pp **poolChainElt[T]) *poolChainElt[T] {
	return (*poolChainElt[T])(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(pp))))
}
